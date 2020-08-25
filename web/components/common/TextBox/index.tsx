import classNames from 'classnames';
import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react';
import styles from './styles.module.css';

export interface TextBoxProps {
  value: string;
  placeholder?: string;
  onChange?: (value: string) => void;
}

const TextBox: React.FC<TextBoxProps> = ({ value, placeholder, onChange }) => {
  const [filled, setFilled] = useState(false);
  const textRef = useRef<HTMLDivElement>();
  const valueRef = useRef('');

  const handleFilled = (value: string) => {
    const nextFilled = !!value;
    console.log(value, filled, nextFilled);
    if (filled !== nextFilled) setFilled(nextFilled);
  };

  const handleChange = useCallback(() => {
    const oldValue = valueRef.current;
    const value = textRef.current.innerText;
    if (oldValue === value) return;
    valueRef.current = value;
    handleFilled(value);
    if (onChange) onChange(value);
  }, [onChange]);

  const textPlaceholderClasses = useMemo(
    () => classNames(styles.textPlaceholder, { [styles.hidden]: filled }),
    [filled]
  );

  useEffect(() => {
    const oldValue = valueRef.current;
    if (oldValue !== value) {
      textRef.current.innerText = value;
      handleFilled(value);
    }
  }, [value]);

  return (
    <div className={styles.textContainer}>
      <div
        ref={textRef}
        role='textbox'
        className={styles.textEditor}
        contentEditable
        onInput={handleChange}
        onBlur={handleChange}
      />
      <div className={textPlaceholderClasses}>{placeholder}</div>
    </div>
  );
};

export default TextBox;
