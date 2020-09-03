import clsx from 'clsx';
import React, {
  useCallback,
  useEffect,
  useMemo,
  useRef,
  useState,
} from 'react';
import styles from './styles.module.css';

export interface TextBoxProps {
  className?: string;
  value: string;
  placeholder?: string;
  onChange?: (value: string) => void;
  onPaste?: (evt: React.ClipboardEvent) => void;
}

const TextBox: React.FC<TextBoxProps> = ({
  className,
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
      if (value) {
        document.execCommand('inserttext', false, value);
        handleValue(value);
      }
      if (onPaste) onPaste(e);
    },
    [handleChange, onPaste]
  );

  useEffect(() => {
    const oldValue = valueRef.current;
    if (oldValue !== value) {
      textRef.current.innerText = value;
      handleValue(value);
    }
  }, [value]);

  const containerCls = useMemo(() => clsx(styles.container, className), [
    className,
  ]);

  const placeholderCls = useMemo(
    () => clsx(styles.placeholder, filled && styles.hidden),
    [filled]
  );

  return (
    <div className={containerCls}>
      <div
        ref={textRef}
        role='textbox'
        className={styles.editor}
        contentEditable
        onInput={handleChange}
        onBlur={handleChange}
        onPaste={handlePaste}
      />
      <div className={placeholderCls}>{placeholder}</div>
    </div>
  );
};

export default TextBox;
