import {
  createRouter,
  createRootRoute,
  createRoute,
} from '@tanstack/react-router';
import { Root } from './routes/root';
import { Index } from './routes';
import { Section } from './routes/section';

const rootRoute = createRootRoute({
  component: Root,
});

const indexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  component: Index,
});

const sectionRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: 'section/$sectionId',
  component: Section,
});

const routeTree = rootRoute.addChildren([indexRoute, sectionRoute]);

export const router = createRouter({ routeTree });

// 👇 tells TS what router you're using
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
