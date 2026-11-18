import '../models/transaction.dart';
import 'api_client.dart';

class TransactionService {
  final _client = ApiClient();

  Future<List<Transaction>> getTransactions() async {
    final data = await _client.get('/api/transactions') as List;
    return data.cast<Map<String, dynamic>>().map(Transaction.fromJson).toList();
  }

  /// [confirmed] tells the backend to apply the category to every
  /// transaction from the same merchant and use it as a future example.
  Future<void> updateTransaction(
    int id, {
    required String category,
    required bool isExpense,
    required String description,
    required bool confirmed,
  }) async {
    await _client.patch('/api/transactions/$id', {
      'category': category,
      'is_expense': isExpense,
      'description': description,
      'confirmed': confirmed,
    });
  }
}
