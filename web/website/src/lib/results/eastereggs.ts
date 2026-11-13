// Small surprises for curious people: a line shown on the result for a few
// special links. Order matters; the first match wins.

import { TAGLINE_LINKS } from '../data/tagline';
import { REPO } from '../site';

type Egg = { test: (u: URL, host: string) => boolean; say: string };

const RICKROLL = 'dQw4w9WgXcQ';
const REPO_URL = new URL(REPO);

const EGGS: Egg[] = [
  {
    test: (_, h) => h === 'url.vet',
    say: '检查我自己？挺大胆。换我也会这么做。',
  },
  {
    test: (u, h) =>
      h === REPO_URL.hostname &&
      u.pathname.toLowerCase().startsWith(REPO_URL.pathname.toLowerCase()),
    say: '那是我住的地方。点个星也不亏。',
  },
  {
    test: (_, h) => h === 'abhizaik.com',
    say: '那是把我做出来的人。看起来挺可靠，不过我有偏心。',
  },
  {
    test: (u, h) =>
      ((h === 'youtube.com' || h === 'm.youtube.com') && u.searchParams.get('v') === RICKROLL) ||
      (h === 'youtu.be' && u.pathname === `/${RICKROLL}`),
    say: '点开是安全的。你想不想点就是另一回事了。Never gonna give you up。',
  },
  {
    // Lookalikes spelled with Cyrillic letters, like the examples on the home page.
    test: (_, h) => /[Ѐ-ӿ]/.test(h) || h.startsWith('xn--'),
    say: '看出差别了吗？肉眼根本看不出来，这正是它的把戏。',
  },
  {
    test: (_, h) => TAGLINE_LINKS.includes(h),
    say: '这个来自我们首页的标语。是我们编的，不过早晚会有人真的注册它。',
  },
  {
    test: (_, h) => h === 'paypal.com',
    say: '这是真的那个。假的会把「а」换成西里尔字母；这个没有。',
  },
  {
    test: (_, h) => h === 'apple.com',
    say: '这是真的那个。假的会拼成 аррӏе，每个字母都换成了西里尔字母的孪生兄弟。',
  },
  {
    test: (_, h) => h === 'google.com',
    say: '啊对，就是你用来搜索网站的那个网站。',
  },
  {
    test: (_, h) => ['example.com', 'example.org', 'example.net'].includes(h),
    say: '互联网上最无辜的网站。它的存在就是为了让别的页面可以指向它。',
  },
  {
    test: (_, h) => h === 'wikipedia.org' || h.endsWith('.wikipedia.org'),
    say: '需要引用来源？这个不用。',
  },
  {
    test: (_, h) => h === 'neverssl.com',
    say: '一个故意不用 HTTPS 的网站，好让 Wi-Fi 登录页能加载。这是「不安全」反而成了重点的少数情况。',
  },
  {
    test: (_, h) => ['1.1.1.1', '1.0.0.1', '8.8.8.8', '8.8.4.4', '9.9.9.9'].includes(h),
    say: '那是一个 DNS 解析器：互联网的电话簿，其实不算网站。',
  },
  {
    test: (_, h) => h.endsWith('.zip') || h.endsWith('.mov'),
    say: '一个不是文件的「.zip」。现在的它们是真实的网址，骗子也注意到了。',
  },
];

function parse(raw: string): { url: URL; host: string } | null {
  try {
    const url = new URL(/^https?:\/\//i.test(raw) ? raw : `http://${raw}`);
    const host = url.hostname.toLowerCase().replace(/^www\./, '');
    return { url, host };
  } catch {
    return null;
  }
}

/** A quip for a special link, or null. */
export function easterEgg(raw: string): string | null {
  const p = parse(raw);
  if (!p) return null;
  return EGGS.find((e) => e.test(p.url, p.host))?.say ?? null;
}

/** A friendly reason not to scan addresses that only exist on the visitor's side. */
export function localAddressMessage(raw: string): string | null {
  const p = parse(raw);
  if (!p) return null;
  const h = p.url.hostname.toLowerCase().replace(/^\[|\]$/g, '');
  if (
    h === 'localhost' ||
    h.endsWith('.localhost') ||
    h === '::1' ||
    /^127\./.test(h) ||
    h === '0.0.0.0'
  ) {
    return '那是你自己的电脑。我们从这里看不到它，别人也看不到。';
  }
  if (
    /^10\./.test(h) ||
    /^192\.168\./.test(h) ||
    /^172\.(1[6-9]|2\d|3[01])\./.test(h) ||
    h.endsWith('.local')
  ) {
    return '那是私有网络里的地址，比如家用路由器。从互联网无法访问它，所以我们没什么可检查的。';
  }
  return null;
}
