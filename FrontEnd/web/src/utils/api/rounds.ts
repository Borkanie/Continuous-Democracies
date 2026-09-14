import type { VotingRound, HydratedVote } from '../types';
import { API_BASE_URL } from './config';

export const listVotingRounds = async (): Promise<VotingRound[]> => {
  const response = await fetch(`${API_BASE_URL}/votingRounds`);
  if (!response.ok) {
    throw new Error(`Failed to fetch voting rounds: ${response.status}`);
  }

  return response.json();
};

export const getVotingRoundById = async (
  roundId: string,
): Promise<VotingRound> => {
  const response = await fetch(`${API_BASE_URL}/votingRounds/${roundId}`);
  if (!response.ok) {
    throw new Error(`Failed to fetch voting round ${roundId}: ${response.status}`);
  }

  return response.json();
};

export const getVotesByVotingRound = async (
  roundId: string,
): Promise<HydratedVote[]> => {
  const response = await fetch(
    `${API_BASE_URL}/votingRounds/${roundId}/votes`,
  );
  if (!response.ok) {
    throw new Error(`Failed to fetch votes for round ${roundId}: ${response.status}`);
  }

  return response.json();
};
