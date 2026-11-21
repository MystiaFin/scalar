import 'package:flutter/material.dart';

import '../utils/format.dart';

/// App bar title that picks one of the last twelve months.
class MonthDropdown extends StatelessWidget {
  final DateTime value;
  final ValueChanged<DateTime> onChanged;

  const MonthDropdown({
    super.key,
    required this.value,
    required this.onChanged,
  });

  @override
  Widget build(BuildContext context) {
    return DropdownButton<DateTime>(
      value: value,
      underline: const SizedBox(),
      icon: const Icon(Icons.expand_more),
      style: Theme.of(context).textTheme.titleLarge,
      items: [
        for (final m in recentMonths())
          DropdownMenuItem(value: m, child: Text(formatMonth(m))),
      ],
      onChanged: (val) {
        if (val != null) onChanged(val);
      },
    );
  }
}
