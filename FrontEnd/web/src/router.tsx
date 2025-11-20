import {
  createRouter,
  createRootRoute,
  createRoute,
  redirect,
} from '@tanstack/react-router';
import { App } from './routes/App';
import { RoundSection } from './routes/RoundSection';
import { getAllRounds } from './utils/api/rounds';
import { RoundBreakdown } from './routes/RoundBreakdown';
import { PartySection } from './routes/PartySection';

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
  component: RoundSection,
});

const politicianListRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/round/$roundId/section/$sectionId/party/$partyId',
  component: PartySection,
});

const routeTree = rootRoute.addChildren([
  rootIndexRoute,
  roundBreakdownRoute,
  sectionRoute,
  politicianListRoute,
]);

export const router = createRouter({ routeTree });

// 👇 tells TS what router we're using
declare module '@tanstack/react-router' {
  interface Register {
    router: typeof router;
  }
}
