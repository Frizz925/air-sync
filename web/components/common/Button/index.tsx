import clsx from 'clsx';
import React from 'react';
import styles from './styles.module.css';

export type ButtonColor = 'primary' | 'red';

export interface ButtonProps {
  color?: ButtonColor;
  className?: string;
  disabled?: boolean;
  rounded?: boolean;
  onClick?: () => void;
}

const Button: React.FC<ButtonProps> = ({
  children,
  color,
  className,
  disabled,
  rounded,
  onClick,
}) => {
  const classes = clsx(
    styles.button,
    rounded && styles.rounded,
    styles[color],
    className
  );
  return (
    <button className={classes} onClick={onClick} disabled={disabled}>
      {children}
    </button>
  );
};

export default Button;
