import classNames from 'classnames';
import styles from './Text.module.css';

const { textComp, trunc } = styles;

type Size = 'small' | 'medium' | 'large' | 'xlarge';

type Props = {
  text: string;
  title: string;
  truncate?: boolean;
  maxLines?: number;
  size?: Size;
};

export const Text = (props: Props) => {
  const { text, title, truncate = false, maxLines = 3, size = 'large' } = props;

  const sizeStyle = {
    small: { fontSize: 'var(--theme-font-size-s)' },
    medium: { fontSize: 'var(--theme-font-size-m)' },
    large: { fontSize: 'var(--theme-font-size-l)' },
    xlarge: { fontSize: 'var(--theme-font-size-xl)' },
  };

  const lineClamp = truncate
    ? ({ '--text-line-clamp': String(maxLines) } as React.CSSProperties)
    : {};

  return (
    <p
      className={classNames('Text', textComp, truncate && trunc)}
      style={{ ...lineClamp, ...sizeStyle[size] }}
      title={title}
    >
      {text}
    </p>
  );
};
