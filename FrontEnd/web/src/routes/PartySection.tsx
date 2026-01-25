import { useParams } from '@tanstack/react-router';
import { Header } from '../components/section/header';
import sharedStyles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getRound } from '../utils/api/rounds';
import classNames from 'classnames';
import { UiText } from '../components/ui/text/UiText';

const { Div, separator, content, chartContainer, bold, chartTitle } =
  sharedStyles;

export const PartySection = () => {
  const { roundId, sectionId } = useParams({ strict: false });

  const { data: roundData } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getRound(queryKey[1] || ''),
  });

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
        <div className={chartContainer}>
          <div>
            <UiText
              className={classNames(bold, chartTitle)}
              text={`Lista politicieni - sectiunea ${sectionId}`}
            />
          </div>
        </div>
      </div>
    </div>
  );
};
