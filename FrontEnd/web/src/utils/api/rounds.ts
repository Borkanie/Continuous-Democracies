export type Round = {
  id: string;
  title: string;
  description: string;
  voteDate: Date;
  voteId: number;
  name: string;
};

export const getAllRounds = async (): Promise<Round[]> => {
  const response = await fetch('/api/Voting/getAllRounds');

  return response.json();
};

export const getRound = async (roundId: string): Promise<Round> => {
  const response = await fetch(`/api/Voting/getRoundById?voteId=${roundId}`);

  return response.json();
};
