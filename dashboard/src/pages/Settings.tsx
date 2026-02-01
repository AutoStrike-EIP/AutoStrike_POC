import { useState, useEffect } from 'react';
import toast from 'react-hot-toast';

// Configuration constants
const DEFAULT_SERVER_URL = import.meta.env.VITE_SERVER_URL || 'https://localhost:8443';
const DEFAULT_HEARTBEAT_INTERVAL = 30;
const DEFAULT_STALE_TIMEOUT = 120;
const STORAGE_KEY = 'autostrike_settings';

interface Settings {
  serverUrl: string;
  safeMode: boolean;
  heartbeatInterval: number;
  staleTimeout: number;
  caCertPath: string;
  serverCertPath: string;
  serverKeyPath: string;
}

const defaultSettings: Settings = {
  serverUrl: DEFAULT_SERVER_URL,
  safeMode: true,
  heartbeatInterval: DEFAULT_HEARTBEAT_INTERVAL,
  staleTimeout: DEFAULT_STALE_TIMEOUT,
  caCertPath: '',
  serverCertPath: '',
  serverKeyPath: '',
};

export default function Settings() {
  const [settings, setSettings] = useState<Settings>(defaultSettings);

  // Load settings from localStorage on mount
  useEffect(() => {
    const saved = localStorage.getItem(STORAGE_KEY);
    if (saved) {
      try {
        const parsed = JSON.parse(saved);
        setSettings({ ...defaultSettings, ...parsed });
      } catch {
        // Use defaults if parsing fails
      }
    }
  }, []);

  const handleSave = () => {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(settings));
    toast.success('Settings saved successfully');
  };

  const updateSetting = <K extends keyof Settings>(key: K, value: Settings[K]) => {
    setSettings(prev => ({ ...prev, [key]: value }));
  };

  return (
    <div>
      <h1 className="text-3xl font-bold mb-8">Settings</h1>

      <div className="max-w-2xl space-y-6">
        {/* Server Configuration */}
        <div className="card">
          <h2 className="text-lg font-semibold mb-4">Server Configuration</h2>
          <div className="space-y-4">
            <div>
              <label htmlFor="server-url" className="block text-sm font-medium text-gray-700 mb-1">
                Server URL
              </label>
              <input
                id="server-url"
                type="text"
                className="input"
                value={settings.serverUrl}
                onChange={(e) => updateSetting('serverUrl', e.target.value)}
              />
            </div>
          </div>
        </div>

        {/* Execution Settings */}
        <div className="card">
          <h2 className="text-lg font-semibold mb-4">Execution Settings</h2>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium">Safe Mode by Default</p>
                <p className="text-sm text-gray-500">
                  Only run safe techniques that don't modify the system
                </p>
              </div>
              <button
                onClick={() => updateSetting('safeMode', !settings.safeMode)}
                className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors ${
                  settings.safeMode ? 'bg-primary-600' : 'bg-gray-200'
                }`}
              >
                <span
                  className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${
                    settings.safeMode ? 'translate-x-6' : 'translate-x-1'
                  }`}
                />
              </button>
            </div>
          </div>
        </div>

        {/* Agent Settings */}
        <div className="card">
          <h2 className="text-lg font-semibold mb-4">Agent Settings</h2>
          <div className="space-y-4">
            <div>
              <label htmlFor="heartbeat-interval" className="block text-sm font-medium text-gray-700 mb-1">
                Heartbeat Interval (seconds)
              </label>
              <input
                id="heartbeat-interval"
                type="number"
                className="input"
                value={settings.heartbeatInterval}
                onChange={(e) => updateSetting('heartbeatInterval', Number.parseInt(e.target.value, 10) || DEFAULT_HEARTBEAT_INTERVAL)}
              />
            </div>
            <div>
              <label htmlFor="stale-timeout" className="block text-sm font-medium text-gray-700 mb-1">
                Stale Agent Timeout (seconds)
              </label>
              <input
                id="stale-timeout"
                type="number"
                className="input"
                value={settings.staleTimeout}
                onChange={(e) => updateSetting('staleTimeout', Number.parseInt(e.target.value, 10) || DEFAULT_STALE_TIMEOUT)}
              />
            </div>
          </div>
        </div>

        {/* TLS Settings */}
        <div className="card">
          <h2 className="text-lg font-semibold mb-4">TLS / mTLS Configuration</h2>
          <div className="space-y-4">
            <div>
              <label htmlFor="ca-cert-path" className="block text-sm font-medium text-gray-700 mb-1">
                CA Certificate Path
              </label>
              <input
                id="ca-cert-path"
                type="text"
                className="input"
                placeholder="/path/to/ca.crt"
                value={settings.caCertPath}
                onChange={(e) => updateSetting('caCertPath', e.target.value)}
              />
            </div>
            <div>
              <label htmlFor="server-cert-path" className="block text-sm font-medium text-gray-700 mb-1">
                Server Certificate Path
              </label>
              <input
                id="server-cert-path"
                type="text"
                className="input"
                placeholder="/path/to/server.crt"
                value={settings.serverCertPath}
                onChange={(e) => updateSetting('serverCertPath', e.target.value)}
              />
            </div>
            <div>
              <label htmlFor="server-key-path" className="block text-sm font-medium text-gray-700 mb-1">
                Server Key Path
              </label>
              <input
                id="server-key-path"
                type="text"
                className="input"
                placeholder="/path/to/server.key"
                value={settings.serverKeyPath}
                onChange={(e) => updateSetting('serverKeyPath', e.target.value)}
              />
            </div>
          </div>
        </div>

        <div className="flex justify-end">
          <button className="btn-primary" onClick={handleSave}>
            Save Settings
          </button>
        </div>
      </div>
    </div>
  );
}
