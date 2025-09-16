import {
  createRouter,
  createRootRoute,
  createRoute,
  redirect,
} from '@tanstack/react-router';
import { App } from './routes/App';
import { Section } from './routes/section';
import { getAllRounds } from './utils/api/rounds';
import { RoundBreakdown } from './routes/RoundBreakdown';

const rootRoute = createRootRoute({
  component: App,
});

const rootIndexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  beforeLoad: async () => {
    const rounds = await getAllRounds();
    if (rounds && rounds.length > 0) {
      throw redirect({
        to: '/round/$roundId',
        params: { roundId: rounds[0].voteId.toString() },
      });
    }
  },
  component: () => null,
});

const roundBreakdownRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/round/$roundId',
  component: RoundBreakdown,
});

const sectionRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/round/$roundId/section/$sectionId',
  component: Section,
});

const routeTree = rootRoute.addChildren([
  rootIndexRoute,
  roundBreakdownRoute,
  sectionRoute,
]);

export const router = createRouter({ routeTree });

// 👇 tells TS what router you're using
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
