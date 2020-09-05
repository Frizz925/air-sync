import { IconProp } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import clsx from 'clsx';
import React from 'react';
import styles from './styles.module.css';

export interface IconButtonProps {
  size?: string | number;
  color?: string;
  icon: IconProp;
  title?: string;
  onClick?: () => void;
}

const IconButton: React.FC<IconButtonProps> = ({
  size,
  color,
  icon,
  title,
  onClick,
}) => {
  const classes = clsx(styles.container, styles[color]);
  const style = size
    ? ({ height: size, width: size } as React.CSSProperties)
    : null;
  return (
    <button className={classes} style={style} title={title} onClick={onClick}>
      <div className={styles.icon}>
        <FontAwesomeIcon icon={icon} />
      </div>
    </button>
  );
};

export default IconButton;
