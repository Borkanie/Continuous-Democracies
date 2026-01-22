import type { IconDefinition } from '@fortawesome/free-solid-svg-icons';
import styles from './Button.module.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classNames from 'classnames';

const { button, iconCls } = styles;

type Props = {
  onClick?: () => void;
  text?: string;
  icon?: IconDefinition;
  className?: string;
};

export const Button = (props: Props) => {
  const { onClick, text, icon, className } = props;

  return (
    <button className={classNames(button, className)} onClick={onClick}>
      {icon && <FontAwesomeIcon className={iconCls} icon={icon} />}
      {text && <span>{text}</span>}
    </button>
  );
};
