import 'package:flutter/material.dart';

import '../app_colors.dart';

/// Catppuccin Mocha.
class Mocha implements AppColors {
  const Mocha();

  @override
  Color get primary => const Color(0xFFcba6f7); // mauve
  @override
  Color get secondary => const Color(0xFF89b4fa); // blue
  @override
  Color get surface => const Color(0xFF313244); // surface0
  @override
  Color get background => const Color(0xFF11111b); // crust
  @override
  Color get bar => const Color(0xFF1e1e2e); // base
  @override
  Color get error => const Color(0xFFf38ba8); // red
  @override
  Color get text => const Color(0xFFcdd6f4); // text
  @override
  bool get isDark => true;
}
