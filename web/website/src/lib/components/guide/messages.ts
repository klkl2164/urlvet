// Everything Vetty says, in one place. Tone: the casual url.vet voice,
// with a wink at the old office assistant.
import type { AnalyzeResult } from '../../types';

export const NAME = 'Vetty';

export const INTRO = [
  `你好，我是 ${NAME}。既然能直接问我，何必自己到处点？`,
  '你需要的东西都在这里。',
];

/** Said when someone keeps tapping Vetty, one line per extra tap. */
export const POKES = ['嗯？', '我就在这儿。', '你对所有回形针都这么动手动脚吗？', '好吧，有点痒。'];

export const NUDGE = '你好，我是 Vetty。看不懂这些是什么意思？我可以帮忙。';
/** The visitor's two replies to the nudge. */
export const NUDGE_YES = '你能做什么？';
export const NUDGE_NO = '不用了，谢谢';

/** Greeting when reopening Vetty, by page. */
export const GREETINGS: Record<string, string[]> = {
  home: [
    '看起来你正在检查一个链接。需要帮忙吗？',
    '又回来了？很好。',
    '需要帮忙吗？我手不多，但都归你用。',
  ],
  result: ['要我带你看看这个结果吗？', '细则我已经读完了，你不用读。有什么就问吧。'],
  scanning: ['等等，我正在读细则……'],
  error: ['这次没成功。谁都会遇到。'],
  'how-it-works': ['在看说明文档？佩服。', '啊，又一个爱看文档的人。'],
  about: ['来了解我们？我是比较友好的那个。', '没错，我也在这儿工作。'],
  other: ['迷路了？我知道一个检查链接的好地方。'],
};

export const TIPS = [
  '小提示：在任意位置按 / 就能跳到搜索框。',
  '小提示：短链接会隐藏真正的目的地。粘贴到这里，就能看到它一路去哪儿。',
  '小提示：一个页面可以长得和你的银行一模一样，但它并不是你的银行。URL 也一样。',
  '小提示：每个结果都有分享链接，你可以把它发回给把链接发给你的人。',
];

/** Score bands, matching the backend thresholds. */
export const SCORE_BANDS = [
  {
    verdict: 'Risky',
    from: 0,
    to: 29,
    dot: 'bg-red-500',
    text: 'text-red-600 dark:text-red-400',
    meaning: '有明显的钓鱼或诈骗迹象。离远一点。',
  },
  {
    verdict: 'Suspicious',
    from: 30,
    to: 64,
    dot: 'bg-yellow-500',
    text: 'text-yellow-700 dark:text-yellow-400',
    meaning: '信号有好有坏。登录和付款要小心。',
  },
  {
    verdict: 'Safe',
    from: 65,
    to: 100,
    dot: 'bg-emerald-500',
    text: 'text-emerald-700 dark:text-emerald-400',
    meaning: '没有发现让人担心的地方。',
  },
];

export const SCORES_FOOTNOTE = '我也可能出错，所以每个结果都有「报告」按钮。每一份报告都有真人看。';

export const SCAN_FAILED = [
  '扫描失败通常有三个原因：',
  '• 地址打错了，',
  '• 网站挂了，',
  '• 或者它拦截了自动访问（有些银行会这样）。',
  '再检查一遍链接，然后重试。',
];

export const HIDE_BYE = '好吧。我去页脚待着。反正我以前也被收起来过。';

/** "Who is Vetty?": what he is, what he can do, and why he's here. */
export const MEET = {
  lead: `我是 ${NAME}，url.vet 的小助手。`,
  why: '检查一个链接不该需要安全学位。我在这里，是为了让任何人都能用 url.vet：我用大白话解释结果，还会告诉你每个东西在哪儿。',
  canTitle: '我能做什么',
  can: [
    '帮你检查一个链接，或者一整条带链接的消息',
    '用大白话解释结果',
    '告诉你分数是什么意思',
    '带你逛一遍页面',
    '帮你报告看起来不对的结果',
  ],
  note: '我只记得你在这里检查过的链接，你离开后我就忘了。随时可以把我收起来，我会在页脚等着。',
};

/** A message to send back to whoever shared the link, worded for the verdict. */
export function warnMessage(verdict: string, domain: string, score: number, link: string): string {
  if (verdict === 'Risky')
    return `我在 url.vet 上检查了这个链接，看起来是诈骗（${domain}，${score}/100）。请不要打开它，也不要在那里填写任何信息。\n\n原因：${link}`;
  if (verdict === 'Suspicious')
    return `我在 url.vet 上检查了这个链接。${domain} 看起来可疑（${score}/100），最好别在那里填写任何信息。\n\n原因：${link}`;
  return `我在 url.vet 上检查了这个链接，${domain} 看起来是安全的（${score}/100）。\n\n详情：${link}`;
}

export const WARN_INTRO: Record<string, string> = {
  Risky: '我们提醒一下发链接的人吧。想改的话可以编辑这条消息：',
  Suspicious: '在他们打开之前先提醒一下。想改的话可以编辑这条消息：',
  Safe: '告诉他们这个没问题。想改的话可以编辑这条消息：',
};

export const REAL_SITE = {
  brand: (d: string) => `这个页面在冒充 ${d}。`,
  lookalike: (d: string) => `这个地址是故意做得像 ${d} 的。`,
  advice: (d: string) =>
    `如果你本来想访问 ${d}，别用别人发给你的链接。直接去真正的网站，或者先检查一下。`,
};

export const CLEAN_LINK = {
  lead: '这个链接带着跟踪参数。它们会告诉发送者谁点了、从哪里点的。',
  done: '同一个页面，去掉这些参数：',
};

export const CHECK_FOR_ME = '把链接粘进来，或者把整条消息粘进来。我会把它放进搜索框并检查。';

export const TRY_EXAMPLE = '挑一个吧。绿色的是真实链接，红色的是用西里尔字母拼出来的假链接。';

/** Plain-language explanation of a scan result. */
export function explain(result: AnalyzeResult): {
  lead: string;
  lists: { title: string; items: string[] }[];
  advice: string;
} {
  const verdict = result.result?.verdict ?? 'Suspicious';
  const score = result.result?.final_score ?? 0;
  const bad = (result.result?.reasons?.bad_reasons ?? []).slice(0, 3);
  const good = (result.result?.reasons?.good_reasons ?? []).slice(0, 3);
  const domain = result.domain;

  if (verdict === 'Safe') {
    return {
      lead: `好消息：${domain} 看起来是真的（${score}/100）。`,
      lists: [
        { title: '我为什么放心', items: good },
        { title: '我注意到的小问题', items: bad },
      ].filter((l) => l.items.length),
      advice: '不过，如果这个链接来自一条你没预料到的消息，自己手动输入网址比直接点它更安全。',
    };
  }
  if (verdict === 'Risky') {
    return {
      lead: `我建议离 ${domain} 远一点（${score}/100）。`,
      lists: [{ title: '让我担心的地方', items: bad }].filter((l) => l.items.length),
      advice: '别在那里填任何东西。如果你觉得我弄错了，报告一下，会有真人来核查。',
    };
  }
  return {
    lead: `对 ${domain} 我拿不准（${score}/100）。`,
    lists: [
      { title: '让我担心的地方', items: bad },
      { title: '对它有利的地方', items: good },
    ].filter((l) => l.items.length),
    advice: '如果这个页面让你登录或付款，别在这里做。自己在浏览器里输入网址，再从那里登录。',
  };
}

export const pick = <T>(list: T[]): T => list[Math.floor(Math.random() * list.length)];
