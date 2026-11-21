import 'package:flutter_test/flutter_test.dart';

import 'package:scalar_app/models/transaction.dart';
import 'package:scalar_app/utils/totals.dart';

Transaction _tx(
  String category,
  double amount,
  DateTime date, {
  bool isExpense = true,
}) {
  return Transaction(
    id: 0,
    merchant: 'm',
    amount: amount,
    category: category,
    isExpense: isExpense,
    date: date,
    description: '',
    confirmed: false,
  );
}

void main() {
  test('sums expenses per category for the month, largest first', () {
    final nov = DateTime(2026, 11);
    final txns = [
      _tx('transport', 25000, DateTime(2026, 11, 2)),
      _tx('food & drink', 60000, DateTime(2026, 11, 3)),
      _tx('transport', 50000, DateTime(2026, 11, 9)),
      _tx('shopping', 400000, DateTime(2026, 10, 30)), // other month
      _tx('transfer', 100000, DateTime(2026, 11, 4), isExpense: false),
    ];

    final totals = expenseTotals(txns, nov);

    expect(totals.map((e) => e.key), ['transport', 'food & drink']);
    expect(totals.map((e) => e.value), [75000, 60000]);
  });

  test('is empty when the month has no expenses', () {
    expect(expenseTotals([], DateTime(2026, 11)), isEmpty);
  });
}
