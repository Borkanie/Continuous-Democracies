import {
  useNavigate,
  useParams,
  Outlet,
  useMatches,
} from '@tanstack/react-router';
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
import { positionColor } from '../utils/helper';
import { UiText } from '../components/ui/text/UiText';

const {
  Div,
  separator,
  content,
  chartContainer,
  bold,
  chartTitle,
  highlight,
  info,
} = sharedStyles;

export const RoundSection = () => {
  const { roundId, sectionId } = useParams({ strict: false });
  const navigate = useNavigate();
  const matches = useMatches();
  const hasChildRoute = matches.length > 3;

  const { data: roundData, isFetching } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getRound(queryKey[1] || ''),
  });
  const { data: groupedRoundResults } = useResultsByRoundId(roundId);

  const partyPieData = buildPartyPieData(groupedRoundResults, sectionId);

  return (
    <div className={Div}>
      {isFetching ? (
        <Spinner />
      ) : (
        <>
          {hasChildRoute ? (
            <Outlet />
          ) : (
            <>
              <Header
                title={roundData?.title || ''}
                extraDetails={{ voteDate: roundData?.voteDate }}
                description={roundData?.description}
              />

              <div className={separator}></div>

              {/* Party distribution pie + legend for selected section */}
              <div className={content}>
                <div className={chartContainer}>
                  <div>
                    <p className={classNames(bold, chartTitle)}>
                      Impartirea pe partide pentru:{' '}
                      <span
                        className={highlight}
                        style={{
                          color: positionColor(Number(sectionId) as Position),
                        }}
                      >
                        {SECTION_LABELS[Number(sectionId) as Position]}
                      </span>
                    </p>
                    <PieChart
                      data={partyPieData}
                      onSliceClick={(id) => navigate({ to: `party/${id}` })}
                    />
                    <UiText
                      className={info}
                      size={'small'}
                      text={
                        'Apasati click pe o sectiune pentru a vedea lista votantilor din partid'
                      }
                    />
                  </div>
                  <div>
                    <Legend slices={partyPieData.slices} />
                  </div>
                </div>
              </div>
            </>
          )}
        </>
      )}
    </div>
  );
};

export const groupVotesByPartyForPosition = (
  grouped: Record<Position, VoteResult[]> | undefined,
  sectionId?: string,
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

export const buildPartyPieData = (
  groupedByPosition: Record<Position, VoteResult[]> | undefined,
  sectionId?: string,
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
      acronym: party?.acronym || '',
    };
  });

  return { slices };
};
