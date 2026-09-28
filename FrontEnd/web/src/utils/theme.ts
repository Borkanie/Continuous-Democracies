import { createTheme, type Theme } from '@mui/material/styles';
import type { Theme as AppTheme } from './context/ThemeContext';

/*
 * Mirrors the :root[data-theme='light'|'dark'] values in
 * routes/styles/App.module.css - keep these in sync if those change.
 *
 * These have to be literal hex, not var(--theme-*) strings: several MUI
 * components (Button, ToggleButton) call theme.alpha()/lighten()/darken()
 * on palette.primary.main and palette.text.primary while building their
 * base styles, and none of those color-math helpers can parse a var()
 * string. Everything in `components.*.styleOverrides` below is our own
 * plain CSS though, so it stays var()-driven and safe.
 */
const TOKENS: Record<AppTheme, Record<string, string>> = {
  light: {
    backgroundPrimary: '#ffffff',
    backgroundSecondary: '#f3f8ff',
    colorPrimary: '#2563eb',
    fontPrimary: '#231f20',
    fontSecondary: '#8a8a8a',
  },
  dark: {
    backgroundPrimary: '#232524',
    backgroundSecondary: '#171717',
    colorPrimary: '#2563eb',
    fontPrimary: '#e5e7eb',
    fontSecondary: '#9ca3af',
  },
};

export const createAppTheme = (mode: AppTheme): Theme => {
  const t = TOKENS[mode];

  return createTheme({
    palette: {
      mode,
      primary: {
        main: t.colorPrimary,
        contrastText: '#ffffff',
      },
      background: {
        default: t.backgroundSecondary,
        paper: t.backgroundPrimary,
      },
      text: {
        primary: t.fontPrimary,
        secondary: t.fontSecondary,
      },
    },
    typography: {
      fontFamily: "'Montserrat', sans-serif",
    },
    shape: {
      borderRadius: 4,
    },
    components: {
      MuiButton: {
        styleOverrides: {
          root: {
            border: 'var(--theme-border-primary)',
            borderRadius: 4,
            textTransform: 'none',
            color: 'var(--theme-font-color-primary)',
            transition: 'all 0.1s',
            '&:hover': {
              border: '1px solid var(--theme-color-primary)',
              backgroundColor: 'var(--theme-color-primary)',
              color: 'var(--background-color-primary)',
            },
          },
          startIcon: {
            color: 'var(--theme-font-color-secondary)',
            transition: 'all 0.1s',
          },
        },
      },
      MuiIconButton: {
        styleOverrides: {
          root: {
            border: 'var(--theme-border-primary)',
            backgroundColor: 'var(--background-color-tertiary)',
            borderRadius: 6,
            color: 'var(--theme-font-color-secondary)',
          },
        },
      },
      MuiToggleButton: {
        styleOverrides: {
          root: {
            position: 'relative',
            zIndex: 1,
            padding: '6px 14px',
            border: 'none',
            borderRadius: 999,
            textTransform: 'none',
            fontSize: 'var(--theme-font-size-s)',
            fontWeight: 600,
            color: 'var(--theme-font-color-secondary)',
            backgroundColor: 'transparent',
            transition: 'color 0.35s ease',
            '&.Mui-selected': {
              color: 'var(--theme-font-color-primary)',
              backgroundColor: 'transparent',
            },
            '&.Mui-selected:hover': {
              backgroundColor: 'transparent',
            },
          },
        },
      },
    },
  });
};
