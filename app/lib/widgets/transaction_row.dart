import 'package:flutter/material.dart';

import '../models/category.dart';
import '../models/transaction.dart';
import '../utils/format.dart';
import 'edit_transaction_sheet.dart';

class TransactionRow extends StatelessWidget {
  final Transaction transaction;

  /// Called after the transaction was saved, so the list can reload.
  final VoidCallback onChanged;

  const TransactionRow({
    super.key,
    required this.transaction,
    required this.onChanged,
  });

  Future<void> _openEditSheet(BuildContext context) async {
    final saved = await showModalBottomSheet<bool>(
      context: context,
      isScrollControlled: true,
      builder: (_) => EditTransactionSheet(transaction: transaction),
    );
    if (saved ?? false) onChanged();
  }

  Widget _buildIcon(ColorScheme colors) {
    return Container(
      padding: const EdgeInsets.all(8),
      decoration: BoxDecoration(
        color: colors.primary.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(8),
      ),
      child: Icon(
        iconForCategory(transaction.category),
        color: colors.primary,
        size: 20,
      ),
    );
  }

  Widget _buildLabels(ColorScheme colors) {
    final subtitle = transaction.description.isEmpty
        ? transaction.category
        : '${transaction.category} · ${transaction.description}';

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Row(
          children: [
            Flexible(
              child: Text(
                transaction.merchant,
                overflow: TextOverflow.ellipsis,
                style: TextStyle(
                  color: colors.onSurface,
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
            if (transaction.confirmed) ...[
              const SizedBox(width: 4),
              Icon(Icons.verified_outlined, size: 14, color: colors.secondary),
            ],
          ],
        ),
        const SizedBox(height: 2),
        Text(
          subtitle,
          maxLines: 1,
          overflow: TextOverflow.ellipsis,
          style: TextStyle(
            color: colors.onSurface.withValues(alpha: 0.5),
            fontSize: 12,
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    final colors = Theme.of(context).colorScheme;

    return InkWell(
      onTap: () => _openEditSheet(context),
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
        child: Row(
          children: [
            _buildIcon(colors),
            const SizedBox(width: 12),
            Expanded(child: _buildLabels(colors)),
            const SizedBox(width: 12),
            Text(
              formatSignedAmount(transaction.amount, transaction.isExpense),
              style: TextStyle(
                color: transaction.isExpense ? colors.error : colors.primary,
                fontWeight: FontWeight.bold,
                fontSize: 15,
              ),
            ),
          ],
        ),
      ),
    );
  }
}
