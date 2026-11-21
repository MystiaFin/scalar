const _monthNames = [
  'January',
  'February',
  'March',
  'April',
  'May',
  'June',
  'July',
  'August',
  'September',
  'October',
  'November',
  'December',
];

/// "November 2026"
String formatMonth(DateTime dt) => '${_monthNames[dt.month - 1]} ${dt.year}';

/// "November 21, 2026"
String formatDate(DateTime dt) =>
    '${_monthNames[dt.month - 1]} ${dt.day}, ${dt.year}';

/// 1250000 -> "1,250,000", 12.5 -> "12.50"
String formatAmount(double amount) {
  final whole = amount.truncateToDouble() == amount;
  return amount
      .toStringAsFixed(whole ? 0 : 2)
      .replaceAllMapped(
        RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))'),
        (m) => '${m[1]},',
      );
}

/// Expenses get a minus sign, income gets a plus.
String formatSignedAmount(double amount, bool isExpense) =>
    '${isExpense ? '-' : '+'}${formatAmount(amount)}';

DateTime monthOf(DateTime dt) => DateTime(dt.year, dt.month);

bool isSameMonth(DateTime a, DateTime b) =>
    a.year == b.year && a.month == b.month;

/// The current month first, then the [count] - 1 months before it.
List<DateTime> recentMonths({DateTime? now, int count = 12}) {
  final today = now ?? DateTime.now();
  return List.generate(count, (i) => DateTime(today.year, today.month - i));
}
