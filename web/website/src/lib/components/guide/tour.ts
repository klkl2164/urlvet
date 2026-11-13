// Guided tour steps. Each points at an element tagged data-guide="<target>";
// steps whose target isn't on the page right now are skipped.

/** landing: before a scan. results: once a result is on screen. */
export type TourStep = {
  target: string;
  title: string;
  text: string;
  on: 'landing' | 'results';
  /** Room around the highlight, in px (default 8); less where neighbours are close. */
  pad?: number;
};

export const TOUR_STEPS: TourStep[] = [
  {
    target: 'search',
    on: 'landing',
    title: '在这里粘贴链接',
    text: '完整链接、短链接，或者只有 example.com 这样的域名都行。也可以粘贴一整条消息，我会从中找出链接。',
  },
  {
    target: 'paste',
    on: 'landing',
    pad: 2, // Scan sits right beside it
    title: '或者一键粘贴',
    text: '粘贴你上次复制的内容。',
  },
  {
    target: 'scan',
    on: 'landing',
    title: '然后点「扫描」',
    text: '按回车也行。检查需要几秒钟。',
  },
  {
    target: 'examples',
    on: 'landing',
    title: '只是想试试？',
    text: '这些会检查真实网站。红色的是用同形字符拼出来的假域名。',
  },
  {
    target: 'verdict',
    on: 'results',
    title: '判定结果',
    text: '安全、可疑或危险，附带一个满分 100 的分数。分数越高越安全。',
  },
  {
    target: 'screenshot',
    on: 'results',
    title: '页面长什么样',
    text: '一张页面截图，让你不用亲自打开就能看到。',
  },
  {
    target: 'flags',
    on: 'results',
    title: '为什么',
    text: '红色标记是让我担心的地方，绿色标记是让我放心的地方。',
  },
  {
    target: 'sections',
    on: 'results',
    title: '全部细节',
    text: '展开任意一行，就能看到我具体检查了什么。右边的小圆点是简要结论。',
  },
  {
    target: 'report',
    on: 'results',
    title: '觉得我弄错了？',
    text: '告诉我们。每一份报告都有真人看。',
  },
  {
    target: 'share',
    on: 'results',
    title: '分享结果',
    text: '把它发给把链接发给你的人，让他们也能看到。',
  },
  {
    target: 'learn-more',
    on: 'landing',
    title: '想了解更多？',
    text: '「工作原理」解释了每一项检查。「隐私」说明了检查过的链接会被怎么处理。',
  },
];

/** Steps for the current page (landing or results) whose target is on screen. */
export function availableSteps(): TourStep[] {
  const has = (target: string) => !!document.querySelector(`[data-guide="${target}"]`);
  const page = has('verdict') ? 'results' : 'landing';
  return TOUR_STEPS.filter((s) => s.on === page && has(s.target));
}
