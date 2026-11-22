# run inside the dev shell: `nix develop`, then `just <recipe>`
# or in one go: `nix develop -c just dev`

set dotenv-load
set dotenv-filename := "backend/.env"

ollama_url := env("OLLAMA_URL", "http://localhost:11434")
ollama_model := env("OLLAMA_MODEL", "llama3:8b")
api_port := "7225"
api_url := "http://localhost:" + api_port

# list recipes
default:
    @just --list

# first time setup: .env, dependencies, database and the ollama model
setup: env-file ollama-up
    cd backend && go mod download
    cd app && flutter pub get
    cd backend && go run ./cmd/migrate
    @just pull-model
    @echo "setup done, run: just dev"

# run ollama, the backend and the app together (ctrl+c or q stops everything)
dev device="linux" api=api_url: env-file ollama-up migrate
    #!/usr/bin/env bash
    set -euo pipefail

    mkdir -p .dev
    (cd backend && go build -o ../.dev/scalar-server ./cmd/server)

    (cd backend && exec ../.dev/scalar-server) > .dev/backend.log 2>&1 &
    backend=$!
    trap 'kill $backend 2>/dev/null || true; just stop' EXIT

    for _ in $(seq 1 40); do
        curl -sf "{{api_url}}/api/transactions" > /dev/null && break
        kill -0 $backend 2>/dev/null || { cat .dev/backend.log; exit 1; }
        sleep 0.25
    done

    echo "backend on {{api_url}}, log: .dev/backend.log"

    cd app && flutter run -d {{device}} --dart-define=SCALAR_API_URL={{api}}

# start only the backend, in the foreground
serve: env-file ollama-up
    cd backend && go run ./cmd/server

# start only the app, against a backend that is already running
app device="linux" api=api_url:
    cd app && flutter run -d {{device}} --dart-define=SCALAR_API_URL={{api}}

# create or update backend/scalar.db
migrate:
    cd backend && go run ./cmd/migrate

# migrate, then add example transactions
seed:
    cd backend && go run ./cmd/seed

# delete the database and create it again
reset-db:
    rm -f backend/scalar.db
    @just migrate

# backend and app tests
test:
    cd backend && go test ./...
    cd app && flutter test

# go vet and flutter analyze
check:
    cd backend && go vet ./...
    cd app && flutter analyze

# start ollama unless it is already reachable at OLLAMA_URL
ollama-up:
    #!/usr/bin/env bash
    set -euo pipefail

    if curl -sf "{{ollama_url}}/api/tags" > /dev/null; then
        echo "ollama is running at {{ollama_url}}"
        exit 0
    fi

    mkdir -p .dev
    echo "starting ollama at {{ollama_url}}, log: .dev/ollama.log"

    OLLAMA_HOST="{{ollama_url}}" nohup ollama serve > .dev/ollama.log 2>&1 &
    echo $! > .dev/ollama.pid

    for _ in $(seq 1 40); do
        curl -sf "{{ollama_url}}/api/tags" > /dev/null && exit 0
        sleep 0.25
    done

    echo "ollama did not start, see .dev/ollama.log"
    exit 1

# download OLLAMA_MODEL if it is not there yet
pull-model: ollama-up
    OLLAMA_HOST="{{ollama_url}}" ollama show {{ollama_model}} > /dev/null 2>&1 \
        || OLLAMA_HOST="{{ollama_url}}" ollama pull {{ollama_model}}

# stop the ollama this justfile started (one that was already running is left alone)
stop:
    #!/usr/bin/env bash
    if [ -f .dev/ollama.pid ]; then
        kill "$(cat .dev/ollama.pid)" 2>/dev/null || true
        rm -f .dev/ollama.pid
        echo "stopped ollama"
    fi

# create backend/.env from the example if it does not exist (email left blank so sync stays off)
env-file:
    #!/usr/bin/env bash
    if [ ! -f backend/.env ]; then
        sed -E 's/^(SCALAR_EMAIL|SCALAR_APP_PASSWORD)=.*/\1=/' backend/.env.example > backend/.env
        echo "created backend/.env, fill in SCALAR_EMAIL and SCALAR_APP_PASSWORD to sync email"
    fi
