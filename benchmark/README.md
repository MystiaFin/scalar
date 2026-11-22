# benchmark

Reproduces the zero-shot vs few-shot comparison from the paper (Table 5): Llama 3 8B (`llama3:8b` on Ollama) extracts merchant, amount, category and is_expense from 127 emails that pass Scalar's keyword filter (113 transactions, 14 non-transactions). Each email body is cut to 500 characters.

The scripts and results are the ones used for the paper, unchanged.

| File | What it is |
|---|---|
| `run_eval.py` | Sends each email to Ollama, zero-shot or few-shot, and writes one CSV row per email |
| `metrics.py` | Scores the CSVs: accuracy, per-class precision/recall/F1, macro-F1, errors |
| `prompts/extract.txt` | The prompt template (same text as `PROMPT` in `run_eval.py`) |
| `prompts/schema.json` | The JSON schema passed to Ollama's `format` |
| `prompts/fewshot_example_001.txt` | What the few-shot block looks like for email 001 |
| `results/2026-11-11/` | The run reported in the paper: both CSVs, `report.txt` and the run logs |
| `data/sample.jsonl` | 6 made-up emails, only for checking that everything runs |

Settings: temperature 0, seed 42. Few-shot adds up to 10 lines of `merchant: X → category: Y, is_expense: Z` taken from the labels of other emails, never the same merchant as the email being tested.

## Data

The real dataset comes from the configured Gmail account and holds personal transactions, so it is not in this repository. If you have a copy, save it as `data/dataset.jsonl`. Git ignores everything in `data/` except the sample.

Each line is one email:

```json
{"id": "s01", "subject": "...", "body": "...", "gold": {"is_receipt": 1, "merchant": "...", "amount": 18000, "category": "food & drink", "is_expense": true}}
```

`body` is the text Scalar sends to the model, with personal details replaced by tags such as `[NAMA_PENGGUNA]` and `[REKENING]`. `is_receipt` is 0 for emails that are not real transactions (declined, pending, Rp0). See `data/sample.jsonl` for complete lines.

## Run it on Windows

1. Install Python and Ollama, then reopen the terminal:
   ```
   winget install Python.Python.3.12
   winget install Ollama.Ollama
   ```
2. Download the model (about 4.7 GB):
   ```
   ollama pull llama3:8b
   ```
3. From this folder, check that everything works on the sample, then run the real dataset:
   ```
   .\run data\sample.jsonl
   .\run
   ```

## Run it on Linux or macOS

Install Python 3 and [Ollama](https://ollama.com/download), then:

```bash
ollama pull llama3:8b
./run.sh data/sample.jsonl
./run.sh
```

## Output

The run writes `out/eval_zeroshot_snippet.csv`, `out/eval_fewshot_snippet.csv` and `out/report.txt`. Compare `out/report.txt` with `results/2026-11-11/report.txt`. The paper's run took about 9 minutes per mode on an RTX 3050. Other GPUs can give slightly different answers even at temperature 0.

The scripts expect Ollama on `http://localhost:11434`. The report's "of 14" and "of 127" labels assume the full dataset, so they are wrong on the sample.

To score saved results without running the model:

```
python metrics.py results/2026-11-11/eval_zeroshot_snippet.csv results/2026-11-11/eval_fewshot_snippet.csv
```

## What this does not reproduce

Tables 3 and 4 of the paper (full email text, 77.0% merchant, 81.4% category) come from a separate Google Colab notebook that is not included here. This benchmark covers Table 5 only.
