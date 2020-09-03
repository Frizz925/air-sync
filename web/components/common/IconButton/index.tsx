import { IconProp } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classNames from 'classnames';
import React from 'react';
import styles from './styles.module.css';

export interface IconButtonProps {
  size?: string | number;
  color?: string;
  icon: IconProp;
  onClick: () => void;
}

const IconButton: React.FC<IconButtonProps> = ({
  size,
  color,
  icon,
  onClick,
}) => {
  const classes = classNames(styles.container, styles[color]);
  const style = size
    ? ({ height: size, width: size } as React.CSSProperties)
    : null;
  return (
    <button className={classes} style={style} onClick={onClick}>
      <div className={styles.icon}>
        <FontAwesomeIcon icon={icon} />
      </div>
    </button>
  );
};

export default IconButton;
