import 'package:flutter/material.dart';

import '../services/settings_service.dart';
import '../state/app_state.dart';
import '../widgets/section_label.dart';
import 'theme_selection_page.dart';

class SettingsPage extends StatefulWidget {
  const SettingsPage({super.key});

  @override
  State<SettingsPage> createState() => _SettingsPageState();
}

class _SettingsPageState extends State<SettingsPage> {
  final _settings = SettingsService();
  final _urlController = TextEditingController();
  bool _isDirty = false;

  @override
  void initState() {
    super.initState();
    _loadUrl();
  }

  @override
  void dispose() {
    _urlController.dispose();
    super.dispose();
  }

  Future<void> _loadUrl() async {
    final url = await _settings.getBaseUrl();
    if (mounted) _urlController.text = url;
  }

  Future<void> _saveUrl() async {
    FocusScope.of(context).unfocus();
    await _settings.saveBaseUrl(_urlController.text);
    _urlController.text = await _settings.getBaseUrl();
    requestRefresh();

    if (!mounted) return;
    setState(() => _isDirty = false);
    ScaffoldMessenger.of(context)
        .showSnackBar(const SnackBar(content: Text('Endpoint saved.')));
  }

  Widget _buildApiSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionLabel('API'),
        const SizedBox(height: 10),
        TextField(
          controller: _urlController,
          keyboardType: TextInputType.url,
          autocorrect: false,
          decoration: const InputDecoration(
            labelText: 'Base URL',
            hintText: 'http://192.168.1.10:7225',
            helperText: 'Android emulator: http://10.0.2.2:7225',
            border: OutlineInputBorder(),
          ),
          onChanged: (_) => setState(() => _isDirty = true),
          onSubmitted: (_) => _isDirty ? _saveUrl() : null,
        ),
        const SizedBox(height: 12),
        SizedBox(
          width: double.infinity,
          child: FilledButton(
            onPressed: _isDirty ? _saveUrl : null,
            child: const Text('Save'),
          ),
        ),
      ],
    );
  }

  Widget _buildAppearanceSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        const SectionLabel('Appearance'),
        const SizedBox(height: 4),
        ListTile(
          contentPadding: EdgeInsets.zero,
          leading: const Icon(Icons.palette_outlined),
          title: const Text('Themes'),
          subtitle: const Text('Color scheme and dark mode'),
          trailing: const Icon(Icons.chevron_right),
          onTap: () => Navigator.push(
            context,
            MaterialPageRoute(builder: (_) => const ThemeSelectionPage()),
          ),
        ),
      ],
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Settings')),
      body: ListView(
        padding: const EdgeInsets.all(20),
        children: [
          _buildApiSection(),
          const SizedBox(height: 32),
          const Divider(),
          const SizedBox(height: 24),
          _buildAppearanceSection(),
        ],
      ),
    );
  }
}
