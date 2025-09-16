import { useNavigate, useParams } from '@tanstack/react-router';
import { PieChart, type DrilldownEntry } from '../components/chart/PieChart';
import styles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getRound } from '../utils/api/rounds';
import { Status } from '../components/status/Status';
import { DateComp } from '../components/date/DateComp';
import { Legend } from '../components/legend/Legend';
import classNames from 'classnames';
import { useResultsByRoundId } from '../utils/hooks/useResultsByRoundId';
import type { Position } from '../utils/types';

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

  const topLevelData: DrilldownEntry = {
    id: '0',
    label: 'Vot',
    children: Object.entries(groupedRoundResults).map(([key, votes]) => ({
      id: key,
      label: positionLabel(Number(key) as Position),
      value: votes.length,
    })),
  };

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
              data={topLevelData}
              onSliceClick={(id) => navigate({ to: `section/${id}` })}
            />
            <p className={info}>
              Apasati click pe o sectiune pentru a vedea impartirea pe partide
            </p>
          </div>
          <div>
            <Legend />
          </div>
        </div>
      </div>
    </div>
  );
};
