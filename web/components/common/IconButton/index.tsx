import { IconProp } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classNames from 'classnames';
import React from 'react';
import styles from './styles.module.css';

export interface IconButtonProps {
  color?: string;
  icon: IconProp;
  onClick: () => void;
}

const IconButton: React.FC<IconButtonProps> = ({ color, icon, onClick }) => {
  const classes = classNames(styles.container, styles[color]);
  return (
    <button className={classes} onClick={onClick}>
      <div className={styles.icon}>
        <FontAwesomeIcon icon={icon} />
      </div>
    </button>
  );
};

export default IconButton;
