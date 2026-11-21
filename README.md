# scalar

Go backend in `backend/`, Flutter app in `app/`.

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
