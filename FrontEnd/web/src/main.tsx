import './reset.css';
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { RouterProvider } from '@tanstack/react-router';
import { router } from './router.tsx';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { DrawerProvider } from './utils/context/DrawerContext';
import { ThemeProvider } from './utils/context/ThemeContext.tsx';
import { ChamberProvider } from './utils/context/ChamberContext.tsx';
import { MuiThemeBridge } from './utils/context/MuiThemeBridge.tsx';

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      refetchOnWindowFocus: false,
      staleTime: 5 * 60 * 1000, // cache data for 5 minutes
    },
  },
});

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <ThemeProvider>
        <MuiThemeBridge>
          <ChamberProvider>
            <DrawerProvider>
              <RouterProvider router={router} />
            </DrawerProvider>
          </ChamberProvider>
        </MuiThemeBridge>
      </ThemeProvider>
    </QueryClientProvider>
  </StrictMode>,
);
