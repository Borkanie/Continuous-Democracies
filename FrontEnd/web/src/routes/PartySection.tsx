import { useParams, useRouter } from '@tanstack/react-router';
import { Header } from '../components/section/header';
import sharedStyles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getRound } from '../utils/api/rounds';
import classNames from 'classnames';

const { Div, separator, content, chartContainer, bold, chartTitle } =
  sharedStyles;

export const PartySection = () => {
  const { roundId, sectionId } = useParams({ strict: false });

  const router = useRouter();

  const { data: roundData, isFetching } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getRound(queryKey[1] || ''),
    staleTime: 5 * 60 * 1000, // cache data for 5 minutes
  });

  const title = `${roundData?.title} `;

  return (
    <div className={Div}>
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
              Lista politicieni - sectiunea {sectionId}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};
