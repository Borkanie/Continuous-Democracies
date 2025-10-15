import { useNavigate, useParams } from '@tanstack/react-router';
import { PieChart, type PieChartData } from '../components/chart/PieChart';
import styles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getRound } from '../utils/api/rounds';
import { Status } from '../components/status/Status';
import { DateComp } from '../components/date/DateComp';
import { Legend } from '../components/legend/Legend';
import classNames from 'classnames';
import {
  useResultsByRoundId,
  type GroupedVotes,
} from '../utils/hooks/useResultsByRoundId';
import type { Position } from '../utils/types';
import { VOTERS_TOTAL_NUMBER } from '../utils/constants';

const {
  Div,
  header,
  title,
  extraDetails,
  separator,
  content,
  bold,
  chartContainer,
  chartTitle,
  info,
} = styles;

export const RoundBreakdown = () => {
  const navigate = useNavigate();
  const params = useParams({ strict: false });
  const { roundId } = params;

  const { data: roundData } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getRound(queryKey[1] || ''),
    staleTime: 5 * 60 * 1000, // cache data for 5 minutes
  });

  const { data: groupedRoundResults } = useResultsByRoundId(roundId);

  if (!roundData || !groupedRoundResults) {
    // TODO: Add empty state here
    return <></>;
  }

  const roundBreakdownData = buildRoundBreakdownData(groupedRoundResults);

  return (
    <div className={Div}>
      <div className={header}>
        <div className={title}>
          <h3>{roundData.title}</h3>
          <Status text={'ACTIV'} />
        </div>
        {/* TODO: Use roundData.description when available */}
        <p>
          Comprehensive legislation to reduce carbon emissions by 50% by 2030
        </p>
        <div className={extraDetails}>
          <DateComp text={roundData.voteDate} />
        </div>
      </div>
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
              Apasati click pe o sectiune pentru a vedea impartirea pe partide
            </p>
          </div>
          <div>
            <Legend slices={roundBreakdownData.slices} />
          </div>
        </div>
      </div>
    </div>
  );
};

const buildRoundBreakdownData = (
  groupedRoundResults: GroupedVotes
): PieChartData => {
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
  });

  return { slices };
};

const positionLabel = (position: Position): string => {
  switch (position) {
    case 0:
      return 'Da';
    case 1:
      return 'Nu';
    case 2:
      return 'Abtinere';
    case 3:
      return 'Absent';
    default:
      return 'Unknown';
  }
};
