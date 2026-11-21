import 'package:flutter/material.dart';

import '../theme/theme_notifier.dart';
import '../widgets/section_label.dart';

class ThemeSelectionPage extends StatelessWidget {
  const ThemeSelectionPage({super.key});

  Widget _buildSchemeSection() {
    return const Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        SectionLabel('Color scheme'),
        SizedBox(height: 4),
        // Only one scheme for now, more palettes go next to it.
        ListTile(
          contentPadding: EdgeInsets.zero,
          leading: Icon(Icons.color_lens_outlined),
          title: Text('Catppuccin'),
          subtitle: Text('Mocha when dark, Latte when light'),
          trailing: Icon(Icons.check),
        ),
      ],
    );
  }

  Widget _buildModeSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionLabel('Mode'),
        ValueListenableBuilder<ThemeMode>(
          valueListenable: themeNotifier,
          builder: (context, mode, _) {
            final isDark = mode == ThemeMode.dark;
            return SwitchListTile(
              contentPadding: EdgeInsets.zero,
              secondary: Icon(
                isDark ? Icons.dark_mode_outlined : Icons.light_mode_outlined,
              ),
              title: const Text('Dark mode'),
              subtitle: Text(isDark ? 'Mocha' : 'Latte'),
              value: isDark,
              onChanged: (_) => toggleTheme(),
            );
          },
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Themes')),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          _buildSchemeSection(),
          const SizedBox(height: 32),
          const Divider(),
          const SizedBox(height: 16),
          _buildModeSection(),
        ],
      ),
    );
  }
}
