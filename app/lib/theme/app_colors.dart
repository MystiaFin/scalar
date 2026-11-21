import 'package:flutter/material.dart';

/// The handful of colors a palette has to provide. Everything else in the
/// Material color scheme is derived from these in [AppTheme.from].
abstract class AppColors {
  Color get primary;
  Color get secondary;
  Color get surface;
  Color get background;
  Color get bar;
  Color get error;
  Color get text;
  bool get isDark;
}
