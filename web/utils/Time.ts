import moment from 'moment';

const millisecond = 1;
const second = 1000 * millisecond;
const minute = 60 * second;
const hour = 60 * minute;
const day = 24 * hour;
const week = 7 * day;
const month = 30 * day;

interface UnitMap {
  [key: number]: [number, string];
}

const unitMap: UnitMap = {
  [second]: [second, 'second'],
  [minute]: [minute, 'minute'],
  [hour]: [hour, 'hour'],
  [day]: [day, 'day'],
  [week]: [week, 'week'],
  [month]: [month, 'month'],
};

export const formatTimestamp = (ts: number) => {
  return moment(ts).format('YYYY-MM-DD hh:mm:ss');
};

export const formatShortTimestamp = (ts: number, now: number) => {
  const diff = now - ts;
  if (diff <= 0) return 'Just now';
  else if (diff > month) return moment(ts).format('YYYY-MM-DD');
  return moment(ts).from(now, false);
};
