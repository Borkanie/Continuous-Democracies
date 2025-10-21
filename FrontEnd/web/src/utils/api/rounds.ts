import type { Round, VoteResult } from '../types';

export const getAllRounds = async (searchTerm?: string): Promise<Round[]> => {
  const params = new URLSearchParams();

  if (searchTerm) {
    params.set('keywords', encodeURIComponent(searchTerm));
  }

  const url = `/api/Voting/getAllRounds${
    params.toString() ? `?${params.toString()}` : ''
  }`;
  const response = await fetch(url);

  return response.json();
};

export const getRound = async (roundId: string): Promise<Round> => {
  const response = await fetch(
    `/api/Voting/getRoundById?voteNumber=${roundId}`
  );

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
