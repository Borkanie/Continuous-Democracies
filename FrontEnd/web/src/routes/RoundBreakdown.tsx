import { useParams } from '@tanstack/react-router';
import { PieChart } from '../components/chart/PieChart';
import styles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getResultsByRoundId, getRound } from '../utils/api/rounds';
import { Status } from '../components/status/Status';
import { DateComp } from '../components/Date/DateComp';

const { Div, details, extraDetails, date } = styles;

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
      <div className={details}>
        <h2>{roundData.title}</h2>
        {/* TODO: Use roundData.description when available */}
        <p>
          Comprehensive legislation to reduce carbon emissions by 50% by 2030
        </p>
        <div className={extraDetails}>
          <Status text={'ACTIV'} />
          <DateComp text={roundData.voteDate} />
        </div>
      </div>
      <PieChart isFirstLevel={true} />
    </div>
  );
};
