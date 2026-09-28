import { useNavigate, useParams, Outlet, useMatches } from '@tanstack/react-router';
import sharedStyles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getVotingRoundById } from '../utils/api/rounds';
import { Spinner } from '../components/spinner/Spinner';
import {
  useResultsByRoundId,
  type GroupedVotes,
} from '../utils/hooks/useResultsByRoundId';
import type { PartyColor, VoteValue, HydratedVote } from '../utils/types';
import { PieChart, type PieChartData, type Slice } from '../components/chart/PieChart';
import { Legend } from '../components/legend/Legend';
import classNames from 'classnames';
import { Header } from '../components/section/header';
import { SECTION_LABELS } from '../utils/constants';
import {
  positionColor,
  positionLabel,
  haze,
  getVoteStatus,
} from '../utils/helper';
import { useTheme, type Theme } from '../utils/context/ThemeContext';
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

const VOTE_VALUES: VoteValue[] = ['Yes', 'No', 'Abstain', 'Absent'];

/*
 * Radial offsets (px) for the exploded wheel: the selected value's party
 * slices pop out to stand out; the other three vote-value slices are pulled
 * toward center so they read as visually smaller/recessed - Chart.js applies
 * `offset` as a plain radial translation, so a negative value pulls inward.
 */
const SELECTED_SLICE_OFFSET = 18;
const UNSELECTED_SLICE_OFFSET = -12;

export const RoundSection = () => {
  const { roundId, sectionId } = useParams({ strict: false });
  const navigate = useNavigate();
  const matches = useMatches();
  const hasChildRoute = matches.length > 3;

  const { data: roundData, isFetching } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getVotingRoundById(queryKey[1] || ''),
  });
  const { data: groupedRoundResults } = useResultsByRoundId(roundId);
  const { theme } = useTheme();

  const explodedData = buildExplodedSectionData(
    groupedRoundResults,
    sectionId as VoteValue | undefined,
    theme,
  );

  const handleSliceClick = (id: string | number) => {
    const value = id.toString();
    if (VOTE_VALUES.includes(value as VoteValue)) {
      if (value !== sectionId) {
        navigate({
          to: '/round/$roundId/section/$sectionId',
          params: { roundId: roundId || '', sectionId: value },
        });
      }
      return;
    }

    navigate({ to: `party/${value}` });
  };

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
                    ? getVoteStatus(groupedRoundResults, roundData.majorityType, roundData.chamber)
                    : undefined
                }
                majorityType={roundData?.majorityType}
              />

              <div className={separator}></div>

              {/* Full round wheel, with the selected section split into parties and the rest dimmed */}
              <div className={content}>
                <div className={chartContainer}>
                  <div>
                    <p className={classNames(bold, chartTitle)}>
                      Impartirea pe partide pentru:{' '}
                      <span
                        className={highlight}
                        style={{
                          color: positionColor(sectionId as VoteValue, theme),
                        }}
                      >
                        {SECTION_LABELS[sectionId as VoteValue]}
                      </span>
                    </p>
                    <PieChart data={explodedData} onSliceClick={handleSliceClick} />
                    <UiText
                      className={info}
                      size={'medium'}
                      text={
                        'Apasa pe o sectiune din grafic sau pe legenda pentru a vedea lista votantilor per partid'
                      }
                    />
                  </div>
                  <div>
                    <Legend
                      text={'Legenda partidelor'}
                      slices={explodedData.slices}
                      onSliceClick={handleSliceClick}
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

export const groupVotesByPartyForPosition = (
  grouped: GroupedVotes | undefined,
  sectionId?: string,
): Record<string, HydratedVote[]> => {
  if (!grouped || !sectionId || !VOTE_VALUES.includes(sectionId as VoteValue)) {
    return {};
  }

  const votes = grouped[sectionId as VoteValue] ?? [];
  return votes.reduce<Record<string, HydratedVote[]>>((acc, vote) => {
    const partyId = vote.party?.id ?? 'unknown';

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
  return color;
};

/**
 * Builds one wheel covering the whole round: the selected vote value is split
 * into one slice per party (full opacity), the other three vote values stay
 * as single, dimmed slices in their original position. Percentages for every
 * slice are relative to the whole round's vote count, so they sum to 100%
 * together (matching the angles chart.js actually draws, which are
 * proportional to raw counts).
 */
const buildExplodedSectionData = (
  groupedRoundResults: GroupedVotes | undefined,
  selectedValue: VoteValue | undefined,
  theme: Theme,
): PieChartData => {
  if (!groupedRoundResults) {
    return { slices: [] };
  }

  const total = VOTE_VALUES.reduce(
    (sum, value) => sum + (groupedRoundResults[value]?.length ?? 0),
    0,
  );
  const percentageOf = (count: number) =>
    total > 0 ? Number(((count / total) * 100).toFixed(2)) : 0;

  const slices: Slice[] = VOTE_VALUES.flatMap((value) => {
    if (value === selectedValue) {
      const byParty = groupVotesByPartyForPosition(
        groupedRoundResults,
        selectedValue,
      );
      return Object.entries(byParty).map(([partyId, votes]) => {
        const party = votes[0]?.party;
        const count = votes.length;

        return {
          id: partyId,
          label: party?.name ?? party?.acronym ?? partyId,
          value: { count, percentage: percentageOf(count) },
          color: party ? partyColorToCss(party.color) : '#888888',
          acronym: party?.acronym || '',
          offset: SELECTED_SLICE_OFFSET,
        };
      });
    }

    const count = groupedRoundResults[value]?.length ?? 0;
    return [
      {
        id: value,
        label: positionLabel(value),
        value: { count, percentage: percentageOf(count) },
        color: haze(positionColor(value, theme), theme),
        acronym: '',
        offset: UNSELECTED_SLICE_OFFSET,
      },
    ];
  });

  return { slices };
};
