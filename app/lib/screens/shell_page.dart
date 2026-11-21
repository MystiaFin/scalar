import 'package:flutter/material.dart';

import '../state/app_state.dart';
import 'chart_page.dart';
import 'home_page.dart';
import 'settings_page.dart';

/// Bottom navigation between the three tabs. IndexedStack keeps every page
/// alive so switching tabs does not refetch or lose the selected month.
class ShellPage extends StatelessWidget {
  const ShellPage({super.key});

  static const _pages = [HomePage(), ChartPage(), SettingsPage()];

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<int>(
      valueListenable: tabNotifier,
      builder: (context, index, _) => Scaffold(
        body: IndexedStack(index: index, children: _pages),
        bottomNavigationBar: NavigationBar(
          selectedIndex: index,
          onDestinationSelected: (i) => tabNotifier.value = i,
          destinations: const [
            NavigationDestination(
              icon: Icon(Icons.receipt_long_outlined),
              selectedIcon: Icon(Icons.receipt_long),
              label: 'Transactions',
            ),
            NavigationDestination(
              icon: Icon(Icons.pie_chart_outline),
              selectedIcon: Icon(Icons.pie_chart),
              label: 'Overview',
            ),
            NavigationDestination(
              icon: Icon(Icons.settings_outlined),
              selectedIcon: Icon(Icons.settings),
              label: 'Settings',
            ),
          ],
        ),
      ),
    );
  }
}
