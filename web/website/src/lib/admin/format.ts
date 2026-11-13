// Display helpers shared by the admin tabs.

import { verdictCopy } from '../verdict';

/** Chinese display name for a verdict value coming from the API. */
export function verdictLabel(v: string | undefined): string {
  return verdictCopy(v).label;
}

export function formatTTL(sec: number): string {
  if (sec === -1) return '永不过期';
  if (sec === -2) return '已失效';
  if (sec < 60) return `${sec} 秒`;
  if (sec < 3600) return `${Math.floor(sec / 60)} 分 ${sec % 60} 秒`;
  return `${Math.floor(sec / 3600)} 小时 ${Math.floor((sec % 3600) / 60)} 分`;
}

export function ttlColor(sec: number): string {
  if (sec < 0) return 'text-gray-400 dark:text-gray-500';
  if (sec < 60) return 'text-red-500 dark:text-red-400';
  if (sec < 300) return 'text-yellow-600 dark:text-yellow-400';
  return 'text-emerald-600 dark:text-emerald-400';
}

export function verdictColor(v: string): string {
  if (v === 'Safe')
    return 'text-emerald-700 dark:text-emerald-300 bg-emerald-50 dark:bg-emerald-500/10 border-emerald-200 dark:border-emerald-500/20';
  if (v === 'Risky')
    return 'text-red-700 dark:text-red-300 bg-red-50 dark:bg-red-500/10 border-red-200 dark:border-red-500/20';
  return 'text-yellow-700 dark:text-yellow-300 bg-yellow-50 dark:bg-yellow-500/10 border-yellow-200 dark:border-yellow-500/20';
}

export function verdictDot(v: string): string {
  if (v === 'Safe') return 'bg-emerald-500';
  if (v === 'Risky') return 'bg-red-500';
  return 'bg-yellow-500';
}

export function verdictTextColor(v: string): string {
  if (v === 'Safe') return 'text-emerald-600 dark:text-emerald-400';
  if (v === 'Risky') return 'text-red-600 dark:text-red-400';
  return 'text-yellow-600 dark:text-yellow-400';
}

export function relativeTime(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  const s = Math.floor(diff / 1000);
  if (s < 60) return `${s} 秒前`;
  const m = Math.floor(s / 60);
  if (m < 60) return `${m} 分钟前`;
  const h = Math.floor(m / 60);
  if (h < 24) return `${h} 小时前`;
  return `${Math.floor(h / 24)} 天前`;
}

export function formatJSON(val: unknown): string {
  try {
    return JSON.stringify(val, null, 2);
  } catch {
    return String(val);
  }
}

export function hostOf(url: string): string {
  try {
    return new URL(url.includes('://') ? url : `https://${url}`).hostname;
  } catch {
    return url;
  }
}

/** Copies text and calls back with it, then with null after a short delay. */
export async function copyWithFeedback(text: string, onChange: (v: string | null) => void) {
  try {
    await navigator.clipboard.writeText(text);
    onChange(text);
    setTimeout(() => onChange(null), 1200);
  } catch {
    /* clipboard blocked */
  }
}
