// Verdict wording shared by the results page and the social preview image,
// so both always say the same thing.

export type Verdict = 'Safe' | 'Suspicious' | 'Risky';

export const VERDICT_COPY: Record<Verdict, { name: string; label: string; quip: string }> = {
  Safe: {
    name: '安全',
    label: '可信',
    quip: '看起来没问题，放心点就行。 👌',
  },
  Suspicious: {
    name: '可疑',
    label: '需谨慎',
    quip: '嗯，有点不对劲，小心一点。 🤔',
  },
  Risky: {
    name: '危险',
    label: '高风险',
    quip: '很可能不安全，最好别点。 🚫',
  },
};

export function verdictCopy(verdict: string | undefined) {
  return VERDICT_COPY[(verdict ?? '') as Verdict] ?? VERDICT_COPY.Suspicious;
}

/** Plain-Chinese name for the verdict value itself, like 安全 or 危险. */
export function verdictName(verdict: string | undefined): string {
  return verdictCopy(verdict).name;
}
