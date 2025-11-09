import { useParams, useRouter } from '@tanstack/react-router';
import sharedStyles from './styles/RoundBreakdown.module.css';
import styles from './styles/RoundSection.module.css';
import { Status } from '../components/status/Status';
import { useQuery } from '@tanstack/react-query';
import { getRound } from '../utils/api/rounds';
import { Spinner } from '../components/spinner/Spinner';
import { faArrowLeft } from '@fortawesome/free-solid-svg-icons';
import { Button } from '../components/button/Button';
import { DateComp } from '../components/Date/DateComp';

const { Div, header, title, extraDetails, separator } = sharedStyles;
const { flex } = styles;

export const RoundSection = () => {
  const { roundId, sectionId } = useParams({ strict: false });
  const router = useRouter();

  const { data: roundData, isFetching } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getRound(queryKey[1] || ''),
    staleTime: 5 * 60 * 1000, // cache data for 5 minutes
  });

  return (
    <div className={Div}>
      {isFetching ? (
        <Spinner />
      ) : (
        <>
          <div className={header}>
            <div className={title}>
              <div className={flex}>
                <Button
                  icon={faArrowLeft}
                  onClick={() => router.history.back()}
                />
                <h3>{roundData?.title}</h3>
              </div>
              <Status text={'ACTIV'} />
            </div>
            {/* TODO: Use roundData.description when available */}
            <p>
              Comprehensive legislation to reduce carbon emissions by 50% by
              2030
            </p>
            <div className={extraDetails}>
              {roundData?.voteDate && <DateComp text={roundData.voteDate} />}
            </div>
          </div>
          <div className={separator}></div>
        </>
      )}
    </div>
  );
};
