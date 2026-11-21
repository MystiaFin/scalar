# scalar_app

Flutter client for the scalar backend. It lists the transactions the backend pulled out of your email, lets you fix the category the extractor guessed, and shows where the money went each month.

## Features

- **Transactions** grouped by date, filtered by month and category
- **Edit** a transaction's category, type and note. "Save as default" applies the category to every transaction from that merchant and teaches the extractor
- **Overview** pie chart of the month's expenses per category. Tap a slice to see those transactions
- **Settings** for the backend url and Catppuccin Mocha / Latte theme, both saved on the device

## Layout

```
lib/
├── main.dart                      # theme wiring, starts ShellPage
├── models/
│   ├── category.dart              # category list and icons, same as the backend
│   └── transaction.dart
├── screens/
│   ├── shell_page.dart            # bottom navigation
│   ├── home_page.dart             # transactions list
│   ├── chart_page.dart            # overview pie chart
│   ├── settings_page.dart
│   └── theme_selection_page.dart
├── services/
│   ├── api_client.dart            # GET / PATCH against the saved base url
│   ├── settings_service.dart
│   └── transaction_service.dart
├── state/app_state.dart           # tab, month, filter and refresh notifiers
├── theme/                         # AppColors, Mocha, Latte, AppTheme
├── utils/                         # formatting and totals
└── widgets/                       # transaction row, edit sheet, shared bits
```

## Backend

The app uses these routes from `backend/internal/routes`:

| Method | Path | Used for |
|---|---|---|
| `GET` | `/api/transactions` | transactions list and overview |
| `PATCH` | `/api/transactions/:id` | edit sheet, body has `category`, `is_expense`, `description`, `confirmed` |

Start it from the repo root with `./run.sh serve` (or `.\run serve` on Windows). It listens on port `7225`. Run `seed` first if you want example data.

## Running

```bash
flutter pub get
flutter run
```

Then open **Settings** and set the base url:

| Where the app runs | Base URL |
|---|---|
| Desktop or web on the same machine | `http://localhost:7225` |
| Android emulator | `http://10.0.2.2:7225` |
| Phone on the same wifi | `http://<your computer's ip>:7225` |

## Checks

```bash
flutter analyze
flutter test
```
