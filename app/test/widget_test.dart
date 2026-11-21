import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:shared_preferences/shared_preferences.dart';

import 'package:scalar_app/main.dart';
import 'package:scalar_app/models/transaction.dart';
import 'package:scalar_app/widgets/transaction_row.dart';

void main() {
  setUp(() => SharedPreferences.setMockInitialValues({}));

  testWidgets('asks for an endpoint when none is saved', (tester) async {
    await tester.pumpWidget(const MyApp());
    await tester.pumpAndSettle();

    expect(find.text('Transactions'), findsOneWidget);
    expect(find.text('Overview'), findsOneWidget);
    expect(
      find.text('No endpoint set. Configure it in settings.'),
      findsWidgets,
    );
  });

  testWidgets('settings tab shows the base url field', (tester) async {
    await tester.pumpWidget(const MyApp());
    await tester.pumpAndSettle();

    await tester.tap(find.byIcon(Icons.settings_outlined));
    await tester.pumpAndSettle();

    expect(find.widgetWithText(TextField, 'Base URL'), findsOneWidget);
    expect(find.text('Themes'), findsOneWidget);
  });

  testWidgets('transaction row shows merchant, category and amount', (
    tester,
  ) async {
    final tx = Transaction(
      id: 1,
      merchant: 'Starbucks',
      amount: 60000,
      category: 'food & drink',
      isExpense: true,
      date: DateTime(2026, 11, 21),
      description: '',
      confirmed: false,
    );

    await tester.pumpWidget(
      MaterialApp(
        home: Scaffold(
          body: TransactionRow(transaction: tx, onChanged: () {}),
        ),
      ),
    );

    expect(find.text('Starbucks'), findsOneWidget);
    expect(find.text('food & drink'), findsOneWidget);
    expect(find.text('-60,000'), findsOneWidget);
    expect(find.byIcon(Icons.restaurant_outlined), findsOneWidget);
  });
}
