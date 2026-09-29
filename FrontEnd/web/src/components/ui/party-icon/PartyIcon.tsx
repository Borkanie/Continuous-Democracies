import classNames from 'classnames';
import type { Party } from '../../../utils/types';
import styles from './PartyIcon.module.css';

const { badge } = styles;

type Props = {
  party?: Party;
  size?: number;
  className?: string;
};

export const PartyIcon = (props: Props) => {
  const { party, size = 24, className } = props;

  return (
    <div
      className={classNames(badge, className)}
      title={party?.name}
      style={{
        width: size,
        height: size,
        fontSize: size * 0.32,
        backgroundColor: party?.logoUrl ? undefined : party?.color || '#888888',
      }}
    >
      {party?.logoUrl ? (
        <img src={party.logoUrl} alt={party.name} />
      ) : (
        <span>{party?.acronym?.slice(0, 4)}</span>
      )}
    </div>
  );
};
