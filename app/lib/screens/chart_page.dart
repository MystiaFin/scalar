import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';

import '../models/category.dart';
import '../models/transaction.dart';
import '../services/transaction_service.dart';
import '../state/app_state.dart';
import '../utils/format.dart';
import '../utils/totals.dart';
import '../widgets/month_dropdown.dart';
import '../widgets/status_view.dart';

class ChartPage extends StatefulWidget {
  const ChartPage({super.key});

  @override
  State<ChartPage> createState() => _ChartPageState();
}

class _ChartPageState extends State<ChartPage> {
  final _service = TransactionService();
  late Future<List<Transaction>> _transactions;

  DateTime get _selectedMonth => monthNotifier.value;

  @override
  void initState() {
    super.initState();
    _transactions = _service.getTransactions();
    refreshNotifier.addListener(_reload);
    monthNotifier.addListener(_onMonthChanged);
  }

  @override
  void dispose() {
    refreshNotifier.removeListener(_reload);
    monthNotifier.removeListener(_onMonthChanged);
    super.dispose();
  }

  void _onMonthChanged() => setState(() {});

  void _reload() => setState(() => _transactions = _service.getTransactions());

  /// Jumps to the transactions tab showing only [category].
  void _showCategory(String category) {
    categoryFilterNotifier.value = {category};
    tabNotifier.value = 0;
  }

  List<Color> _sliceColors() {
    final c = Theme.of(context).colorScheme;
    return [
      c.primary,
      c.secondary,
      c.tertiary,
      c.error,
      c.primaryContainer,
      c.secondaryContainer,
      c.tertiaryContainer,
      c.errorContainer,
    ];
  }

  Color _colorAt(int i) {
    final colors = _sliceColors();
    return colors[i % colors.length];
  }

  Widget _buildContent(List<Transaction> txns) {
    final entries = expenseTotals(txns, _selectedMonth);
    if (entries.isEmpty) {
      return const StatusView(
        icon: Icons.pie_chart_outline,
        message: 'No expenses this month',
      );
    }

    final grandTotal = entries.fold(0.0, (sum, e) => sum + e.value);

    return Column(
      children: [
        const SizedBox(height: 24),
        _buildChart(entries, grandTotal),
        _buildSummaryRow(grandTotal),
        const Divider(height: 1),
        Expanded(child: _buildLegend(entries, grandTotal)),
      ],
    );
  }

  Widget _buildChart(
    List<MapEntry<String, double>> entries,
    double grandTotal,
  ) {
    final surface = Theme.of(context).colorScheme.surface;
    final sections = [
      for (var i = 0; i < entries.length; i++)
        PieChartSectionData(
          value: entries[i].value,
          color: _colorAt(i),
          radius: 90,
          title: '${(entries[i].value / grandTotal * 100).toStringAsFixed(1)}%',
          titleStyle: TextStyle(
            fontSize: 11,
            fontWeight: FontWeight.bold,
            color: surface,
          ),
        ),
    ];

    return SizedBox(
      height: 260,
      child: PieChart(
        PieChartData(
          sections: sections,
          sectionsSpace: 3,
          centerSpaceRadius: 50,
          pieTouchData: PieTouchData(
            touchCallback: (event, response) {
              if (event is! FlTapUpEvent) return;
              final i = response?.touchedSection?.touchedSectionIndex ?? -1;
              if (i < 0 || i >= entries.length) return;
              _showCategory(entries[i].key);
            },
          ),
        ),
      ),
    );
  }

  Widget _buildSummaryRow(double grandTotal) {
    final colors = Theme.of(context).colorScheme;
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: Row(
        children: [
          Text(
            'Total expenses',
            style: TextStyle(
              fontSize: 13,
              color: colors.onSurface.withValues(alpha: 0.6),
            ),
          ),
          const Spacer(),
          Text(
            formatAmount(grandTotal),
            style: TextStyle(fontWeight: FontWeight.bold, color: colors.error),
          ),
        ],
      ),
    );
  }

  Widget _buildLegend(
    List<MapEntry<String, double>> entries,
    double grandTotal,
  ) {
    return ListView.builder(
      itemCount: entries.length,
      itemBuilder: (context, i) => _buildLegendItem(entries[i], grandTotal, i),
    );
  }

  Widget _buildLegendItem(
    MapEntry<String, double> entry,
    double grandTotal,
    int i,
  ) {
    final colors = Theme.of(context).colorScheme;
    final pct = entry.value / grandTotal * 100;

    return ListTile(
      onTap: () => _showCategory(entry.key),
      leading: CircleAvatar(
        radius: 16,
        backgroundColor: _colorAt(i).withValues(alpha: 0.2),
        child: Icon(iconForCategory(entry.key), size: 16, color: _colorAt(i)),
      ),
      title: Text(entry.key, style: const TextStyle(fontSize: 14)),
      trailing: Column(
        crossAxisAlignment: CrossAxisAlignment.end,
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Text(
            formatAmount(entry.value),
            style: const TextStyle(fontWeight: FontWeight.bold),
          ),
          Text(
            '${pct.toStringAsFixed(1)}%',
            style: TextStyle(
              fontSize: 11,
              color: colors.onSurface.withValues(alpha: 0.5),
            ),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: MonthDropdown(
          value: _selectedMonth,
          onChanged: (m) => monthNotifier.value = m,
        ),
      ),
      body: FutureBuilder<List<Transaction>>(
        future: _transactions,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const Center(child: CircularProgressIndicator());
          }
          if (snapshot.hasError) {
            return StatusView(
              icon: Icons.cloud_off_outlined,
              message: '${snapshot.error}',
              onRetry: _reload,
            );
          }
          return _buildContent(snapshot.data!);
        },
      ),
    );
  }
}
