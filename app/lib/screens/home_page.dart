import 'package:flutter/material.dart';

import '../models/category.dart';
import '../models/transaction.dart';
import '../services/transaction_service.dart';
import '../state/app_state.dart';
import '../utils/format.dart';
import '../widgets/month_dropdown.dart';
import '../widgets/status_view.dart';
import '../widgets/transaction_row.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _service = TransactionService();
  late Future<List<Transaction>> _transactions;
  DateTime _selectedMonth = monthOf(DateTime.now());

  @override
  void initState() {
    super.initState();
    _transactions = _service.getTransactions();
    categoryFilterNotifier.addListener(_onFilterChanged);
    refreshNotifier.addListener(_reload);
  }

  @override
  void dispose() {
    categoryFilterNotifier.removeListener(_onFilterChanged);
    refreshNotifier.removeListener(_reload);
    super.dispose();
  }

  void _onFilterChanged() => setState(() {});

  void _reload() => setState(() => _transactions = _service.getTransactions());

  Future<void> _pullToRefresh() async {
    _reload();
    await _transactions.catchError((_) => <Transaction>[]);
  }

  List<Transaction> _applyFilters(List<Transaction> txns) {
    final active = categoryFilterNotifier.value;
    return txns.where((t) {
      final matchesMonth = isSameMonth(t.date, _selectedMonth);
      final matchesCategory =
          active.isEmpty || active.contains(t.category.toLowerCase());
      return matchesMonth && matchesCategory;
    }).toList();
  }

  Map<String, List<Transaction>> _groupByDate(List<Transaction> txns) {
    final grouped = <String, List<Transaction>>{};
    for (final t in txns) {
      grouped.putIfAbsent(formatDate(t.date), () => []).add(t);
    }
    return grouped;
  }

  void _openCategoryFilter() {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (context) => ValueListenableBuilder<Set<String>>(
        valueListenable: categoryFilterNotifier,
        builder: (context, active, _) => Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(16, 16, 8, 0),
              child: Row(
                children: [
                  const Text(
                    'Filter by category',
                    style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                  ),
                  const Spacer(),
                  TextButton(
                    onPressed: active.isEmpty
                        ? null
                        : () => categoryFilterNotifier.value = {},
                    child: const Text('Clear all'),
                  ),
                ],
              ),
            ),
            const Divider(),
            Flexible(
              child: ListView(
                shrinkWrap: true,
                children: [
                  for (final cat in categories)
                    CheckboxListTile(
                      title: Text(cat),
                      secondary: Icon(iconForCategory(cat)),
                      value: active.contains(cat),
                      controlAffinity: ListTileControlAffinity.leading,
                      onChanged: (checked) {
                        final updated = {...active};
                        checked! ? updated.add(cat) : updated.remove(cat);
                        categoryFilterNotifier.value = updated;
                      },
                    ),
                ],
              ),
            ),
            const SizedBox(height: 16),
          ],
        ),
      ),
    );
  }

  Widget _buildDateHeader(String label) {
    final colors = Theme.of(context).colorScheme;
    return Container(
      width: double.infinity,
      color: colors.surface,
      padding: const EdgeInsets.fromLTRB(16, 10, 16, 10),
      child: Text(
        label,
        style: TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.w600,
          letterSpacing: 0.5,
          color: colors.primary,
        ),
      ),
    );
  }

  Widget _buildList(List<Transaction> txns) {
    final grouped = _groupByDate(txns);
    return ListView(
      physics: const AlwaysScrollableScrollPhysics(),
      children: [
        for (final entry in grouped.entries) ...[
          _buildDateHeader(entry.key),
          for (final t in entry.value)
            TransactionRow(transaction: t, onChanged: requestRefresh),
        ],
      ],
    );
  }

  /// Keeps a non-list state scrollable so pull to refresh still works.
  Widget _buildScrollable(Widget child) {
    return LayoutBuilder(
      builder: (context, constraints) => SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        child: SizedBox(height: constraints.maxHeight, child: child),
      ),
    );
  }

  Widget _buildBody() {
    return FutureBuilder<List<Transaction>>(
      future: _transactions,
      builder: (context, snapshot) {
        if (snapshot.connectionState == ConnectionState.waiting) {
          return const Center(child: CircularProgressIndicator());
        }
        if (snapshot.hasError) {
          return _buildScrollable(
            StatusView(
              icon: Icons.cloud_off_outlined,
              message: '${snapshot.error}',
              onRetry: _reload,
            ),
          );
        }

        final filtered = _applyFilters(snapshot.data!);
        if (filtered.isEmpty) {
          return _buildScrollable(
            const StatusView(
              icon: Icons.receipt_long_outlined,
              message: 'No transactions found',
            ),
          );
        }
        return _buildList(filtered);
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    final active = categoryFilterNotifier.value;

    return Scaffold(
      appBar: AppBar(
        title: MonthDropdown(
          value: _selectedMonth,
          onChanged: (m) => setState(() => _selectedMonth = m),
        ),
        actions: [
          IconButton(
            tooltip: 'Filter by category',
            icon: Badge(
              isLabelVisible: active.isNotEmpty,
              label: Text('${active.length}'),
              child: const Icon(Icons.filter_list),
            ),
            onPressed: _openCategoryFilter,
          ),
        ],
      ),
      body: RefreshIndicator(onRefresh: _pullToRefresh, child: _buildBody()),
    );
  }
}
