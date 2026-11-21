import 'package:flutter/material.dart';

/// Same list the backend accepts in `PATCH /api/transactions/:id`.
const categories = [
  'food & drink',
  'transport',
  'utilities',
  'entertainment',
  'shopping',
  'health',
  'transfer',
  'other',
];

IconData iconForCategory(String category) {
  switch (category.toLowerCase()) {
    case 'food & drink':
      return Icons.restaurant_outlined;
    case 'transport':
      return Icons.directions_car_outlined;
    case 'utilities':
      return Icons.bolt_outlined;
    case 'entertainment':
      return Icons.movie_outlined;
    case 'shopping':
      return Icons.shopping_bag_outlined;
    case 'health':
      return Icons.favorite_outline;
    case 'transfer':
      return Icons.swap_horiz_outlined;
    default:
      return Icons.receipt_outlined;
  }
}
