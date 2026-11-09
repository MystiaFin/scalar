@echo off
rem run migrate   create or update the database
rem run seed      migrate, then add example transactions
rem run serve     start the api server
rem run test      run the tests

cd /d "%~dp0backend"

where go >nul 2>nul
if errorlevel 1 (
  echo go not found, run this, reopen the terminal, then try again:
  echo winget install GoLang.Go
  exit /b 1
)

if "%~1"=="migrate" (
  go run ./cmd/migrate
) else if "%~1"=="seed" (
  go run ./cmd/seed
) else if "%~1"=="serve" (
  go run ./cmd/server
) else if "%~1"=="test" (
  go test ./...
) else (
  echo usage: run [migrate^|seed^|serve^|test]
  exit /b 1
)
