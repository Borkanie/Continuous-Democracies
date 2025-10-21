import styles from './EmptyState.module.css';

const { Div, pText } = styles;

type Props = { text?: string };

export const EmptyState = (props: Props) => {
  const { text } = props;

  return <div className={Div}>{text && <p className={pText}>{text}</p>}</div>;
};
