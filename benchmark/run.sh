#!/usr/bin/env bash
# ./run.sh                       run the benchmark on data/dataset.jsonl
# ./run.sh data/sample.jsonl     run it on another file, such as the synthetic sample
set -e

cd "$(dirname "$0")"

DATA="${1:-data/dataset.jsonl}"

if [ ! -f "$DATA" ]; then
	echo "$DATA not found, see README.md, or try: ./run.sh data/sample.jsonl"
	exit 1
fi

PY=python3
command -v python3 >/dev/null || PY=python

mkdir -p out
$PY run_eval.py "$DATA" zeroshot out/eval_zeroshot_snippet.csv
$PY run_eval.py "$DATA" fewshot out/eval_fewshot_snippet.csv
$PY metrics.py out/eval_zeroshot_snippet.csv out/eval_fewshot_snippet.csv | tee out/report.txt
