import 'package:shared_preferences/shared_preferences.dart';

class SettingsService {
  static const _key = 'base_url';

  Future<String> getBaseUrl() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString(_key) ?? '';
  }

  /// Saves the url without a trailing slash so paths can be appended as is.
  Future<void> saveBaseUrl(String url) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_key, url.trim().replaceFirst(RegExp(r'/+$'), ''));
  }
}
