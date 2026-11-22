"""Score eval CSVs from run_eval.py. Prints plain-text tables.

    python3 metrics.py eval_zeroshot_snippet.csv eval_fewshot_snippet.csv
"""
import csv
import sys
from collections import Counter

CATS = ["food & drink", "grocery", "transport", "subscription", "health", "shopping",
        "rent", "installment", "utilities", "entertainment", "transfer", "other"]


def norm(text):
    return " ".join(str(text).lower().split())


def to_float(v):
    try:
        return float(v)
    except (TypeError, ValueError):
        return None


def load(path):
    rows = list(csv.DictReader(open(path, encoding="utf-8")))
    for r in rows:
        r["ok"] = not r["error"]
    return rows


def score(rows):
    rec = [r for r in rows if r["is_receipt"] == "1"]
    non = [r for r in rows if r["is_receipt"] == "0"]
    n = len(rec)
    m = sum(r["ok"] and norm(r["gold_merchant"]) == norm(r["pred_merchant"]) for r in rec)
    a = sum(r["ok"] and to_float(r["gold_amount"]) == to_float(r["pred_amount"]) for r in rec)
    e = sum(r["ok"] and r["gold_is_expense"].lower() == r["pred_is_expense"].lower() for r in rec)
    c = sum(r["ok"] and r["gold_category"] == r["pred_category"] for r in rec)
    s = {"n": n, "merchant": m / n, "amount": a / n, "is_expense": e / n, "category": c / n,
         "fake": sum(r["ok"] for r in non), "n_non": len(non),
         "failed": sum(not r["ok"] for r in rows),
         "seconds": sum(float(r["seconds"]) for r in rows)}
    per = {}
    gold_classes = sorted({r["gold_category"] for r in rec}, key=CATS.index)
    for k in CATS:
        tp = sum(r["ok"] and r["gold_category"] == k and r["pred_category"] == k for r in rec)
        fp = sum(r["ok"] and r["gold_category"] != k and r["pred_category"] == k for r in rec)
        fn = sum(r["gold_category"] == k and not (r["ok"] and r["pred_category"] == k) for r in rec)
        p = tp / (tp + fp) if tp + fp else 0.0
        rc = tp / (tp + fn) if tp + fn else 0.0
        f1 = 2 * p * rc / (p + rc) if p + rc else 0.0
        per[k] = (tp + fn, tp, fp, fn, p, rc, f1)
    s["macro_f1"] = sum(per[k][6] for k in gold_classes) / len(gold_classes)
    s["per"], s["gold_classes"] = per, gold_classes
    s["conf"] = Counter((r["gold_category"], r["pred_category"] if r["ok"] else "<failed>")
                        for r in rec if not (r["ok"] and r["gold_category"] == r["pred_category"]))
    s["merchant_err"] = [(r["id"], r["gold_merchant"], r["pred_merchant"] if r["ok"] else "<failed>")
                         for r in rec if not (r["ok"] and norm(r["gold_merchant"]) == norm(r["pred_merchant"]))]
    s["amount_err"] = [(r["id"], r["gold_amount"], r["pred_amount"] if r["ok"] else "<failed>")
                       for r in rec if not (r["ok"] and to_float(r["gold_amount"]) == to_float(r["pred_amount"]))]
    s["non"] = [(r["id"], r["pred_merchant"], r["pred_amount"], r["pred_category"], r["pred_is_expense"])
                for r in non]
    s["pred_cat"] = {r["id"]: (r["gold_category"], r["pred_category"] if r["ok"] else "<failed>") for r in rec}
    return s


def main(paths):
    names = [p.split("_")[1] for p in paths]
    S = [score(load(p)) for p in paths]
    print("SUMMARY (receipts n=%d)" % S[0]["n"])
    print("%-28s" % "metric" + "".join("%14s" % x for x in names))
    for k in ("merchant", "amount", "is_expense", "category", "macro_f1"):
        print("%-28s" % (k + (" acc" if k != "macro_f1" else "")) + "".join("%14.3f" % s[k] for s in S))
    print("%-28s" % "fake transactions (of 14)" + "".join("%14d" % s["fake"] for s in S))
    print("%-28s" % "failed calls (of 127)" + "".join("%14d" % s["failed"] for s in S))
    print("%-28s" % "runtime (min)" + "".join("%14.1f" % (s["seconds"] / 60) for s in S))

    for name, s in zip(names, S):
        print("\nPER-CLASS %s  (classes with gold samples in macro-F1; * = no gold samples)" % name)
        print("%-15s%5s%5s%5s%5s%8s%8s%8s" % ("class", "n", "TP", "FP", "FN", "P", "R", "F1"))
        for k in CATS:
            n, tp, fp, fn, p, rc, f1 = s["per"][k]
            if n or fp:
                star = "" if k in s["gold_classes"] else "*"
                print("%-15s%5d%5d%5d%5d%8.3f%8.3f%8.3f" % (k + star, n, tp, fp, fn, p, rc, f1))
        print("\nCONFUSION (wrong only) %s: gold -> pred : count" % name)
        for (g, p), c in sorted(s["conf"].items(), key=lambda x: -x[1]):
            print("  %s -> %s : %d" % (g, p, c))
        print("\nMERCHANT ERRORS %s (%d): id | gold | pred" % (name, len(s["merchant_err"])))
        for x in s["merchant_err"]:
            print("  %s | %s | %s" % x)
        print("\nAMOUNT ERRORS %s (%d): id | gold | pred" % (name, len(s["amount_err"])))
        for x in s["amount_err"]:
            print("  %s | %s | %s" % x)
        print("\nNON-RECEIPTS %s: id | merchant | amount | category | is_expense" % name)
        for x in s["non"]:
            print("  %s | %s | %s | %s | %s" % x)

    if len(S) == 2:
        print("\nCATEGORY CHANGES %s -> %s: id | gold | %s | %s" % (names[0], names[1], *names))
        for i, (g, p0) in S[0]["pred_cat"].items():
            p1 = S[1]["pred_cat"][i][1]
            if p0 != p1:
                mark = "fixed" if p1 == g else ("broke" if p0 == g else "still wrong")
                print("  %s | %s | %s | %s  (%s)" % (i, g, p0, p1, mark))


if __name__ == "__main__":
    main(sys.argv[1:])
