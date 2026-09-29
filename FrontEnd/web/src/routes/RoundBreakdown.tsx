import {
  useNavigate,
  useParams,
  Outlet,
  useMatches,
} from '@tanstack/react-router';
import { useMemo, useState } from 'react';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import FormControl from '@mui/material/FormControl';
import InputLabel from '@mui/material/InputLabel';
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContent from '@mui/material/DialogContent';
import { faXmark, faList } from '@fortawesome/free-solid-svg-icons';
import { PieChart, type PieChartData } from '../components/chart/PieChart';
import styles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getVotingRoundById } from '../utils/api/rounds';
import { Legend } from '../components/legend/Legend';
import classNames from 'classnames';
import {
  useResultsByRoundId,
  groupVotesByValue,
  type GroupedVotes,
} from '../utils/hooks/useResultsByRoundId';
import type { VoteValue, HydratedVote, Party } from '../utils/types';
import { Spinner } from '../components/spinner/Spinner';
import { Header } from '../components/section/header';
import { PoliticiansList } from '../components/politicians-list/PoliticiansList';
import { UiButton } from '../components/ui/button/UiButton';
import { positionColor, positionLabel, getVoteStatus } from '../utils/helper';
import { useTheme, type Theme } from '../utils/context/ThemeContext';
import { UiText } from '../components/ui/text/UiText';

const {
  Div,
  separator,
  content,
  bold,
  chartContainer,
  chartTitle,
  info,
  filterRow,
  dialogContent,
} = styles;

const VOTE_VALUES: VoteValue[] = ['Yes', 'No', 'Abstain', 'Absent'];

export const RoundBreakdown = () => {
  const navigate = useNavigate();
  const params = useParams({ strict: false });
  const { roundId } = params;
  const matches = useMatches();
  const hasChildRoute = matches.length > 2;
  const [selectedPartyId, setSelectedPartyId] = useState('');
  const [allVotesOpen, setAllVotesOpen] = useState(false);
  const { theme } = useTheme();

  const { data: roundData, isFetching } = useQuery({
    queryKey: ['roundById', roundId],
    enabled: !!roundId,
    queryFn: ({ queryKey }) => getVotingRoundById(queryKey[1] || ''),
  });

  const { data: groupedRoundResults } = useResultsByRoundId(roundId);

  const allVotes = useMemo(
    () => flattenVotes(groupedRoundResults),
    [groupedRoundResults],
  );
  const parties = useMemo(() => uniqueParties(allVotes), [allVotes]);
  const selectedParty = parties.find((p) => p.id === selectedPartyId);
  const groupedForChart = useMemo(
    () =>
      selectedPartyId
        ? groupVotesByValue(
            allVotes.filter((v) => v.party?.id === selectedPartyId),
          )
        : groupedRoundResults,
    [allVotes, selectedPartyId, groupedRoundResults],
  );
  const roundBreakdownData = buildRoundBreakdownData(groupedForChart, theme);

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

              <div className={content}>
                <div className={chartContainer}>
                  <div>
                    <UiText
                      className={classNames(bold, chartTitle)}
                      text={
                        selectedParty
                          ? `Distributia voturilor - ${selectedParty.name}`
                          : 'Distributia voturilor'
                      }
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

                <div className={filterRow}>
                  <FormControl size={'small'} sx={{ minWidth: 220 }}>
                    <InputLabel id={'party-filter-label'}>
                      Filtreaza pe partid
                    </InputLabel>
                    <Select
                      labelId={'party-filter-label'}
                      label={'Filtreaza pe partid'}
                      value={selectedPartyId}
                      onChange={(e) => setSelectedPartyId(e.target.value)}
                    >
                      <MenuItem value={''}>
                        <em>Niciunul</em>
                      </MenuItem>
                      {parties.map((party) => (
                        <MenuItem key={party.id} value={party.id}>
                          {party.name} ({party.acronym})
                        </MenuItem>
                      ))}
                    </Select>
                  </FormControl>

                  <UiButton
                    icon={faList}
                    text={'Vezi toate voturile'}
                    onClick={() => setAllVotesOpen(true)}
                  />
                </div>
              </div>

              <Dialog
                open={allVotesOpen}
                onClose={() => setAllVotesOpen(false)}
                fullWidth
              >
                <DialogTitle
                  sx={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                  }}
                >
                  Toate voturile
                  <UiButton
                    icon={faXmark}
                    title={'Inchide'}
                    onClick={() => setAllVotesOpen(false)}
                  />
                </DialogTitle>
                <DialogContent className={dialogContent}>
                  <PoliticiansList vote={allVotes} />
                </DialogContent>
              </Dialog>
            </>
          )}
        </>
      )}
    </div>
  );
};

const flattenVotes = (
  groupedRoundResults: GroupedVotes | undefined,
): HydratedVote[] => {
  if (!groupedRoundResults) {
    return [];
  }
  return VOTE_VALUES.flatMap((value) => groupedRoundResults[value] ?? []);
};

const uniqueParties = (votes: HydratedVote[]): Party[] => {
  const byId = new Map<string, Party>();
  votes.forEach((v) => {
    if (v.party && !byId.has(v.party.id)) {
      byId.set(v.party.id, v.party);
    }
  });
  return Array.from(byId.values()).sort((a, b) => a.name.localeCompare(b.name));
};

const buildRoundBreakdownData = (
  groupedRoundResults: GroupedVotes | undefined,
  theme: Theme,
): PieChartData => {
  if (!groupedRoundResults) {
    return { slices: [] };
  }

  const total = VOTE_VALUES.reduce(
    (sum, value) => sum + (groupedRoundResults[value]?.length ?? 0),
    0,
  );

  const slices = VOTE_VALUES.map((value) => {
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
