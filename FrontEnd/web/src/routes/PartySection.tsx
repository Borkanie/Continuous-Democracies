import { useParams } from '@tanstack/react-router';
import { Header } from '../components/section/header';
import sharedStyles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getRound } from '../utils/api/rounds';
import classNames from 'classnames';
import { UiText } from '../components/ui/text/UiText';
import {
  buildPartyPieData,
  groupVotesByPartyForPosition,
} from './RoundSection';
import { useResultsByRoundId } from '../utils/hooks/useResultsByRoundId';
import type { PieChartData } from '../components/chart/PieChart';
import { PoliticiansList } from '../components/politicians-list/PoliticiansList';

const { Div, separator, content, bold, chartTitle } = sharedStyles;

export const PartySection = () => {
  const { roundId, sectionId, partyId } = useParams({ strict: false });

  const { data: roundData } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getRound(queryKey[1] || ''),
  });

  const { data: groupedRoundResults } = useResultsByRoundId(roundId);
  const partyPieData = buildPartyPieData(groupedRoundResults, sectionId);
  const getPartyById = (partyPieData: PieChartData, partyId: string) => {
    return partyPieData.slices.find((party) => party.id === partyId);
  };
  const party = getPartyById(partyPieData, partyId || '');
  const votes = groupVotesByPartyForPosition(groupedRoundResults, sectionId)[
    partyId || ''
  ];

  const title = `${roundData?.title} `;

  return (
    <div className={Div}>
      <Header
        title={title}
        description={roundData?.description}
        extraDetails={{ voteDate: roundData?.voteDate }}
      />
      <div className={separator}></div>

      {/* Party distribution pie + legend for selected section */}
      <div className={content}>
        <UiText
          className={classNames(bold, chartTitle)}
          text={`Votanti per partid - ${party?.label} (${party?.acronym || ''})`}
        />
        <PoliticiansList vote={votes} />
      </div>
    </div>
  );
};
