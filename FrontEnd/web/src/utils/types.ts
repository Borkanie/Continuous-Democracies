export type PartyColor = string;

export type Party = {
  acronym: string;
  logoUrl: string | null;
  color: PartyColor;
  active: boolean;
  id: string;
  name: string;
};

export type Politician = {
  gender: number;
  imageUrl: string | null;
  partyId: string;
  active: boolean;
  workLocation: number;
  id: string;
  name: string;
};

export type MajorityType = 'simple' | 'absolute' | 'qualified';

export type Chamber = 'parliament' | 'senate';

export type VotingRound = {
  id: number;
  title: string;
  description: string;
  voteDate: Date;
  majorityType: MajorityType;
  chamber: Chamber;
};

export type VoteValue = 'Yes' | 'No' | 'Abstain' | 'Absent';

export type HydratedVote = {
  politicianId: string;
  partyId: string;
  value: VoteValue;
  politician: Politician;
  party: Party;
};
