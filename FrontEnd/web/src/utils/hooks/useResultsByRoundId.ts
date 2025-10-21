import { useQuery } from '@tanstack/react-query';
import type { Position, VoteResult } from '../types';
import { getResultsByRoundId } from '../api/rounds';

export type GroupedVotes = Record<Position, VoteResult[]>;

export const useResultsByRoundId = (roundId: string | undefined) => {
  return useQuery({
    queryKey: ['roundResults', roundId],
    queryFn: () => {
      return getResultsByRoundId(roundId!);
    },
    select: (results: VoteResult[]) => {
      const grouped: GroupedVotes = { 0: [], 1: [], 2: [], 3: [] };
      results.forEach((vote) => grouped[vote.position].push(vote));
      return grouped;
    },
    enabled: !!roundId, // query only runs if roundId is truthy
    staleTime: 5 * 60 * 1000,
    refetchOnWindowFocus: false,
  });
};
