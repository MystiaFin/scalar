import 'package:flutter/material.dart';

import '../models/category.dart';
import '../models/transaction.dart';
import '../services/transaction_service.dart';

/// Bottom sheet for correcting what the extractor guessed.
/// Pops with `true` once the backend accepted the change.
class EditTransactionSheet extends StatefulWidget {
  final Transaction transaction;

  const EditTransactionSheet({super.key, required this.transaction});

  @override
  State<EditTransactionSheet> createState() => _EditTransactionSheetState();
}

class _EditTransactionSheetState extends State<EditTransactionSheet> {
  final _service = TransactionService();
  late String _category;
  late bool _isExpense;
  late bool _confirmed;
  late final TextEditingController _descController;
  bool _isSaving = false;

  @override
  void initState() {
    super.initState();
    final t = widget.transaction;
    _category = categories.contains(t.category) ? t.category : 'other';
    _isExpense = t.isExpense;
    _confirmed = t.confirmed;
    _descController = TextEditingController(text: t.description);
  }

  @override
  void dispose() {
    _descController.dispose();
    super.dispose();
  }

  Future<void> _submit() async {
    setState(() => _isSaving = true);
    try {
      await _service.updateTransaction(
        widget.transaction.id,
        category: _category,
        isExpense: _isExpense,
        description: _descController.text.trim(),
        confirmed: _confirmed,
      );
      if (mounted) Navigator.pop(context, true);
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context)
            .showSnackBar(SnackBar(content: Text('Error: $e')));
      }
    } finally {
      if (mounted) setState(() => _isSaving = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    final bottomInset = MediaQuery.of(context).viewInsets.bottom;

    return Padding(
      padding: EdgeInsets.fromLTRB(20, 20, 20, bottomInset + 20),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.stretch,
        children: [
          Text(
            'Edit ${widget.transaction.merchant}',
            style: const TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 20),
          DropdownButtonFormField<String>(
            initialValue: _category,
            decoration: const InputDecoration(
              labelText: 'Category',
              border: OutlineInputBorder(),
            ),
            items: [
              for (final c in categories)
                DropdownMenuItem(value: c, child: Text(c)),
            ],
            onChanged: (val) => setState(() => _category = val!),
          ),
          const SizedBox(height: 16),
          SwitchListTile(
            title: const Text('Is an expense?'),
            subtitle: const Text('Turn off if this was money received'),
            value: _isExpense,
            onChanged: (val) => setState(() => _isExpense = val),
            contentPadding: EdgeInsets.zero,
          ),
          CheckboxListTile(
            title: const Text('Save as default for this merchant'),
            subtitle: const Text(
              'Updates past entries and trains the AI for future ones',
            ),
            value: _confirmed,
            onChanged: (val) => setState(() => _confirmed = val!),
            contentPadding: EdgeInsets.zero,
            controlAffinity: ListTileControlAffinity.leading,
          ),
          const SizedBox(height: 8),
          TextField(
            controller: _descController,
            decoration: const InputDecoration(
              labelText: 'Description / note',
              border: OutlineInputBorder(),
            ),
            maxLines: 2,
          ),
          const SizedBox(height: 24),
          FilledButton(
            onPressed: _isSaving ? null : _submit,
            child: _isSaving
                ? const SizedBox.square(
                    dimension: 20,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Text('Save changes'),
          ),
        ],
      ),
    );
  }
}
