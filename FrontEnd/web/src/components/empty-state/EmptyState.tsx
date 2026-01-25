import { UiText } from '../ui/text/UiText';
import styles from './EmptyState.module.css';

const { Div, pText } = styles;

type Props = { text?: string };

export const EmptyState = (props: Props) => {
  const { text } = props;

  return (
    <div className={Div}>
      {text && <UiText className={pText} text={text} />}
    </div>
  );
};
