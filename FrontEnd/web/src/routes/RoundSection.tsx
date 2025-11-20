import { useNavigate, useParams, useRouter } from '@tanstack/react-router';
import sharedStyles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getRound } from '../utils/api/rounds';
import { Spinner } from '../components/spinner/Spinner';
import { useResultsByRoundId } from '../utils/hooks/useResultsByRoundId';
import type { PartyColor, Position, VoteResult } from '../utils/types';
import { PieChart, type PieChartData } from '../components/chart/PieChart';
import { Legend } from '../components/legend/Legend';
import classNames from 'classnames';
import { Header } from '../components/section/header';
import { SECTION_LABELS } from '../utils/constants';

const { Div, separator, content, chartContainer, bold, chartTitle } =
  sharedStyles;

export const RoundSection = () => {
  const { roundId, sectionId } = useParams({ strict: false });
  const router = useRouter();
  const navigate = useNavigate();

  const { data: roundData, isFetching } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getRound(queryKey[1] || ''),
    staleTime: 5 * 60 * 1000, // cache data for 5 minutes
  });
  const { data: groupedRoundResults } = useResultsByRoundId(roundId);

  const partyPieData = buildPartyPieData(groupedRoundResults, sectionId);

  const title = `${roundData?.title} - ${
    SECTION_LABELS[Number(sectionId) as Position] || ''
  }`;

  return (
    <div className={Div}>
      {isFetching ? (
        <Spinner />
      ) : (
        <>
          <Header
            title={title}
            onBack={() => router.history.back()}
            status={'ACTIV'}
            extraDetails={{ voteDate: roundData?.voteDate }}
          />
          <div className={separator}></div>

          {/* Party distribution pie + legend for selected section */}
          <div className={content}>
            <div className={chartContainer}>
              <div>
                <p className={classNames(bold, chartTitle)}>
                  Distributia pe partide
                </p>
                <PieChart
                  data={partyPieData}
                  onSliceClick={(id) => navigate({ to: `party/${id}` })}
                />
              </div>
              <div>
                <Legend slices={partyPieData.slices} />
              </div>
            </div>
          </div>
        </>
      )}
    </div>
  );
};

const groupVotesByPartyForPosition = (
  grouped: Record<Position, VoteResult[]> | undefined,
  sectionId?: string
): Record<string, VoteResult[]> => {
  if (!grouped || !sectionId) {
    return {};
  }

  const pos = Number(sectionId);

  if (!Number.isInteger(pos)) {
    return {};
  }

  const votes = grouped[pos as Position] ?? [];
  return votes.reduce<Record<string, VoteResult[]>>((acc, vote) => {
    const partyId = vote.politician?.party?.id ?? 'unknown';

    if (!acc[partyId]) {
      acc[partyId] = [];
    }

    acc[partyId].push(vote);
    return acc;
  }, {});
};

const partyColorToCss = (color?: PartyColor): string => {
  if (!color) {
    return '#888888';
  }

  const r = color.r ?? 0;
  const g = color.g ?? 0;
  const b = color.b ?? 0;
  const a = color.a ?? 1;

  const toHex = (value: number) => {
    const hex = Math.round(value).toString(16).padStart(2, '0');
    return hex;
  };

  const rgbHex = `#${toHex(r)}${toHex(g)}${toHex(b)}`;
  const alphaHex = a < 1 ? toHex(a * 255) : '';

  return `${rgbHex}${alphaHex}`;
};

const buildPartyPieData = (
  groupedByPosition: Record<Position, VoteResult[]> | undefined,
  sectionId?: string
): PieChartData => {
  const grouped = groupVotesByPartyForPosition(groupedByPosition, sectionId);
  const entries = Object.entries(grouped);
  const total = entries.reduce((s, [, votes]) => s + votes.length, 0);

  const slices = entries.map(([partyId, votes]) => {
    const politician = votes[0]?.politician;
    const party = politician?.party;
    const count = votes.length;
    const percentage =
      total > 0 ? Number(((count / total) * 100).toFixed(2)) : 0;

    return {
      id: partyId,
      label: party?.name ?? party?.acronym ?? partyId,
      value: { count, percentage },
      color: party ? partyColorToCss(party.color) : '#888888',
    };
  });

  return { slices };
};
