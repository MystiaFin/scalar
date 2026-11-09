#!/usr/bin/env bash
# ./run.sh migrate   create or update the database
# ./run.sh seed      migrate, then add example transactions
# ./run.sh test      run the tests
set -e

cd "$(dirname "$0")"

if ! command -v go >/dev/null && command -v nix >/dev/null; then
	exec nix develop -c ./run.sh "$@"
fi

cd backend

case "$1" in
	migrate) go run ./cmd/migrate ;;
	seed) go run ./cmd/seed ;;
	test) go test ./... ;;
	*)
		echo "usage: ./run.sh [migrate|seed|test]"
		exit 1
		;;
esac
