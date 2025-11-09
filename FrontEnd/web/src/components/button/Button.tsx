import type { IconDefinition } from '@fortawesome/free-solid-svg-icons';
import styles from './Button.module.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

const { button, iconCls } = styles;

type Props = {
  onClick?: () => void;
  text?: string;
  icon?: IconDefinition;
};

export const Button = (props: Props) => {
  const { onClick, text, icon } = props;

  return (
    <button className={button} onClick={onClick}>
      {icon && <FontAwesomeIcon className={iconCls} icon={icon} />}
      {text && <span>{text}</span>}
    </button>
  );
};
