import classNames from 'classnames';
import React from 'react';
import styles from './styles.module.css';

export interface ButtonProps {
  color?: 'primary';
  className?: string;
  disabled?: boolean;
  onClick?: () => void;
}

const Button: React.FC<ButtonProps> = ({
  children,
  color,
  className,
  disabled,
  onClick,
}) => {
  const classes = classNames(styles.button, styles[color], className);
  return (
    <button className={classes} onClick={onClick} disabled={disabled}>
      {children}
    </button>
  );
};

export default Button;
