import { useMatches } from '@tanstack/react-router';

export type Breadcrumb = {
  label: string;
  to: string;
};

/**
 * Hook to retrieve breadcrumbs from the current route matches.
 * @returns An array of breadcrumb objects.
 */
export function useBreadcrumbs(): Breadcrumb[] {
  const matches = useMatches();

  return matches
    .filter((m) => m.staticData.breadcrumb)
    .map((m) => ({
      label: m.staticData.breadcrumb || '',
      to: m.pathname,
    }));
}
