import 'package:flutter_test/flutter_test.dart';

import 'package:scalar_app/models/transaction.dart';

void main() {
  test('parses a row from the backend', () {
    final tx = Transaction.fromJson({
      'id': 3,
      'email_hash': 'seed-03',
      'merchant': 'Indomaret',
      'amount': 85000,
      'category': 'food & drink',
      'is_expense': true,
      'date': '2026-11-19',
      'description': 'groceries',
      'confirmed': false,
      'updated_at': '2026-11-19 10:00:00',
    });

    expect(tx.id, 3);
    expect(tx.merchant, 'Indomaret');
    expect(tx.amount, 85000.0);
    expect(tx.category, 'food & drink');
    expect(tx.isExpense, isTrue);
    expect(tx.date, DateTime(2026, 11, 19));
    expect(tx.description, 'groceries');
    expect(tx.confirmed, isFalse);
  });

  test('falls back when optional fields are missing', () {
    final tx = Transaction.fromJson({
      'id': 1,
      'merchant': 'GoPay',
      'amount': 12.5,
      'category': 'transfer',
      'is_expense': false,
      'date': '2026-11-01',
    });

    expect(tx.amount, 12.5);
    expect(tx.description, '');
    expect(tx.confirmed, isFalse);
  });
}
