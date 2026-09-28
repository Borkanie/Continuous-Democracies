import {
  createRouter,
  createRootRoute,
  createRoute,
  redirect,
} from '@tanstack/react-router';
import { App } from './routes/App';
import { RoundSection } from './routes/RoundSection';
import { listVotingRounds } from './utils/api/rounds';
import { RoundBreakdown } from './routes/RoundBreakdown';
import { PartySection } from './routes/PartySection';
import type { Chamber } from './utils/types';

declare module '@tanstack/react-router' {
  interface StaticDataRouteOption {
    breadcrumb?: string;
  }
}

const rootRoute = createRootRoute({
  component: App,
});

const rootIndexRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/',
  beforeLoad: async () => {
    // Runs outside React context, so read the chamber directly from
    // localStorage (mirrors how ChamberContext/ThemeContext bootstrap). The
    // backend has no chamber filter, so filter the full list client-side.
    const chamber =
      (localStorage.getItem('chamber') as Chamber) || 'parliament';
    const rounds = await listVotingRounds();
    const roundsInChamber = rounds?.filter((round) => round.chamber === chamber);
    if (roundsInChamber && roundsInChamber.length > 0) {
      throw redirect({
        to: '/round/$roundId',
        params: { roundId: roundsInChamber[0].id.toString() },
      });
    }
  },
  component: () => null,
});

const roundBreakdownRoute = createRoute({
  getParentRoute: () => rootRoute,
  path: '/round/$roundId',
  component: RoundBreakdown,
  staticData: {
    breadcrumb: 'Distributia voturilor',
  },
});

const sectionRoute = createRoute({
  getParentRoute: () => roundBreakdownRoute,
  path: '/section/$sectionId',
  component: RoundSection,
  staticData: {
    breadcrumb: 'Impartirea pe partide',
  },
});

const politicianListRoute = createRoute({
  getParentRoute: () => sectionRoute,
  path: '/party/$partyId',
  component: PartySection,
  staticData: {
    breadcrumb: 'Votanti per partid',
  },
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
