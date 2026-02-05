import { createContext, useContext, useEffect, useState, useMemo, ReactNode, useCallback } from 'react';

/**
 * Available theme options
 */
export type Theme = 'light' | 'dark' | 'system';

/**
 * Theme context value interface
 */
interface ThemeContextType {
  /** Current theme setting */
  readonly theme: Theme;
  /** Resolved theme (light or dark based on system preference if theme is 'system') */
  readonly resolvedTheme: 'light' | 'dark';
  /** Set theme preference */
  readonly setTheme: (theme: Theme) => void;
  /** Toggle between light and dark */
  readonly toggleTheme: () => void;
}

const STORAGE_KEY = 'autostrike_theme';

const ThemeContext = createContext<ThemeContextType | null>(null);

/**
 * Get system color scheme preference
 */
function getSystemTheme(): 'light' | 'dark' {
  if (globalThis.window === undefined) return 'light';
  return globalThis.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light';
}

/**
 * Get stored theme from localStorage
 */
function getStoredTheme(): Theme {
  if (globalThis.window === undefined) return 'system';
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored === 'light' || stored === 'dark' || stored === 'system') {
    return stored;
  }
  return 'system';
}

interface ThemeProviderProps {
  readonly children: ReactNode;
}

/**
 * Theme provider component that manages dark/light mode
 * with persistence and system preference detection.
 */
export function ThemeProvider({ children }: ThemeProviderProps) {
  const [theme, setTheme] = useState<Theme>(getStoredTheme);
  const [resolvedTheme, setResolvedTheme] = useState<'light' | 'dark'>(() => {
    const stored = getStoredTheme();
    return stored === 'system' ? getSystemTheme() : stored;
  });

  // Apply theme to document
  const applyTheme = useCallback((t: Theme) => {
    const resolved = t === 'system' ? getSystemTheme() : t;
    setResolvedTheme(resolved);

    const root = document.documentElement;
    if (resolved === 'dark') {
      root.classList.add('dark');
    } else {
      root.classList.remove('dark');
    }
  }, []);

  // Initialize and watch for system theme changes
  useEffect(() => {
    applyTheme(theme);

    // Listen for system theme changes when in system mode
    if (theme === 'system') {
      const mediaQuery = globalThis.matchMedia('(prefers-color-scheme: dark)');
      const handler = () => applyTheme('system');
      mediaQuery.addEventListener('change', handler);
      return () => mediaQuery.removeEventListener('change', handler);
    }
  }, [theme, applyTheme]);

  const handleSetTheme = useCallback((newTheme: Theme) => {
    setTheme(newTheme);
    localStorage.setItem(STORAGE_KEY, newTheme);
    applyTheme(newTheme);
  }, [applyTheme, setTheme]);

  const toggleTheme = useCallback(() => {
    handleSetTheme(resolvedTheme === 'light' ? 'dark' : 'light');
  }, [resolvedTheme, handleSetTheme]);

  const value = useMemo(() => ({
    theme,
    resolvedTheme,
    setTheme: handleSetTheme,
    toggleTheme,
  }), [theme, resolvedTheme, handleSetTheme, toggleTheme]);

  return (
    <ThemeContext.Provider value={value}>
      {children}
    </ThemeContext.Provider>
  );
}

/**
 * Hook to access theme context
 */
// eslint-disable-next-line react-refresh/only-export-components
export function useTheme(): ThemeContextType {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme must be used within a ThemeProvider');
  }
  return context;
}
