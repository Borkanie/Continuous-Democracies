import { faArrowLeft } from '@fortawesome/free-solid-svg-icons';
import styles from './header.module.css';
import { Button } from '../button/Button';
import { Status } from '../status/Status';
import { DateComponent } from '../date/DateComponent';

const { header, title, extraDetails, flex } = styles;

type Props = {
  title: string;
  subtitle?: string;
  description?: string;
  onBack?: () => void;
  status?: string;
  extraDetails?: {
    voteDate?: Date;
  };
};

export const Header = (props: Props) => {
  const {
    title: headerTitle,
    description,
    onBack,
    status,
    extraDetails: details,
  } = props;

  return (
    <div className={header}>
      <div className={title}>
        <div className={flex}>
          {onBack && <Button icon={faArrowLeft} onClick={onBack} />}
          <h3>{headerTitle}</h3>
        </div>
        {status && <Status text={'ACTIV'} />}
      </div>
      {/* TODO: Use roundData.description when available */}
      {/* {description && <p>{description}</p>} */}
      <p>Comprehensive legislation to reduce carbon emissions by 50% by 2030</p>
      <div className={extraDetails}>
        {details?.voteDate && <DateComponent text={details.voteDate} />}
      </div>
    </div>
  );
};
