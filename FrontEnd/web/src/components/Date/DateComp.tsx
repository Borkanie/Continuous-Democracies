import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import styles from './DateComp.module.css';
import { faCalendar } from '@fortawesome/free-regular-svg-icons';

const { Div, icon, date } = styles;

type Props = { text: Date };

export const DateComp = (props: Props) => {
  const { text } = props;

  const datee = new Date(text);

  return (
    <div className={Div}>
      <FontAwesomeIcon icon={faCalendar} className={icon} />
      <p className={date}>
        {datee.toLocaleDateString('ro-RO', {
          year: 'numeric',
          month: 'numeric',
          day: 'numeric',
        })}
      </p>
    </div>
  );
};
