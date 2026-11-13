// Mentions of url.vet shown on the landing page.

export type ShoutoutType = 'tweet' | 'paper' | 'email' | 'package' | 'newsletter' | 'linkedin';

export type Shoutout = {
  name: string;
  handle: string;
  date: string;
  text: string;
  url: string;
  type: ShoutoutType;
  flag: string;
};

export const SHOUTOUTS: Shoutout[] = [
  {
    type: 'paper',
    name: '学术研究',
    handle: 'ijesr.org',
    date: '2025 年 4 月 26 日',
    text: '这个项目的早期版本被一篇学术论文引用并复现（👀）。',
    url: 'https://ijesr.org/index.php/ijesr/article/view/377',
    flag: '🇮🇳',
  },
  {
    type: 'package',
    name: 'Dika Ardianta',
    handle: '@DikaArdnt',
    date: '2026 年 4 月 28 日',
    text: '我把检测引擎移植到了 PHP，并作为非官方的社区 Composer 包发布在 packagist.org 上',
    url: 'https://packagist.org/packages/safesurf/safesurf',
    flag: '🇮🇩',
  },
  {
    type: 'tweet',
    name: 'Tom Dörr',
    handle: '@tom_doerr',
    date: '2026 年 1 月 16 日',
    text: '带网页界面和浏览器扩展的钓鱼检测引擎',
    url: 'https://twitter.com/tom_doerr/status/2012050578721915177',
    flag: '🇩🇪',
  },
  {
    type: 'newsletter',
    name: 'OSINTech',
    handle: 'osintech.substack.com',
    date: '2026 年 5 月 21 日',
    text: 'URLvet。开源钓鱼检测引擎——实时获得信任分、完全可解释的判定结果，以及带页面实时预览、可分享的安全报告。',
    url: 'https://osintech.substack.com/p/osintechs-timeline-163-21052026?open=false#%C2%A7osint-tools-services-and-investigations',
    flag: '🇰🇿',
  },
  {
    type: 'tweet',
    name: 'NeoTeo.com',
    handle: '@NeoteoCom',
    date: '2026 年 1 月 16 日',
    text: 'SafeSurf，一个带网页界面和浏览器扩展的开源钓鱼检测引擎。容易集成，用来保护用户、减少欺诈。',
    url: 'https://twitter.com/NeoteoCom/status/2012102861807554843',
    flag: '🇪🇸',
  },
  {
    type: 'tweet',
    name: 'Bryan',
    handle: '@so_sthbryan',
    date: '2026 年 5 月 16 日',
    text: '可解释的钓鱼检测，实时扫描 URL。SafeSurf 能在你和钓鱼网站交互之前就发现它。',
    url: 'https://x.com/so_sthbryan/status/2055390764339974377',
    flag: '🇺🇸',
  },
  {
    type: 'linkedin',
    name: 'Şevket Can Özöğretmen',
    handle: '@ozogretmen',
    date: '2026 年 5 月 14 日',
    text: '不确定网上的一个链接会把你带到哪里？url.vet 是一个强大的钓鱼检测引擎，实时扫描、评分，并透明地解释一个链接为什么安全或不安全。',
    url: 'https://www.linkedin.com/posts/ozogretmen_cybersecurity-opensource-phishingdetection-ugcPost-7460836352545341440-B4yF/',
    flag: '🇹🇷',
  },
];
