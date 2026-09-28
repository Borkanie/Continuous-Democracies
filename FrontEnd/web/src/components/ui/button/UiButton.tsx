import type { IconDefinition } from '@fortawesome/free-solid-svg-icons';
import Button from '@mui/material/Button';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

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
    <Button
      className={className}
      onClick={onClick}
      title={title}
      startIcon={icon && text ? <FontAwesomeIcon icon={icon} /> : undefined}
      sx={!text ? { minWidth: 0, padding: '4px' } : undefined}
    >
      {!text && icon && <FontAwesomeIcon icon={icon} />}
      {text}
    </Button>
  );
};
