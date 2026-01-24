import type { JSX } from 'react';

type Props = {
  text: string;
  level?: 1 | 2 | 3 | 4 | 5 | 6;
};

export const UiHeading = ({ text, level = 1 }: Props) => {
  const HeadingTag = `h${level}` as keyof JSX.IntrinsicElements;
  return <HeadingTag>{text}</HeadingTag>;
};
