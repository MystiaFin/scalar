import 'package:flutter/material.dart';

import 'screens/shell_page.dart';
import 'theme/app_theme.dart';
import 'theme/dark/mocha.dart';
import 'theme/light/latte.dart';
import 'theme/theme_notifier.dart';

Future<void> main() async {
  WidgetsFlutterBinding.ensureInitialized();
  await loadSavedTheme();
  runApp(const MyApp());
}

class MyApp extends StatelessWidget {
  const MyApp({super.key});

  @override
  Widget build(BuildContext context) {
    return ValueListenableBuilder<ThemeMode>(
      valueListenable: themeNotifier,
      builder: (context, mode, _) => MaterialApp(
        title: 'Scalar',
        debugShowCheckedModeBanner: false,
        theme: AppTheme.from(const Latte()),
        darkTheme: AppTheme.from(const Mocha()),
        themeMode: mode,
        home: const ShellPage(),
      ),
    );
  }
}
