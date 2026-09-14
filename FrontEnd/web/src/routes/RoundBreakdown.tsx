import {
  useNavigate,
  useParams,
  Outlet,
  useMatches,
} from '@tanstack/react-router';
import { PieChart, type PieChartData } from '../components/chart/PieChart';
import styles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getVotingRoundById } from '../utils/api/rounds';
import { Legend } from '../components/legend/Legend';
import classNames from 'classnames';
import {
  useResultsByRoundId,
  type GroupedVotes,
} from '../utils/hooks/useResultsByRoundId';
import type { VoteValue } from '../utils/types';
import { Spinner } from '../components/spinner/Spinner';
import { Header } from '../components/section/header';
import { positionColor, positionLabel, getVoteStatus } from '../utils/helper';
import { useTheme, type Theme } from '../utils/context/ThemeContext';
import { UiText } from '../components/ui/text/UiText';

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
    queryFn: ({ queryKey }) => getVotingRoundById(queryKey[1] || ''),
  });

  const { data: groupedRoundResults } = useResultsByRoundId(roundId);
  const { theme } = useTheme();

  const roundBreakdownData = buildRoundBreakdownData(groupedRoundResults, theme);

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
                status={
                  groupedRoundResults && roundData
                    ? getVoteStatus(groupedRoundResults, roundData.majorityType)
                    : undefined
                }
                majorityType={roundData?.majorityType}
              />

              <div className={separator}></div>

              <div className={content}>
                <div className={chartContainer}>
                  <div>
                    <UiText
                      className={classNames(bold, chartTitle)}
                      text={'Distributia voturilor'}
                    />
                    <PieChart
                      data={roundBreakdownData}
                      onSliceClick={(id) => navigate({ to: `section/${id}` })}
                    />
                    <UiText
                      className={info}
                      text={
                        'Apasa pe o sectiune din grafic sau pe legenda pentru a vedea impartirea pe partide'
                      }
                      size={'medium'}
                    />
                  </div>
                  <div>
                    <Legend
                      text={'Legenda voturilor'}
                      slices={roundBreakdownData.slices}
                      onSliceClick={(id) => navigate({ to: `section/${id}` })}
                    />
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
  theme: Theme,
): PieChartData => {
  if (!groupedRoundResults) {
    return { slices: [] };
  }

  const voteValues: VoteValue[] = ['Yes', 'No', 'Abstain', 'Absent'];
  const total = voteValues.reduce(
    (sum, value) => sum + (groupedRoundResults[value]?.length ?? 0),
    0,
  );

  const slices = voteValues.map((value) => {
    const count = groupedRoundResults[value]?.length ?? 0;
    const percentage =
      total > 0 ? Number(((count / total) * 100).toFixed(2)) : 0;

    return {
      id: value,
      label: positionLabel(value),
      value: { count, percentage },
      color: positionColor(value, theme),
    };
  });

  return { slices };
};
