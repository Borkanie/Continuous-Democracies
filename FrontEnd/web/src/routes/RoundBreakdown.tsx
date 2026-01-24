import {
  useNavigate,
  useParams,
  Outlet,
  useMatches,
} from '@tanstack/react-router';
import { PieChart, type PieChartData } from '../components/chart/PieChart';
import styles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getRound } from '../utils/api/rounds';
import { Legend } from '../components/legend/Legend';
import classNames from 'classnames';
import {
  useResultsByRoundId,
  type GroupedVotes,
} from '../utils/hooks/useResultsByRoundId';
import type { Position } from '../utils/types';
import { VOTERS_TOTAL_NUMBER } from '../utils/constants';
import { Spinner } from '../components/spinner/Spinner';
import { Header } from '../components/section/header';
import { positionColor, positionLabel } from '../utils/helper';

const { Div, separator, content, bold, chartContainer, chartTitle, info } =
  styles;

export const RoundBreakdown = () => {
  const navigate = useNavigate();
  const params = useParams({ strict: false });
  const { roundId } = params;
  const matches = useMatches();
  const hasChildRoute = matches.length > 2;

  const { data: roundData, isFetching } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getRound(queryKey[1] || ''),
    staleTime: 5 * 60 * 1000, // cache data for 5 minutes
  });

  const { data: groupedRoundResults } = useResultsByRoundId(roundId);

  const roundBreakdownData = buildRoundBreakdownData(groupedRoundResults);

  return (
    <div className={Div}>
      {isFetching ? (
        <Spinner />
      ) : (
        <>
          <Outlet />
          {!hasChildRoute && (
            <>
              <Header
                title={roundData?.title || ''}
                extraDetails={{ voteDate: roundData?.voteDate }}
                description={roundData?.description}
              />

              <div className={separator}></div>

              <div className={content}>
                <div className={chartContainer}>
                  <div>
                    <p className={classNames(bold, chartTitle)}>
                      Distributia voturilor
                    </p>
                    <PieChart
                      data={roundBreakdownData}
                      onSliceClick={(id) => navigate({ to: `section/${id}` })}
                    />
                    <p className={info}>
                      Apasati click pe o sectiune pentru a vedea impartirea pe
                      partide
                    </p>
                  </div>
                  <div>
                    <Legend slices={roundBreakdownData.slices} />
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

const buildRoundBreakdownData = (
  groupedRoundResults: GroupedVotes | undefined,
): PieChartData => {
  if (!groupedRoundResults) {
    return { slices: [] };
  }

  const total = VOTERS_TOTAL_NUMBER;
  const nonAbsentPositions: Position[] = [0, 1, 2];

  // Build slices for Da, Nu, Abtinere
  const slices = nonAbsentPositions.map((pos) => {
    const count = groupedRoundResults[pos]?.length ?? 0;
    const percentage =
      total > 0 ? Number(((count / total) * 100).toFixed(2)) : 0;

    return {
      id: pos,
      label: positionLabel(pos),
      value: { count, percentage },
      color: positionColor(pos),
    };
  });

  // Compute "Absent" as the remainder
  const usedVotes = slices.reduce((sum, slice) => sum + slice.value.count, 0);
  const absentCount = Math.max(0, total - usedVotes);
  const absentPercentage =
    total > 0 ? Number(((absentCount / total) * 100).toFixed(2)) : 0;

  slices.push({
    id: 3,
    label: positionLabel(3),
    value: { count: absentCount, percentage: absentPercentage },
    color: positionColor(3),
  });

  return { slices };
};
