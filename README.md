THIS IS MY COLLEGE ASSIGNMENT <br>
NOT A REAL PROJECT

# scalar

Go backend in `backend/`, Flutter app in `app/`.

## Quick start (Linux and macOS with Nix)

The dev shell in `flake.nix` has Go, Flutter, Ollama and `just`.

```bash
nix develop
just setup   # backend/.env, dependencies, database, pulls the ollama model
just dev     # starts ollama, the backend and the app together
```

`just dev` runs the Linux desktop app pointed at `http://localhost:7225`. For another device, pass its id and the url it can reach the backend on, for example `just dev emulator-5554 http://10.0.2.2:7225`. Run `just` to see every recipe.

To sync email, fill in `SCALAR_EMAIL` and `SCALAR_APP_PASSWORD` in `backend/.env`. `OLLAMA_URL` and `OLLAMA_MODEL` are optional. If Ollama is already running, `just` uses it. Otherwise it starts its own and stops it when `just dev` exits.

## Database

| Linux and macOS | Windows | What it does |
|---|---|---|
| `./run.sh migrate` | `.\run migrate` | Creates or updates `backend/scalar.db` |
| `./run.sh seed` | `.\run seed` | Migrates, then adds 20 made-up transactions |
| `./run.sh serve` | `.\run serve` | Starts the API on port 7225 |
| `./run.sh test` | `.\run test` | Tests both on a temporary database |

To change the schema, add the next numbered file to `backend/internal/db/migrations/`, such as `002_add_something.sql`. Don't edit a migration that has already been pushed.

The backend needs Go, and `.\run` tells you how to install it on Windows.

## App

`app/` is the Flutter client. Start the backend with `./run.sh serve`, then:

```bash
cd app
flutter pub get
flutter run
```

Set the backend url in the app's Settings tab. [app/README.md](app/README.md) lists the url to use for an emulator or a phone.
