/// One row from `GET /api/transactions` on the scalar backend.
class Transaction {
  final int id;
  final String merchant;
  final double amount;
  final String category;
  final bool isExpense;
  final DateTime date;
  final String description;
  final bool confirmed;

  const Transaction({
    required this.id,
    required this.merchant,
    required this.amount,
    required this.category,
    required this.isExpense,
    required this.date,
    required this.description,
    required this.confirmed,
  });

  factory Transaction.fromJson(Map<String, dynamic> json) => Transaction(
    id: json['id'] as int,
    merchant: json['merchant'] as String,
    amount: (json['amount'] as num).toDouble(),
    category: json['category'] as String,
    isExpense: json['is_expense'] as bool,
    date: DateTime.parse(json['date'] as String),
    description: json['description'] as String? ?? '',
    confirmed: json['confirmed'] as bool? ?? false,
  );
}
