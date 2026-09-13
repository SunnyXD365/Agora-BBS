export function formatReadingDuration(totalSeconds: number): string {
  const seconds = Math.max(0, Math.floor(totalSeconds));
  if (seconds < 60) return `${seconds} 秒`;
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes} 分钟`;
  const remainingMinutes = minutes % 60;
  return remainingMinutes > 0 ? `${Math.floor(minutes / 60)} 小时 ${remainingMinutes} 分钟` : `${Math.floor(minutes / 60)} 小时`;
}
