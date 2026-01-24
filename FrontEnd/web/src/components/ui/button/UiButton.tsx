import type { IconDefinition } from '@fortawesome/free-solid-svg-icons';
import styles from './UiButton.module.css';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classNames from 'classnames';
import { UiText } from '../text/UiText';

const { button, iconCls, textCls, hasBoth } = styles;

type Props = {
  onClick?: () => void;
  text?: string;
  title?: string;
  icon?: IconDefinition;
  className?: string;
};

export const UiButton = (props: Props) => {
  const { onClick, text, icon, className, title } = props;

  return (
    <button
      className={classNames(button, icon && text && hasBoth, className)}
      onClick={onClick}
      title={title}
    >
      {icon && <FontAwesomeIcon className={iconCls} icon={icon} />}
      {text && <UiText className={textCls} text={text} />}
    </button>
  );
};
