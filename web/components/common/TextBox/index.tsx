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
  onPaste?: (evt: React.ClipboardEvent) => void;
}

const TextBox: React.FC<TextBoxProps> = ({
  value,
  placeholder,
  onChange,
  onPaste,
}) => {
  const [filled, setFilled] = useState(false);
  const textRef = useRef<HTMLDivElement>();
  const valueRef = useRef('');

  const handleFilled = (value: string) => {
    const nextFilled = !!value;
    if (filled !== nextFilled) setFilled(nextFilled);
  };

  const handleValue = (value: string) => {
    const oldValue = valueRef.current;
    if (oldValue === value) return;
    valueRef.current = value;
    handleFilled(value);
    if (onChange) onChange(value);
  };

  const handleChange = useCallback(() => {
    const value = textRef.current.innerText.trim();
    handleValue(value);
  }, [handleValue]);

  const handlePaste = useCallback(
    (e: React.ClipboardEvent) => {
      e.preventDefault();
      e.stopPropagation();
      const value = e.clipboardData.getData('text/plain');
      document.execCommand('inserttext', false, value);
      handleValue(value);
      if (onPaste) onPaste(e);
    },
    [handleChange]
  );

  const textPlaceholderClasses = useMemo(
    () => classNames(styles.textPlaceholder, { [styles.hidden]: filled }),
    [filled]
  );

  useEffect(() => {
    const oldValue = valueRef.current;
    if (oldValue !== value) {
      textRef.current.innerText = value;
      handleValue(value);
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
        onPaste={handlePaste}
      />
      <div className={textPlaceholderClasses}>{placeholder}</div>
    </div>
  );
};

export default TextBox;
