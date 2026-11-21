import 'package:flutter/material.dart';

/// Centered icon and message, used for empty lists and failed requests.
class StatusView extends StatelessWidget {
  final IconData icon;
  final String message;
  final VoidCallback? onRetry;

  const StatusView({
    super.key,
    required this.icon,
    required this.message,
    this.onRetry,
  });

  @override
  Widget build(BuildContext context) {
    final colors = Theme.of(context).colorScheme;
    return Center(
      child: Padding(
        padding: const EdgeInsets.all(24),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Icon(
              icon,
              size: 48,
              color: colors.onSurface.withValues(alpha: 0.25),
            ),
            const SizedBox(height: 12),
            Text(
              message,
              textAlign: TextAlign.center,
              style: TextStyle(color: colors.onSurface.withValues(alpha: 0.45)),
            ),
            if (onRetry != null) ...[
              const SizedBox(height: 16),
              OutlinedButton.icon(
                onPressed: onRetry,
                icon: const Icon(Icons.refresh),
                label: const Text('Retry'),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
