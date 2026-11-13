// The check groups on the How it works page: what each one looks for, and whether a
// finding counts towards risk or trust. The rules live in
// server/internal/analyzer/result.go; no point values here, since those get tuned.
// Result sections link here by id (#check-url, #check-network, ...).

export type Signal = {
  text: string;
  /** Which side of the score it counts towards; 'info' is shown but not scored. */
  kind: 'risk' | 'trust' | 'info';
};

export type CheckGroup = {
  id: string;
  /** The question this group answers, in plain words. */
  title: string;
  /** The technical name, for people who know it. */
  label: string;
  /** How many of the scan's checks belong to this group. */
  checks: number;
  desc: string;
  signals: Signal[];
  /** What the findings are compared against. */
  against?: string;
};

const risk = (text: string): Signal => ({ text, kind: 'risk' });
const trust = (text: string): Signal => ({ text, kind: 'trust' });
const info = (text: string): Signal => ({ text, kind: 'info' });

export const CHECK_GROUPS: CheckGroup[] = [
  {
    id: 'url',
    title: '地址本身看起来有问题吗？',
    label: 'URL 结构',
    checks: 8,
    desc: '在访问任何东西之前先看链接本身。很多骗局光靠地址就能暴露。',
    signals: [
      risk('直接用 IP 地址，而不是域名'),
      risk('编码过的仿冒字符（Punycode）'),
      risk('用其他字母表里外形相同的字母冒充'),
      risk('未受监管的域名后缀'),
      risk('层级非常深的路径'),
      risk('短链接，我们没能追到它指向哪里'),
      risk('高风险后缀，例如 .top 或 .zip'),
      risk('非常长的链接'),
      risk('超过两级子域名'),
      risk('链接里含有「login」或「verify」之类的词'),
      trust('经过验证的后缀，例如 .gov、.edu 或 .bank'),
      info('短链接会被跟进，途中遇到的其他短链接也继续跟，最终检查它们指向的页面'),
    ],
    against: '138 个高风险后缀、238 个已验证后缀和 2,671 个已知的短链接服务。',
  },
  {
    id: 'network',
    title: '链接实际通向哪里？',
    label: '重定向与 HTTP',
    checks: 1,
    desc: '从我们的服务器发出一次真实请求，跟随每一次重定向，直到最终页面。',
    signals: [
      risk('最终落在与链接所示不同的网站'),
      risk('超过三次跳到站外的重定向'),
      trust('网站强制使用安全连接（HSTS）'),
      info('网站返回的响应码'),
    ],
  },
  {
    id: 'dns',
    title: '这个网站的配置像正常网站吗？',
    label: 'DNS',
    checks: 2,
    desc: '查询网站托管在哪里。一次性的诈骗网站常常连基本配置都省掉。',
    signals: [
      risk('名称服务器配置缺失或不完整'),
      risk('域名没有邮件服务器'),
      info('网站解析到的地址'),
    ],
  },
  {
    id: 'tls',
    title: '连接是否安全？',
    label: 'TLS / SSL',
    checks: 1,
    desc: '读取网站的安全证书。证书有效只能证明连接是私密的，不能证明网站是正当的，所以这里大部分内容只是展示给你看。',
    signals: [
      risk('证书签发给了别的地址（浏览器会就此警告）'),
      risk('页面没有 HTTPS，却要求输入密码'),
      info('证书由谁签发，以及何时到期'),
      info('证书是否被公开记录'),
    ],
  },
  {
    id: 'domain',
    title: '这个网站有多久的历史、多有名？',
    label: '域名情报',
    checks: 4,
    desc: '正规网站会积累历史。诈骗网站通常只注册了几天，没人知道，名字还仿照别人。',
    signals: [
      risk('与知名网站只差一两个字母'),
      risk('最近 30 天内注册'),
      risk('别人的地址里含有某个品牌的名称'),
      risk('最近 90 天内注册'),
      risk('不在访问量前一百万名之列，而且很新或注册时间不明'),
      trust('属于访问量最高的 10,000 个网站'),
      trust('属于访问量最高的 50,000 个网站'),
      trust('有排名，但访问量较低'),
      trust('已运行超过一年，三年、五年以上更佳'),
      trust('DNS 记录有签名（DNSSEC）'),
      info('域名名称看起来有多随机'),
      info(
        '托管在托管服务或建站平台上的网站（github.io、vercel.app、weebly.com）会单独评估，不借助平台本身的信誉'
      ),
    ],
    against:
      '访问量排名依据前一百万个网站，仿冒拼写依据访问量最高的 5,000 个网站，地址中的品牌名依据 514 个品牌。',
  },
  {
    id: 'content',
    title: '页面要求你提供什么？',
    label: '页面内容',
    checks: 1,
    desc: '抓取页面并读取其中的表单：表单要你填什么，以及这些信息会被发往哪里。',
    signals: [
      risk('页面自称是某个品牌，但该品牌并不拥有这个地址'),
      risk('表单把密码、银行卡信息或个人资料发往另一个网站'),
      risk('新建或不知名网站上的登录表单'),
      risk('隐藏的框架在加载另一个页面'),
      risk('银行卡或支付字段'),
      risk('表单把邮箱或用户名发往另一个网站'),
      risk('Cloudflare 或托管商已标记或下架该页面'),
      risk('页面用脚本把你转走，转到一个裸 IP 地址或另一个网站'),
      risk('藏在 WordPress 系统目录里的页面，这是网站被入侵的迹象'),
      risk('从不知名的网站下载程序，或者从任何人都能上传文件的网站下载程序'),
      trust('页面提到的品牌确实拥有这个地址'),
      info('追踪像素，以及对个人资料的索取'),
    ],
    against:
      '覆盖银行、支付、加密货币、购物、快递、电信、政府等领域的 514 个品牌，以及它们实际使用的地址。仿冒字符会按它所模仿的普通字符来比对。',
  },
  {
    id: 'threats',
    title: '有没有人已经举报过它？',
    label: '威胁情报源',
    checks: 4,
    desc: '检查已被举报为钓鱼或恶意软件的链接名单：PhishTank、URLhaus 和 Google Safe Browsing。这些名单保存在我们的服务器上，Google 只会收到由链接生成的短代码，永远看不到链接本身。',
    signals: [
      risk('出现在已确认的钓鱼或恶意软件链接名单上'),
      risk('同一网站的其他页面出现在这些名单上'),
      risk('Google 把它列为危险'),
    ],
  },
];

export const TOTAL_CHECKS = CHECK_GROUPS.reduce((n, g) => n + g.checks, 0);
