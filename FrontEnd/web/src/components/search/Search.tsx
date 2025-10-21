import { faSearch } from '@fortawesome/free-solid-svg-icons';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import styles from './Search.module.css';

const { Div, iconContainer, icon } = styles;

type Props = {
  value: string;
  onChange: (value: string) => void;
};

export const Search = (props: Props) => {
  const { value, onChange } = props;

  return (
    <div className={Div}>
      <div className={iconContainer}>
        <FontAwesomeIcon icon={faSearch} className={icon} />
      </div>
      <input
        placeholder={'Cauta...'}
        value={value}
        onChange={(e) => onChange(e.target.value)}
      />
    </div>
  );
};
