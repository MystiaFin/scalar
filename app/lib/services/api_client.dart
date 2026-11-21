import 'dart:async';
import 'dart:convert';

import 'package:http/http.dart' as http;

import 'settings_service.dart';

class ApiException implements Exception {
  final String message;
  const ApiException(this.message);

  @override
  String toString() => message;
}

class ApiClient {
  static const _timeout = Duration(seconds: 10);
  final _settings = SettingsService();

  Future<Uri> _uri(String path) async {
    final base = await _settings.getBaseUrl();
    if (base.isEmpty) {
      throw const ApiException('No endpoint set. Configure it in settings.');
    }
    return Uri.parse('$base$path');
  }

  Future<http.Response> _send(
    String method,
    Future<http.Response> request,
  ) async {
    try {
      return await request.timeout(_timeout);
    } on TimeoutException {
      throw ApiException('$method timed out, is the backend running?');
    } on http.ClientException catch (e) {
      throw ApiException('Could not reach the backend: ${e.message}');
    }
  }

  dynamic _decode(String method, String path, http.Response response) {
    if (response.statusCode == 200 || response.statusCode == 204) {
      return response.body.isEmpty ? null : jsonDecode(response.body);
    }
    throw ApiException('$method $path failed: ${response.statusCode}');
  }

  Future<dynamic> get(String path) async {
    final uri = await _uri(path);
    final response = await _send('GET', http.get(uri));
    return _decode('GET', path, response);
  }

  Future<dynamic> patch(String path, Map<String, dynamic> body) async {
    final uri = await _uri(path);
    final response = await _send(
      'PATCH',
      http.patch(
        uri,
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode(body),
      ),
    );
    return _decode('PATCH', path, response);
  }
}
