import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faCheck, faXmark, faMinus } from '@fortawesome/free-solid-svg-icons';
import type { HydratedVote } from '../../utils/types';
import { PartyIcon } from '../ui/party-icon/PartyIcon';
import { positionColor } from '../../utils/helper';
import { useTheme } from '../../utils/context/ThemeContext';
import styles from './VoteListItem.module.css';

const { card, imageContainer, name, voteIcon } = styles;

const VOTE_ICON = {
  Yes: faCheck,
  No: faXmark,
  Abstain: faMinus,
  Absent: faMinus,
} as const;

type Props = { vote: HydratedVote; showPartyIcon?: boolean };

export const VoteListItem = (props: Props) => {
  const { vote, showPartyIcon = true } = props;
  const { theme } = useTheme();
  const color = positionColor(vote.value, theme);

  return (
    <li className={card}>
      <div className={imageContainer}>
        <img
          src={vote.politician.imageUrl || undefined}
          alt={vote.politician.name}
        />
      </div>
      <div className={name}>{vote.politician.name}</div>
      <FontAwesomeIcon
        icon={VOTE_ICON[vote.value]}
        className={voteIcon}
        style={{ color }}
        title={vote.value}
      />
      {showPartyIcon && <PartyIcon party={vote.party} size={28} />}
    </li>
  );
};
