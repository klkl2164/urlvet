<script lang="ts">
  import { browser } from "$app/environment";
  import { onDestroy } from "svelte";
  import { REPO } from "../../site";
  import { verdictCopy } from "../../verdict";
  // <head> tags: per-scan share card when a domain is known, the site card otherwise.
  export let shareDomain: string;
  export let verdict: string | undefined;
  export let score: number | undefined;
  export let queryUrl: string;
  export let currentUrl: string;
  /** Tab title for the visitor's own scan (client-side only; shared links keep their SEO title). */
  export let liveTitle: string | undefined = undefined;
  /** Colour of the verdict dot, shown as the tab's icon after a scan. */
  export let liveDot: string | undefined = undefined;

  // Browsers keep showing the first icon they loaded, so swap the existing
  // <link rel="icon"> in place (and put the original back) instead of adding one.
  let originalIcon: string | null = null;
  $: if (browser) {
    const link = document.querySelector<HTMLLinkElement>('link[rel="icon"]');
    if (link) {
      originalIcon ??= link.href;
      link.href = liveDot ? dotIcon(liveDot) : originalIcon;
    }
  }
  onDestroy(() => {
    const link = browser ? document.querySelector<HTMLLinkElement>('link[rel="icon"]') : null;
    if (link && originalIcon) link.href = originalIcon;
  });

  const SITE_TITLE = "url.vet (URLvet)：这个链接安全吗？";
  const SITE_DESC = "粘贴任意链接，查看它是安全、可疑还是危险，以及原因。免费，无需注册。";

  const dotIcon = (color: string) =>
    `data:image/svg+xml,${encodeURIComponent(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 16 16"><circle cx="8" cy="8" r="7" fill="${color}"/></svg>`)}`;

  const schemaSoftwareApp = {
    "@context": "https://schema.org",
    "@type": "SoftwareApplication",
    name: "url.vet",
    alternateName: ["URLvet", "urlvet", "url vet"],
    url: "https://url.vet",
    description:
      "免费的链接检查工具，能识别钓鱼、仿冒地址和其他危险链接，并说明每个判定结果背后的原因。",
    applicationCategory: "SecurityApplication",
    applicationSubCategory: "URL Scanner",
    operatingSystem: "Web",
    browserRequirements: "需要 JavaScript",
    offers: {
      "@type": "Offer",
      price: "0",
      priceCurrency: "USD",
    },
    creator: {
      "@type": "Person",
      name: "abhizaik",
      url: "https://abhizaik.com",
    },
    license: `${REPO}/blob/main/LICENSE`,
    codeRepository: REPO,
    featureList: [
      "实时扫描 URL",
      "钓鱼检测",
      "WHOIS 与域名年龄分析",
      "仿冒域名与同形字符检测",
      "重定向链追踪",
      "TLS 证书验证",
      "IP 与基础设施分析",
      "无需注册",
      "免费使用",
      "开源（AGPL-3.0）",
    ],
  };
</script>

<svelte:head>
  {#if shareDomain}
    {@const title = `${shareDomain} 安全吗？ | url.vet`}
    {@const desc = verdict
      ? `url.vet 将 ${shareDomain} 评为${verdictCopy(verdict).label}。查看原因，或免费检查另一个链接。`
      : `看看 ${shareDomain} 是否可以安全打开，以及原因。免费，无需注册。`}
    {@const ogImage = `https://url.vet/og?domain=${encodeURIComponent(shareDomain)}${verdict ? `&v=${encodeURIComponent(verdict)}` : ""}${score !== undefined ? `&s=${score}` : ""}`}
    <title>{liveTitle ?? title}</title>
    <meta name="description" content={desc} />
    <meta property="og:title" content={title} />
    <meta property="og:description" content={desc} />
    <meta property="og:type" content="website" />
    <link rel="canonical" href="https://url.vet" />
    <meta
      property="og:url"
      content={currentUrl || `https://url.vet/?q=${encodeURIComponent(queryUrl)}`}
    />
    <meta property="og:image" content={ogImage} />
    <meta property="og:image:width" content="1200" />
    <meta property="og:image:height" content="630" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:title" content={title} />
    <meta name="twitter:description" content={desc} />
    <meta name="twitter:image" content={ogImage} />
  {:else}
    <title>{liveTitle ?? SITE_TITLE}</title>
    <meta name="description" content={SITE_DESC} />
    <meta property="og:title" content={SITE_TITLE} />
    <meta property="og:description" content={SITE_DESC} />
    <link rel="canonical" href="https://url.vet" />
    <meta property="og:type" content="website" />
    <meta property="og:url" content="https://url.vet" />
    <meta property="og:image" content="https://url.vet/og" />
    <meta property="og:image:width" content="1200" />
    <meta property="og:image:height" content="630" />
    <meta name="twitter:card" content="summary_large_image" />
    <meta name="twitter:image" content="https://url.vet/og" />
    <meta name="twitter:title" content={SITE_TITLE} />
    <meta name="twitter:description" content={SITE_DESC} />
    {@html `<script type="application/ld+json">${JSON.stringify(schemaSoftwareApp)}</script>`}
  {/if}
</svelte:head>
