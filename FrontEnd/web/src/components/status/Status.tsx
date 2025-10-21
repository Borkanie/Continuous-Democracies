import styles from './Status.module.css';

const { Div, statusText } = styles;

type Props = { text: string };

export const Status = (props: Props) => {
  const { text } = props;

  return (
    <div className={Div}>
      <p className={statusText}>{text}</p>
    </div>
  );
};
