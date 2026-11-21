import 'package:flutter_test/flutter_test.dart';

import 'package:scalar_app/utils/format.dart';

void main() {
  group('formatAmount', () {
    test('groups thousands', () {
      expect(formatAmount(0), '0');
      expect(formatAmount(999), '999');
      expect(formatAmount(60000), '60,000');
      expect(formatAmount(1250000), '1,250,000');
    });

    test('keeps two decimals for fractions', () {
      expect(formatAmount(12.5), '12.50');
      expect(formatAmount(1234.56), '1,234.56');
    });

    test('signs expenses and income', () {
      expect(formatSignedAmount(25000, true), '-25,000');
      expect(formatSignedAmount(100000, false), '+100,000');
    });
  });

  test('formats months and dates', () {
    expect(formatMonth(DateTime(2026, 11)), 'November 2026');
    expect(formatDate(DateTime(2026, 1, 5)), 'January 5, 2026');
  });

  test('compares months', () {
    expect(isSameMonth(DateTime(2026, 11, 1), DateTime(2026, 11, 30)), isTrue);
    expect(isSameMonth(DateTime(2026, 11, 1), DateTime(2025, 11, 1)), isFalse);
    expect(monthOf(DateTime(2026, 11, 21, 9, 30)), DateTime(2026, 11));
  });

  test('lists recent months across a year boundary', () {
    final months = recentMonths(now: DateTime(2026, 2, 14), count: 4);
    expect(months, [
      DateTime(2026, 2),
      DateTime(2026, 1),
      DateTime(2025, 12),
      DateTime(2025, 11),
    ]);
  });
}
