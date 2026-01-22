import { faArrowLeft } from '@fortawesome/free-solid-svg-icons';
import styles from './header.module.css';
import { Button } from '../button/Button';
import { DateComponent } from '../date/DateComponent';

const { header, title, extraDetails, backButton } = styles;

type Props = {
  title: string;
  subtitle?: string;
  description?: string;
  onBack?: () => void;
  extraDetails?: {
    voteDate?: Date;
  };
};

export const Header = (props: Props) => {
  const {
    title: headerTitle,
    description,
    onBack,
    extraDetails: details,
  } = props;

  return (
    <div className={header}>
      <div className={title}>
        {onBack && (
          <Button className={backButton} icon={faArrowLeft} onClick={onBack} />
        )}
        <h3>{headerTitle}</h3>
      </div>
      {description && <p>{description}</p>}
      <div className={extraDetails}>
        {details?.voteDate && <DateComponent text={details.voteDate} />}
      </div>
    </div>
  );
};
