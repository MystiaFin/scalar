import 'package:flutter/material.dart';

/// Which tab ShellPage shows.
final tabNotifier = ValueNotifier<int>(0);

/// Categories the transactions list is filtered by, empty means all.
/// ChartPage sets it when a slice is tapped.
final categoryFilterNotifier = ValueNotifier<Set<String>>({});

/// Bumped whenever the data on the backend may have changed (an edit, a new
/// endpoint), so every page that fetched transactions fetches them again.
final refreshNotifier = ValueNotifier<int>(0);

void requestRefresh() => refreshNotifier.value++;
