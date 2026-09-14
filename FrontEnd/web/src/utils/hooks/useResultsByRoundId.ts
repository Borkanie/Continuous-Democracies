import { useQuery } from '@tanstack/react-query';
import type { VoteValue, HydratedVote } from '../types';
import { getVotesByVotingRound } from '../api/rounds';

export type GroupedVotes = Record<VoteValue, HydratedVote[]>;

/**
 * useResultsByRoundId
 *
 * React Query hook that fetches every vote cast in a round and groups it by
 * vote value. The backend already returns an explicit vote (including
 * "Absent") for every politician, so no client-side synthesis is needed.
 */
export const useResultsByRoundId = (roundId: string | undefined) => {
  return useQuery({
    queryKey: ['roundResults', roundId],
    queryFn: () => getVotesByVotingRound(roundId!),
    select: groupVotesByValue,
    enabled: !!roundId,
    refetchOnWindowFocus: false,
  });
};

const groupVotesByValue = (votes: HydratedVote[]): GroupedVotes => {
  const grouped: GroupedVotes = { Yes: [], No: [], Abstain: [], Absent: [] };
  votes.forEach((vote) => grouped[vote.value].push(vote));
  return grouped;
};
