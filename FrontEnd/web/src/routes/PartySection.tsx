import { useParams } from '@tanstack/react-router';
import { Header } from '../components/section/header';
import sharedStyles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getVotingRoundById } from '../utils/api/rounds';
import classNames from 'classnames';
import { UiText } from '../components/ui/text/UiText';
import { useMemo } from 'react';
import { useResultsByRoundId } from '../utils/hooks/useResultsByRoundId';
import { PoliticiansList } from '../components/politicians-list/PoliticiansList';
import { PartyIcon } from '../components/ui/party-icon/PartyIcon';
import { getVoteStatus } from '../utils/helper';
import type { VoteValue } from '../utils/types';

const {
  Div,
  separator,
  content,
  bold,
  chartTitle,
  watermarkWrapper,
  watermark,
} = sharedStyles;

const VOTE_VALUES: VoteValue[] = ['Yes', 'No', 'Abstain', 'Absent'];

export const PartySection = () => {
  const { roundId, partyId } = useParams({ strict: false });

  const { data: roundData } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getVotingRoundById(queryKey[1] || ''),
  });

  const { data: groupedRoundResults } = useResultsByRoundId(roundId);
  // Every member of the party, regardless of how they voted (Yes -> Absent).
  const votes = useMemo(
    () =>
      VOTE_VALUES.flatMap((value) => groupedRoundResults?.[value] ?? []).filter(
        (v) => v.party?.id === partyId,
      ),
    [groupedRoundResults, partyId],
  );
  const party = votes[0]?.party;

  const title = `${roundData?.title} `;

  return (
    <div className={Div}>
      <Header
        title={title}
        description={roundData?.description}
        extraDetails={{ voteDate: roundData?.voteDate }}
        status={
          groupedRoundResults && roundData
            ? getVoteStatus(groupedRoundResults, roundData.majorityType, roundData.chamber)
            : undefined
        }
        majorityType={roundData?.majorityType}
      />
      <div className={separator}></div>

      {/* All members of the selected party with their votes */}
      <div className={content}>
        <UiText
          className={classNames(bold, chartTitle)}
          text={`Votanti per partid - ${party?.name} (${party?.acronym || ''})`}
        />
        <div className={watermarkWrapper}>
          <div className={watermark}>
            <PartyIcon party={party} size={220} />
          </div>
          <PoliticiansList vote={votes} showPartyIcon={false} />
        </div>
      </div>
    </div>
  );
};
