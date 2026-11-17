import 'package:flutter/material.dart';

import 'app_colors.dart';

class AppTheme {
  /// Builds a Material 3 theme from a palette. Every color role is set by
  /// hand so Flutter never falls back to its own seed colors or tints.
  static ThemeData from(AppColors colors) {
    final surface = colors.surface;
    final text = colors.text;
    final primary = colors.primary;
    final secondary = colors.secondary;
    final error = colors.error;
    final tertiary = Color.lerp(primary, secondary, 0.5)!;

    Color shift(Color c, [double amount = 0.85]) =>
        Color.lerp(c, surface, amount)!;

    final scheme = (colors.isDark ? ColorScheme.dark : ColorScheme.light)(
      primary: primary,
      onPrimary: surface,
      primaryContainer: shift(primary),
      onPrimaryContainer: primary,
      secondary: secondary,
      onSecondary: surface,
      secondaryContainer: shift(secondary),
      onSecondaryContainer: secondary,
      tertiary: tertiary,
      onTertiary: surface,
      tertiaryContainer: shift(tertiary),
      onTertiaryContainer: tertiary,
      error: error,
      onError: surface,
      errorContainer: shift(error),
      onErrorContainer: error,
      surface: surface,
      onSurface: text,
      surfaceContainerHighest: shift(text, colors.isDark ? 0.88 : 0.94),
      onSurfaceVariant: text.withValues(alpha: 0.7),
      outline: text.withValues(alpha: 0.30),
      outlineVariant: text.withValues(alpha: 0.12),
      shadow: Colors.black,
      scrim: Colors.black,
      inverseSurface: text,
      onInverseSurface: colors.bar,
      inversePrimary: primary,
      surfaceTint: Colors.transparent,
    );

    return ThemeData(
      useMaterial3: true,
      fontFamily: 'Inter',
      colorScheme: scheme,
      scaffoldBackgroundColor: colors.background,
      applyElevationOverlayColor: false,
      appBarTheme: AppBarTheme(
        backgroundColor: colors.bar,
        foregroundColor: text,
        surfaceTintColor: Colors.transparent,
      ),
      navigationBarTheme: NavigationBarThemeData(
        backgroundColor: colors.bar,
        indicatorColor: shift(primary),
        surfaceTintColor: Colors.transparent,
      ),
      bottomSheetTheme: BottomSheetThemeData(
        backgroundColor: colors.bar,
        surfaceTintColor: Colors.transparent,
      ),
    );
  }
}
