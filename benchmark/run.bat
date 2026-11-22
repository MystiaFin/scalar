@echo off
rem run                       run the benchmark on data\dataset.jsonl
rem run data\sample.jsonl     run it on another file, such as the synthetic sample

cd /d "%~dp0"

set "DATA=%~1"
if "%DATA%"=="" set "DATA=data\dataset.jsonl"

if not exist "%DATA%" (
  echo %DATA% not found, see README.md, or try: run data\sample.jsonl
  exit /b 1
)

where python >nul 2>nul
if errorlevel 1 (
  echo python not found, run this, reopen the terminal, then try again:
  echo winget install Python.Python.3.12
  exit /b 1
)

set PYTHONUTF8=1
if not exist out mkdir out
python run_eval.py "%DATA%" zeroshot out\eval_zeroshot_snippet.csv || exit /b 1
python run_eval.py "%DATA%" fewshot out\eval_fewshot_snippet.csv || exit /b 1
python metrics.py out\eval_zeroshot_snippet.csv out\eval_fewshot_snippet.csv > out\report.txt || exit /b 1
type out\report.txt
