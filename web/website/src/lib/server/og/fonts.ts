// Loads the site's fonts (Instrument Serif, Geist, Geist Mono and a Noto Sans
// SC subset for Chinese) for satori. Without a browser User-Agent, the Google
// Fonts CSS API serves plain TTF files, which satori can read (it cannot read
// WOFF2). Fonts are fetched once and kept in memory.
type Weight = 400 | 500 | 600;
export type OgFont = {
  name: string;
  data: ArrayBuffer;
  weight: Weight;
  style: 'normal' | 'italic';
};

const LATIN_CSS_URL =
  'https://fonts.googleapis.com/css2?family=Instrument+Serif&family=Geist:wght@400;500&family=Geist+Mono:wght@400;500';

// Every Chinese character the cards draw, plus the Latin letters that appear
// inside Chinese sentences (url.vet, URL). A text-only subset keeps the font
// small; satori can't fall back per glyph, so the list has to stay complete.
// Keep it in step with the strings in cards.ts when card wording changes.
const ZH_TEXT =
  '可疑链接？交给就好信任分扫描结果安全危险需谨慎高风险看起来没问题，放心点行嗯有点不对劲小一心很最能别检查任意钓鱼诈骗和重定向源无需注册每条判定都解释url.vetURL ，。、的开';
const ZH_CSS_URL = `https://fonts.googleapis.com/css2?family=Noto+Sans+SC:wght@400;500&text=${encodeURIComponent(ZH_TEXT)}`;

let cache: Promise<OgFont[]> | null = null;

export function loadFonts(): Promise<OgFont[]> {
  cache ??= fetchFonts().catch((err) => {
    cache = null; // retry on the next request
    throw err;
  });
  return cache;
}

async function fetchFonts(): Promise<OgFont[]> {
  const [latin, zh] = await Promise.all([fetchFrom(LATIN_CSS_URL), fetchFrom(ZH_CSS_URL)]);
  return [...latin, ...zh];
}

async function fetchFrom(cssUrl: string): Promise<OgFont[]> {
  const css = await (await fetch(cssUrl)).text();
  const faces = [...css.matchAll(/@font-face\s*{([^}]*)}/g)].map((m) => m[1]);

  return Promise.all(
    faces.map(async (face) => {
      const name = /font-family:\s*'([^']+)'/.exec(face)?.[1] ?? 'Geist';
      const style = /font-style:\s*italic/.test(face) ? 'italic' : 'normal';
      const weight = Number(/font-weight:\s*(\d+)/.exec(face)?.[1] ?? 400) as Weight;
      const url = /url\(([^)]+)\)/.exec(face)?.[1];
      if (!url) throw new Error(`no font url for ${name}`);
      const data = await (await fetch(url)).arrayBuffer();
      return { name, data, weight, style } satisfies OgFont;
    })
  );
}
