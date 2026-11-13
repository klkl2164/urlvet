import type { RedirectInfo, ShortLinkInfo } from '../types';

function host(url: string): string {
  try {
    return new URL(url).hostname;
  } catch {
    return url;
  }
}

/** One sentence saying where a short link led, or null for an ordinary link. */
export function shortLinkNote(link: ShortLinkInfo | undefined): string | null {
  if (!link) return null;
  const via = host(link.url);
  if (!link.resolved || !link.target) {
    return `这是 ${via} 上的一个短链接。我们看不到它指向哪里，所以这个结果只针对短链接本身，最高也只能判定为可疑。`;
  }
  // Short links can hand off to other short links; name each one on the way.
  const hops = [...new Set(link.chain.slice(1).map(host))].filter(
    (h) => h !== via && h !== host(link.target!)
  );
  const through = hops.length ? `（途经 ${hops.join('、')}）` : '';
  return `这是 ${via} 上的一个短链接。它指向 ${host(link.target)}${through}，这个结果是针对那一页的。`;
}

/** One sentence for a link a well-known site redirected elsewhere (an open redirect), or null. */
export function redirectNote(info: RedirectInfo | undefined): string | null {
  if (!info?.target) return null;
  return `${host(info.url)} 上的这个链接会把访客转发到 ${host(info.target)}。这个结果是针对那一页的，也包括链接本身存在的问题。`;
}
