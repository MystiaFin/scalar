"""Re-run Scalar's LLM step on scalar_llm_dataset.jsonl, zero-shot or few-shot.

Prompt layout and JSON schema follow AskOllama in scalar-engine
internal/email/ollama.go, with the 12-category guide used by the Colab eval.
Few-shot examples mirror scalar-engine eval/fewshot_seed.py.

    python3 run_eval.py scalar_llm_dataset.jsonl zeroshot eval_zeroshot_snippet.csv
    python3 run_eval.py scalar_llm_dataset.jsonl fewshot  eval_fewshot_snippet.csv
"""
import csv
import json
import random
import sys
import time
import urllib.request

MODEL = "llama3:8b"
SEED = 42
MAX_EXAMPLES = 10
CATS = ["food & drink", "grocery", "transport", "subscription", "health", "shopping",
        "rent", "installment", "utilities", "entertainment", "transfer", "other"]

SCHEMA = {
    "type": "object",
    "properties": {
        "merchant": {"type": "string"},
        "amount": {"type": "number"},
        "category": {"type": "string", "enum": CATS},
        "is_expense": {"type": "boolean"},
    },
    "required": ["merchant", "amount", "category", "is_expense"],
}

PROMPT = """You are parsing a financial transaction receipt or email.
Extract the merchant name, amount, category, and whether it is an expense.

Category guide:
- food & drink  : restaurants, cafes, street food, coffee shops, food delivery
- grocery       : supermarkets, minimarkets and convenience stores (Alfamart, Indomaret, FamilyMart, Lawson, Superindo), grocery stalls, even when buying a snack or drink
- transport     : ride-sharing (Grab, Gojek), public transit (TransJakarta), toll, parking, fuel, flights
- subscription  : recurring paid plans (Google One, Spotify, Netflix, Claude)
- health        : pharmacy, clinic, hospital, doctor, supplements, jamu
- shopping      : clothing, electronics, books, household items, marketplaces (Shopee, Tokopedia)
- rent          : rent
- installment   : installment (cicilan)
- utilities     : electricity, water, internet, mobile phone plans
- entertainment : cinema, gaming, events, hobbies
- transfer      : bank transfers, e-wallet top-ups, moving money between the user's own accounts
- other         : anything that clearly does not fit the above

{few_shot}Subject: {subject}
Body: {body}

Rules:
- merchant: extract ONLY the core primary brand name. Strip out payment gateways, branches, store codes, corporate suffixes (PT, Tbk, LLC) and ignore order IDs.
- amount: total transaction value as a plain number, including tax and fees. Strip ALL currency symbols. "Rp15.000" = 15000. Use 0 if not found.
- category: pick the single best fit from the list above.
- is_expense: true if money goes out to someone else, false if money comes in or moves between the user's own accounts.

Respond only with the JSON object."""


def norm(text):
    return " ".join(str(text).lower().split())


# --- few-shot selection, same as scalar-engine eval/fewshot_seed.py ---
def build_pool(rows):
    unique = {}
    for row in rows:
        gold = row["gold"]
        if gold["is_receipt"] != 1:
            continue
        key = (norm(gold["merchant"]), gold["category"], gold["is_expense"])
        unique.setdefault(key, {"merchant": gold["merchant"], "category": gold["category"],
                                "is_expense": gold["is_expense"]})
    return [unique[key] for key in sorted(unique, key=str)]


def select_examples(row, pool):
    merchant = row["gold"]["merchant"]
    excluded = norm(merchant) if merchant else None
    candidates = [ex for ex in pool if norm(ex["merchant"]) != excluded]
    rng = random.Random(f"{SEED}:{row['id']}")
    return rng.sample(candidates, min(MAX_EXAMPLES, len(candidates)))


def build_few_shot_block(examples):
    if not examples:
        return ""
    lines = ["Here are past transactions that were confirmed correct — use these as reference:\n"]
    for ex in examples:
        is_expense = "true" if ex["is_expense"] else "false"
        lines.append(f"  merchant: {ex['merchant']} → category: {ex['category']}, is_expense: {is_expense}\n")
    lines.append("\n")
    return "".join(lines)


def ask(subject, body, few_shot):
    req = json.dumps({
        "model": MODEL,
        "prompt": PROMPT.format(few_shot=few_shot, subject=subject, body=body),
        "format": SCHEMA,
        "stream": False,
        "options": {"temperature": 0, "seed": SEED},
    }).encode()
    r = urllib.request.urlopen(urllib.request.Request(
        "http://localhost:11434/api/generate", req, {"Content-Type": "application/json"}), timeout=900)
    raw = json.loads(r.read())["response"]
    return raw, json.loads(raw)


def main(src, mode, dst):
    rows = [json.loads(l) for l in open(src, encoding="utf-8") if l.strip()]
    pool = build_pool(rows)
    fields = ["merchant", "amount", "category", "is_expense"]
    with open(dst, "w", newline="", encoding="utf-8") as f:
        w = csv.writer(f)
        w.writerow(["id", "is_receipt"] + [f"{p}_{k}" for k in fields for p in ("gold", "pred")]
                   + ["seconds", "error", "raw"])
        for i, row in enumerate(rows, 1):
            few_shot = build_few_shot_block(select_examples(row, pool)) if mode == "fewshot" else ""
            t = time.time()
            raw, pred, err = "", {}, ""
            try:
                raw, pred = ask(row["subject"], row["body"], few_shot)
            except Exception as e:
                err = repr(e)
            g = row["gold"]
            vals = []
            for k in fields:
                vals += [g[k], pred.get(k)]
            w.writerow([row["id"], g["is_receipt"]] + vals + [round(time.time() - t, 1), err, raw])
            f.flush()
            print(f"{mode} {i}/{len(rows)} {row['id']} {time.time() - t:.1f}s {err}", flush=True)


if __name__ == "__main__":
    main(*sys.argv[1:4])
