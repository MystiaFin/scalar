import '../models/transaction.dart';
import 'format.dart';

/// Expense totals per category for one month, largest first.
/// Income is left out since the chart only shows where money went.
List<MapEntry<String, double>> expenseTotals(
  List<Transaction> txns,
  DateTime month,
) {
  final totals = <String, double>{};
  for (final t in txns) {
    if (!t.isExpense || !isSameMonth(t.date, month)) continue;
    totals[t.category] = (totals[t.category] ?? 0) + t.amount;
  }
  return totals.entries.toList()..sort((a, b) => b.value.compareTo(a.value));
}
