import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import styles from './DateComponent.module.css';
import { faCalendar } from '@fortawesome/free-regular-svg-icons';

const { Div, icon, date } = styles;

type Props = { text: Date };

export const DateComponent = (props: Props) => {
  const { text } = props;

  const dateText = new Date(text);

  return (
    <div className={Div}>
      <FontAwesomeIcon icon={faCalendar} className={icon} />
      <p className={date}>
        {dateText.toLocaleDateString('ro-RO', {
          year: 'numeric',
          month: 'numeric',
          day: 'numeric',
        })}
      </p>
    </div>
  );
};
