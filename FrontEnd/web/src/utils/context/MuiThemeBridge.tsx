import { useMemo } from 'react';
import { ThemeProvider as MuiThemeProvider } from '@mui/material/styles';
import { useTheme } from './ThemeContext';
import { createAppTheme } from '../theme';

export function MuiThemeBridge({ children }: { children: React.ReactNode }) {
  const { theme } = useTheme();
  const muiTheme = useMemo(() => createAppTheme(theme), [theme]);

  return <MuiThemeProvider theme={muiTheme}>{children}</MuiThemeProvider>;
}
