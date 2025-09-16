import type { Round, VoteResult } from '../types';

export const getAllRounds = async (): Promise<Round[]> => {
  const response = await fetch('/api/Voting/getAllRounds');

  return response.json();
};

export const getRound = async (roundId: string): Promise<Round> => {
  const response = await fetch(`/api/Voting/getRoundById?voteId=${roundId}`);

  return response.json();
};

export const getResultsByRoundId = async (
  roundId: string
): Promise<VoteResult[]> => {
  const response = await fetch(
    `/api/Voting/GetResultForVote?number=${roundId}`
  );

  return response.json();
};
