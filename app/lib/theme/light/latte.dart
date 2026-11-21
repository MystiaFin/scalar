import 'package:flutter/material.dart';

import '../app_colors.dart';

/// Catppuccin Latte.
class Latte implements AppColors {
  const Latte();

  @override
  Color get primary => const Color(0xFF8839ef); // mauve
  @override
  Color get secondary => const Color(0xFF1e66f5); // blue
  @override
  Color get surface => const Color(0xFFccd0da); // surface0
  @override
  Color get background => const Color(0xFFeff1f5); // base
  @override
  Color get bar => const Color(0xFFe6e9ef); // mantle
  @override
  Color get error => const Color(0xFFd20f39); // red
  @override
  Color get text => const Color(0xFF4c4f69); // text
  @override
  bool get isDark => false;
}
