import Alert from '@/components/common/Alert';
import '@/styles/main.css';
import { AlertMessage, subscribe } from '@/utils/Alert';
import type { AppProps } from 'next/app';
import Head from 'next/head';
import React, { useCallback, useEffect, useRef, useState } from 'react';

export default function App({ Component, pageProps }: AppProps) {
  const [shown, setShown] = useState(false);
  const [alert, setAlert] = useState({} as AlertMessage);
  const queueRef = useRef<AlertMessage[]>([]);

  const progressRef = useRef(false);
  const shownRef = useRef(shown);

  const workQueue = useCallback((loop: boolean) => {
    if (!loop && progressRef.current) return;
    progressRef.current = true;
    const alert = queueRef.current.pop();
    if (!alert) {
      progressRef.current = false;
      return;
    }
    setAlert(alert);
    if (!shownRef.current) {
      setShown(true);
    }
    setTimeout(() => {
      setShown(false);
      setTimeout(() => workQueue(true), 1000);
    }, 3000);
  }, []);

  useEffect(() => {
    shownRef.current = shown;
  }, [shown]);

  useEffect(() => {
    const observer = subscribe((alert) => {
      queueRef.current.push(alert);
      workQueue(false);
    });
    return () => observer.unsubscribe();
  }, []);

  return (
    <React.Fragment>
      <Head>
        <meta
          name='viewport'
          content='minimum-scale=1, initial-scale=1, width=device-width'
        />
        <title>Air Sync</title>
      </Head>
      <Component {...pageProps} />
      <Alert shown={shown} alert={alert} />
    </React.Fragment>
  );
}
