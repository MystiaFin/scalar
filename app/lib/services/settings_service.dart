import 'package:shared_preferences/shared_preferences.dart';

class SettingsService {
  static const _key = 'base_url';

  /// Used until a url is saved in Settings, set with
  /// `--dart-define=SCALAR_API_URL=...` (the justfile does this).
  static const _defaultUrl = String.fromEnvironment('SCALAR_API_URL');

  Future<String> getBaseUrl() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_key) ?? _defaultUrl;
  }

  /// Saves the url without a trailing slash so paths can be appended as is.
  Future<void> saveBaseUrl(String url) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_key, url.trim().replaceFirst(RegExp(r'/+$'), ''));
  }
}
