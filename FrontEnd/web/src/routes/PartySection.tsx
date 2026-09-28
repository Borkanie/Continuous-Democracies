import { useParams } from '@tanstack/react-router';
import { Header } from '../components/section/header';
import sharedStyles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getVotingRoundById } from '../utils/api/rounds';
import classNames from 'classnames';
import { UiText } from '../components/ui/text/UiText';
import { groupVotesByPartyForPosition } from './RoundSection';
import { useResultsByRoundId } from '../utils/hooks/useResultsByRoundId';
import { PoliticiansList } from '../components/politicians-list/PoliticiansList';
import { getVoteStatus } from '../utils/helper';

const { Div, separator, content, bold, chartTitle } = sharedStyles;

export const PartySection = () => {
  const { roundId, sectionId, partyId } = useParams({ strict: false });

  const { data: roundData } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getVotingRoundById(queryKey[1] || ''),
  });

  const { data: groupedRoundResults } = useResultsByRoundId(roundId);
  const votes = groupVotesByPartyForPosition(groupedRoundResults, sectionId)[
    partyId || ''
  ];
  const party = votes?.[0]?.party;

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

      {/* Party distribution pie + legend for selected section */}
      <div className={content}>
        <UiText
          className={classNames(bold, chartTitle)}
          text={`Votanti per partid - ${party?.name} (${party?.acronym || ''})`}
        />
        <PoliticiansList vote={votes || []} />
      </div>
    </div>
  );
};
