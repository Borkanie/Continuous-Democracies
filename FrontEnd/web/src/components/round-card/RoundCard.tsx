import type { Round } from '../../utils/api/rounds';
import styles from './RoundCard.module.css';

const { Div, title, description, footer, status, date, icon, header } = styles;

type Props = {
  round: Round;
};

export const RoundCard = (props: Props) => {
  const { round } = props;

  // TODO: Use real data from `round` prop when available for date, description, and status
  return (
    <div className={Div}>
      <div className={header}>
        <h3 className={title}>{round.title}</h3>
      </div>
      <p className={description}>
        Comprehensive legislation to reduce carbon emissions by 50% by 2030
      </p>
      <div className={footer}>
        <p className={date}>25.11.2024</p>
        <div className={status}>
          <div className={icon}></div>
          <p>Activa</p>
        </div>
      </div>
    </div>
  );
};
