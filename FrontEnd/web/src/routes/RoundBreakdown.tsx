import {
  useNavigate,
  useParams,
  Outlet,
  useMatches,
} from '@tanstack/react-router';
import { useMemo, useState } from 'react';
import Dialog from '@mui/material/Dialog';
import DialogTitle from '@mui/material/DialogTitle';
import DialogContent from '@mui/material/DialogContent';
import { faXmark, faList } from '@fortawesome/free-solid-svg-icons';
import styles from './styles/RoundBreakdown.module.css';
import { useQuery } from '@tanstack/react-query';
import { getVotingRoundById } from '../utils/api/rounds';
import {
  useResultsByRoundId,
  type GroupedVotes,
} from '../utils/hooks/useResultsByRoundId';
import type { VoteValue, HydratedVote } from '../utils/types';
import { Spinner } from '../components/spinner/Spinner';
import { Header } from '../components/section/header';
import { PoliticiansList } from '../components/politicians-list/PoliticiansList';
import { UiButton } from '../components/ui/button/UiButton';
import { VoteSummaryBar } from '../components/vote-summary-bar/VoteSummaryBar';
import { getVoteStatus } from '../utils/helper';
import { PartyVoteRows } from '../components/party-vote-rows/PartyVoteRows';

const { Div, separator, content, filterRow, dialogContent } = styles;

const VOTE_VALUES: VoteValue[] = ['Yes', 'No', 'Abstain', 'Absent'];

export const RoundBreakdown = () => {
  const navigate = useNavigate();
  const params = useParams({ strict: false });
  const { roundId } = params;
  const matches = useMatches();
  const hasChildRoute = matches.length > 2;
  const [allVotesOpen, setAllVotesOpen] = useState(false);

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
                {groupedRoundResults && roundData && (
                  <VoteSummaryBar
                    grouped={groupedRoundResults}
                    majorityType={roundData.majorityType}
                    chamber={roundData.chamber}
                  />
                )}

                <PartyVoteRows
                  votes={allVotes}
                  onPartyClick={(partyId) => navigate({ to: `party/${partyId}` })}
                />

                <div className={filterRow}>
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
