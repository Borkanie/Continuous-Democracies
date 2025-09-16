type PartyColor = {
  r: number;
  g: number;
  b: number;
  a: number;
  isKnownColor: boolean;
  isEmpty: boolean;
  isNamedColor: boolean;
  isSystemColor: boolean;
  name: string;
};

type Party = {
  acronym: string;
  logoUrl: string | null;
  color: PartyColor;
  active: boolean;
  id: string;
  name: string;
};

type Politician = {
  gender: number;
  imageUrl: string | null;
  party: Party;
  active: boolean;
  workLocation: number;
  id: string;
  name: string;
};

export type Round = {
  id: string;
  title: string;
  description: string;
  voteDate: Date;
  voteId: number;
  name: string;
};

export type VoteResult = {
  politician: Politician;
  position: number;
  round: Round;
  id: string;
  name: string;
};
