import { useQuery } from '@tanstack/react-query';
import type { Position, VoteResult, Politician } from '../types';
import { getResultsByRoundId } from '../api/rounds';
import { getAllPoliticians } from '../api/politicians';

export type GroupedVotes = Record<Position, VoteResult[]>;

/**
 * useResultsByRoundId
 *
 * React Query hook that fetches round results and the list of politicians,
 * groups results by position and ensures that all politicians that don't
 * appear in positions 0,1,2 are added to position 3 (absent).
 *
 * - returns grouped results as GroupedVotes
 */
export const useResultsByRoundId = (roundId: string | undefined) => {
  // fetch all politicians (used to infer absents)
  const { data: politicians } = useQuery({
    queryKey: ['allPoliticians'],
    queryFn: () => getAllPoliticians(),
  });

  return useQuery({
    queryKey: ['roundResults', roundId],
    queryFn: () => getResultsByRoundId(roundId!),
    select: (results: VoteResult[]) => {
      const grouped = groupResultsByPosition(results);
      return addAbsentPoliticians(grouped, politicians, results);
    },
    enabled: !!roundId, // query only runs if roundId is truthy
    refetchOnWindowFocus: false,
  });
};

/**
 * Simple grouping helper: produce a GroupedVotes object from the flat results array.
 * Ensures all four position keys exist (0, 1, 2, 3), even if some are empty arrays.
 */
const groupResultsByPosition = (results: VoteResult[]): GroupedVotes => {
  const grouped: GroupedVotes = { 0: [], 1: [], 2: [], 3: [] };
  // push vote into the array corresponding to its position
  results.forEach((vote) => grouped[vote.position].push(vote));
  return grouped;
};

/**
 * Mutates and returns the provided grouped map:
 * - Determines which politicians are already present in positions 0, 1, 2 (and existing 3 entries)
 * - For every politician not present, create a synthetic VoteResult marked as position 3
 *   and push it into grouped[3].
 *
 * This keeps the shape of the data consistent for consumers that expect a VoteResult per politician.
 *
 * Note:
 * - If politicians list is missing, the function returns grouped unchanged.
 * - A sample round (from the results) is used to populate the round field of synthetic VoteResult entries.
 */
const addAbsentPoliticians = (
  grouped: GroupedVotes,
  politicians: Politician[] | undefined,
  results: VoteResult[],
): GroupedVotes => {
  if (!politicians || !Array.isArray(politicians)) {
    return grouped;
  }

  const presentIds = new Set<string>();
  [0, 1, 2].forEach((pos) =>
    grouped[pos as Position].forEach((v) => presentIds.add(v.politician.id)),
  );
  const sampleRound = results[0]?.round;

  // Build a fresh absent list for position 3
  const absents: VoteResult[] = [];
  politicians.forEach((pol) => {
    if (!presentIds.has(pol.id)) {
      absents.push({
        politician: pol,
        position: 3,
        round: sampleRound || {
          id: '',
          title: '',
          description: '',
          voteDate: new Date(),
          voteId: 0,
          name: '',
        },
        id: pol.id,
        name: pol.name,
      });
    }
  });

  grouped[3] = absents;
  return grouped;
};
