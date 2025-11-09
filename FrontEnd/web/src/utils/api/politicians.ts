import type { Politician } from '../types';

export const getAllPoliticians = async (): Promise<Politician[]> => {
  const response = await fetch(`/api/Politicians/getAllPoliticians?number=500`);

  return response.json();
};
