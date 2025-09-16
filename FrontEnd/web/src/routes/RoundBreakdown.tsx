import { useParams } from '@tanstack/react-router';
import { PieChart } from '../components/chart/PieChart';
import styles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getResultsByRoundId, getRound } from '../utils/api/rounds';
import { Status } from '../components/status/Status';
import { DateComp } from '../components/date/DateComp';
import { Legend } from '../components/legend/Legend';
import classNames from 'classnames';

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
  const params = useParams({ strict: false });
  const { roundId } = params;

  const { data: roundData } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getRound(queryKey[1] || ''),
    staleTime: 5 * 60 * 1000, // cache data for 5 minutes
  });

  const { data: roundResults } = useQuery({
    queryKey: ['roundResults', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getResultsByRoundId(queryKey[1] || ''),
    staleTime: 5 * 60 * 1000, // cache data for 5 minutes
  });

  if (!roundData) {
    // TODO: Add empty state here
    return <></>;
  }

  console.log({ roundData, roundResults });
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
            <PieChart isFirstLevel={true} />
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
