<script lang="ts">
  import { env } from "$env/dynamic/public";
  import Icon from "$lib/components/Icon.svelte";
  import PageCta from "$lib/components/PageCta.svelte";
  import { CHECK_GROUPS, TOTAL_CHECKS } from "$lib/data/checks";
  import { FAQ_GROUPS, faqId, faqText } from "$lib/data/faq";
  import { REPO } from "$lib/site";
  import { ICON } from "$lib/ui/icons";
  import { LINK } from "$lib/ui/text";
  import { onMount } from "svelte";

  // Swagger UI is served by the API server itself, next to /api/v1.
  const API_ORIGIN = (() => {
    try {
      return new URL(env.PUBLIC_BASE_URL || "http://localhost:8080/api/v1").origin;
    } catch {
      return "http://localhost:8080";
    }
  })();
  // Self-hosting comes first; the public API is only for a quick try.
  const DEV_LINKS = [
    {
      label: "在自己的服务器上运行",
      note: "Docker 部署、参数配置和反向代理，一步步来",
      href: `${REPO}/blob/main/docs/deployment.md`,
    },
    {
      label: "API 参考",
      note: "analyze 接口、它返回什么，以及可能出现的错误",
      href: `${REPO}/blob/main/docs/api.md`,
    },
    {
      label: "在 Swagger UI 里试用 API",
      note: "在浏览器里向 url.vet 的服务器发几个测试请求。有频率限制，不适合生产环境",
      href: `${API_ORIGIN}/swagger/index.html`,
    },
  ];

  // A made-up scan, worded exactly as the engine words its findings.
  const EXAMPLE = {
    domain: "paypa1-login.top",
    score: 7,
    red: [
      "检测到仿冒域名：该域名与「paypal.com」高度相似（相差 1 个字符）。",
      "注册于 3 天前。钓鱼网站通常都是新注册的。",
      "检测到高风险域名后缀（常与垃圾信息关联）。",
      "URL 中发现敏感的安全关键词：login",
      "访问量极低。",
    ],
    green: ["强制使用严格的 HTTPS 安全设置（已启用 HSTS）。"],
  };

  const STEPS = [
    {
      title: "粘贴链接",
      desc: "一个链接，或者一整条带链接的消息。",
    },
    {
      title: "我们的服务器访问它",
      desc: "页面上的任何内容都不会到达你的设备。",
    },
    {
      title: "你得到判定结果和理由",
      desc: "安全、可疑或危险，每条理由都会列出来。",
    },
  ];

  const BANDS = [
    { v: "危险", from: 0, to: 29, dot: "bg-red-500" },
    { v: "可疑", from: 30, to: 64, dot: "bg-yellow-500" },
    { v: "安全", from: 65, to: 100, dot: "bg-emerald-500" },
  ];

  // What a result has besides the verdict and flags shown in the example.
  const PARTS = [
    {
      title: "页面预览",
      desc: "由我们的服务器拍摄的页面截图。如果它看起来像你的银行，地址却和你的银行毫无关系，这一点就能说明很多。",
    },
    {
      title: "详细信息",
      desc: "按类别分组的区块（域名、重定向、安全证书、页面内容、威胁情报源），给想看原始事实的人。",
    },
    {
      title: "分享链接",
      desc: "每个结果都有一个链接，你可以发给当初给你原文的人。它的预览里一眼就能看到判定结果和分数。",
    },
  ];

  // A link to one question (#faq-...) or one check group (#check-...) opens it.
  function openFromHash() {
    const target = /^#(faq|check)-/.test(location.hash)
      ? document.getElementById(location.hash.slice(1))
      : null;
    const details =
      target instanceof HTMLDetailsElement ? target : target?.querySelector("details");
    if (details) {
      details.open = true;
      target?.scrollIntoView({ block: target === details ? "center" : "start" });
    }
  }
  onMount(openFromHash);

  const schemaFAQ = {
    "@context": "https://schema.org",
    "@type": "FAQPage",
    mainEntity: FAQ_GROUPS.flatMap((g) => g.items).map((f) => ({
      "@type": "Question",
      name: f.q,
      acceptedAnswer: { "@type": "Answer", text: faqText(f) },
    })),
  };
</script>

<svelte:window on:hashchange={openFromHash} />

<svelte:head>
  <title>url.vet 如何检查一个链接</title>
  <meta
    name="description"
    content="url.vet 对每个链接执行的检查、它们如何汇总成一个判定结果，以及常见问题的解答。"
  />
  <meta property="og:title" content="url.vet 如何检查一个链接" />
  <meta
    property="og:description"
    content="url.vet 对每个链接执行的检查、它们如何汇总成一个判定结果，以及常见问题的解答。"
  />
  <meta property="og:type" content="website" />
  <meta property="og:url" content="https://url.vet/how-it-works" />
  <meta property="og:image" content="https://url.vet/og" />
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:title" content="url.vet 如何检查一个链接" />
  <meta
    name="twitter:description"
    content="url.vet 对每个链接执行的检查、它们如何汇总成一个判定结果，以及常见问题的解答。"
  />
  <meta name="twitter:image" content="https://url.vet/og" />
  <link rel="canonical" href="https://url.vet/how-it-works" />
  {@html `<script type="application/ld+json">${JSON.stringify(schemaFAQ)}</script>`}
</svelte:head>

<div class="max-w-3xl mx-auto px-6 pt-16 md:pt-24 pb-20 text-gray-900 dark:text-gray-100">
  <p class="font-mono text-xs uppercase tracking-wider text-gray-500">工作原理</p>
  <h1 class="mt-4 font-serif text-5xl md:text-6xl leading-[1.02] tracking-[-0.015em]">
    粘贴一个链接。<span class="italic">看清为什么。</span>
  </h1>
  <p class="mt-6 text-xl leading-relaxed text-gray-700 dark:text-gray-300">
    解释不了的判定结果就没法信任，所以每一个判定结果都带着理由。下面是你发起扫描后发生的事，
    以及怎么看懂返回的内容。
  </p>

  <!-- The whole thing in three steps: a row on desktop, a short timeline on phones -->
  <ol class="mt-12 grid sm:grid-cols-3 sm:gap-x-6">
    {#each STEPS as step, n}
      <li class="relative flex sm:block gap-4 pb-7 last:pb-0 sm:pb-0">
        <!-- line to the next step -->
        {#if n < STEPS.length - 1}
          <span
            class="absolute bg-gray-200 dark:bg-gray-800 left-[17px] top-11 bottom-1 w-px sm:left-12 sm:-right-3 sm:top-[17px] sm:bottom-auto sm:w-auto sm:h-px"
            aria-hidden="true"
          ></span>
        {/if}
        <span
          class="relative shrink-0 flex items-center justify-center w-9 h-9 rounded-full border border-gray-300 dark:border-gray-700 font-mono text-[13px] text-gray-700 dark:text-gray-300"
          aria-hidden="true">{n + 1}</span
        >
        <div class="sm:mt-4 sm:pr-4">
          <p class="font-serif text-[1.35rem] leading-tight">{step.title}</p>
          <p class="mt-1.5 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">
            {step.desc}
          </p>
        </div>
      </li>
    {/each}
  </ol>

  <!-- Example result -->
  <section class="mt-20 scroll-mt-20" id="result">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">你会得到什么</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      url.vet
      检查的是页面本身，而且是实时检查，不只是比对一份已知的坏链接名单。你按下「扫描」几秒钟后，
      就会看到类似这样的东西：
    </p>

    <figure
      class="mt-8 rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 p-5 sm:p-6"
      aria-label="示例结果"
    >
      <div class="flex flex-wrap items-center justify-between gap-3">
        <span
          class="px-3 py-1 rounded-full border border-gray-300 dark:border-gray-800 font-mono text-[13px] text-gray-800 dark:text-gray-200"
          >{EXAMPLE.domain}</span
        >
        <span class="flex items-baseline gap-3">
          <span
            class="inline-flex items-center gap-1.5 text-sm font-medium text-red-600 dark:text-red-400"
            ><span class="w-1.5 h-1.5 rounded-full bg-red-500"></span>危险</span
          >
          <span class="font-serif text-3xl leading-none"
            >{EXAMPLE.score}<span class="text-base text-gray-400">/100</span></span
          >
        </span>
      </div>

      {#each [{ label: "风险标记", items: EXAMPLE.red, dot: "bg-red-500" }, { label: "正面标记", items: EXAMPLE.green, dot: "bg-emerald-500" }] as group}
        <p class="mt-6 font-mono text-[11px] uppercase tracking-wider text-gray-500">
          {group.label}
        </p>
        <ul class="mt-2 divide-y divide-gray-100 dark:divide-gray-800">
          {#each group.items as flag}
            <li class="flex gap-3 py-2.5 text-[15px] leading-snug text-gray-700 dark:text-gray-300">
              <span class="mt-[0.45em] w-1.5 h-1.5 shrink-0 rounded-full {group.dot}"></span>
              <span class="flex-1">{flag}</span>
            </li>
          {/each}
        </ul>
      {/each}
      <figcaption class="mt-4 font-mono text-[11px] text-gray-400">
        示例结果，已简化。并非真实网站。
      </figcaption>
    </figure>

    <dl class="mt-10 grid sm:grid-cols-3 gap-x-8 gap-y-7">
      {#each PARTS as part}
        <div class="pt-4 border-t border-gray-200 dark:border-gray-800">
          <dt class="font-serif text-[1.35rem] leading-tight">{part.title}</dt>
          <dd class="mt-2 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">
            {part.desc}
          </dd>
        </div>
      {/each}
    </dl>
  </section>

  <!-- Checks: what each group looks at, and which side of the score a finding counts towards -->
  <section class="mt-24 scroll-mt-20" id="checks">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">检查项</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      你按下「扫描」的同时，{TOTAL_CHECKS} 项检查会一起开始。每一项都是独立的，所以某一项变慢或者失败，
      不会拖住其他检查。它们合起来回答七个问题，每个发现都会计入分数中的风险一侧或信任一侧。
    </p>
    <p class="mt-3 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">
      对于知名网站和 .gov
      这类经过验证的域名后缀，少数基于名称的规则会放宽，这样真正的银行自己的登录页
      就不会被误标。想知道扫描过程中哪些内容会被发到哪里？请看<a href="/privacy" class={LINK}
        >隐私页面</a
      >。
    </p>

    <div class="mt-8">
      {#each CHECK_GROUPS as group, n}
        <article
          id="check-{group.id}"
          class="scroll-mt-20 py-7 border-t border-gray-200 dark:border-gray-800"
        >
          <p class="font-mono text-[11px] uppercase tracking-wider text-gray-500">
            {String(n + 1).padStart(2, "0")} · {group.label} · {group.checks} 项检查
          </p>
          <h3 class="mt-2 font-serif text-[1.6rem] leading-tight">{group.title}</h3>
          <p class="mt-2 text-[16px] leading-relaxed text-gray-700 dark:text-gray-300">
            {group.desc}
          </p>

          <details class="group mt-3">
            <summary
              class="inline-flex items-center gap-1.5 cursor-pointer list-none [&::-webkit-details-marker]:hidden text-[14px] text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-100 transition-colors focus:outline-hidden focus-visible:ring-2 focus-visible:ring-gray-400 rounded-sm"
            >
              <span
                class="underline underline-offset-4 decoration-gray-300 dark:decoration-gray-700"
                >查看它检查的内容（{group.signals.length}）</span
              >
              <span class="transition-transform duration-200 group-open:rotate-180">
                <Icon path={ICON.chevronDown} class="w-3.5 h-3.5" />
              </span>
            </summary>
            <ul class="mt-4 grid sm:grid-cols-2 gap-x-10">
              {#each group.signals as signal}
                <li
                  class="flex items-baseline justify-between gap-4 py-2 border-b border-gray-100 dark:border-gray-800/70 text-[14px]"
                >
                  <span class="flex items-baseline gap-2.5 text-gray-700 dark:text-gray-300">
                    <span
                      class="relative top-[-1px] w-1.5 h-1.5 rounded-full shrink-0 {signal.kind ===
                      'risk'
                        ? 'bg-red-500'
                        : signal.kind === 'trust'
                          ? 'bg-emerald-500'
                          : 'bg-gray-300 dark:bg-gray-600'}"
                      aria-hidden="true"
                    ></span>
                    {signal.text}
                  </span>
                  <span
                    class="shrink-0 font-mono text-[11px] whitespace-nowrap {signal.kind === 'risk'
                      ? 'text-red-600 dark:text-red-400'
                      : signal.kind === 'trust'
                        ? 'text-emerald-700 dark:text-emerald-400'
                        : 'text-gray-400 dark:text-gray-500'}"
                  >
                    {signal.kind === "risk"
                      ? "计入风险"
                      : signal.kind === "trust"
                        ? "计入信任"
                        : "仅作展示"}
                  </span>
                </li>
              {/each}
            </ul>

            {#if group.against}
              <p class="mt-3 text-[13px] text-gray-500">对照 {group.against}</p>
            {/if}
          </details>
        </article>
      {/each}
    </div>
  </section>

  <!-- Scoring: the idea in words; no point values, since those get tuned -->
  <section class="mt-24 scroll-mt-20" id="score">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">分数是怎么算的</h2>
    <div class="mt-5 space-y-4 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      <p>
        每个发现都会计入两个方向之一：风险或信任。分数从中间开始，也就是 50 分。风险会把它往下拉，
        信任会把它往上抬，分数始终保持在 0 到 100 之间。
      </p>
      <p>
        严重的发现比轻微的发现权重更大，所以一条确认的钓鱼举报，抵得过好几条轻微的良好迹象。
        算出来的数字会落在三个区间之一。上面这个例子停在 {EXAMPLE.score} 分。
      </p>
    </div>

    <!-- scale -->
    <div class="mt-10 mb-2">
      <div class="relative">
        <div class="flex h-2 rounded-full overflow-hidden gap-0.5">
          {#each BANDS as band}
            <div class={band.dot} style="width: {band.to - band.from + 1}%"></div>
          {/each}
        </div>
        <div
          class="absolute -top-6 flex flex-col items-center"
          style="left: {EXAMPLE.score}%; transform: translateX(-50%)"
        >
          <span class="font-mono text-[10px] text-gray-500 leading-none">{EXAMPLE.score}</span>
          <span class="mt-1 w-0.5 h-6 rounded-full bg-gray-900 dark:bg-gray-100"></span>
        </div>
      </div>
      <div class="mt-4 grid grid-cols-3 gap-4">
        {#each BANDS as band}
          <div>
            <p class="flex items-center gap-2 text-[15px] font-medium">
              <span class="w-1.5 h-1.5 rounded-full {band.dot}"></span>{band.v}
            </p>
            <p class="mt-0.5 font-mono text-xs text-gray-500">{band.from}–{band.to}</p>
          </div>
        {/each}
      </div>
    </div>
  </section>

  <!-- FAQ: collapsed, one open at a time; every answer stays in the page for search engines -->
  <section class="mt-24 scroll-mt-20" id="faq">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">常见问题</h2>
    <div class="mt-8">
      {#each FAQ_GROUPS as group}
        <div class="border-t border-gray-200 dark:border-gray-800">
          <h3 class="pt-6 font-mono text-[11px] uppercase tracking-wider text-gray-500">
            {group.title}
          </h3>
          <div class="mt-2 divide-y divide-gray-200 dark:divide-gray-800">
            {#each group.items as item}
              <details id={faqId(item.q)} name="faq" class="group scroll-mt-24">
                <summary
                  class="flex items-start justify-between gap-4 py-4 cursor-pointer list-none [&::-webkit-details-marker]:hidden text-gray-900 dark:text-gray-100 hover:text-gray-600 dark:hover:text-gray-300 transition-colors focus:outline-hidden focus-visible:ring-2 focus-visible:ring-gray-400 rounded-sm"
                >
                  <span class="font-serif text-[1.3rem] leading-snug">{item.q}</span>
                  <span
                    class="mt-1.5 shrink-0 text-gray-400 transition-transform duration-200 group-open:rotate-180"
                  >
                    <Icon path={ICON.chevronDown} />
                  </span>
                </summary>
                <div class="pb-5 pr-8 text-[16px] leading-relaxed text-gray-600 dark:text-gray-400">
                  <p>{item.a}</p>
                  {#if item.points}
                    <ul class="mt-3 space-y-2">
                      {#each item.points as point}
                        <li class="flex gap-2.5">
                          <span
                            class="mt-[0.6em] w-1 h-1 rounded-full bg-gray-400 shrink-0"
                            aria-hidden="true"
                          ></span>
                          <span>{point}</span>
                        </li>
                      {/each}
                    </ul>
                  {/if}
                </div>
              </details>
            {/each}
          </div>
        </div>
      {/each}
    </div>
  </section>

  <!-- For developers -->
  <section class="mt-24 scroll-mt-20" id="developers">
    <p class="font-mono text-[11px] uppercase tracking-wider text-gray-500">面向开发者</p>
    <h2 class="mt-2 font-serif text-3xl md:text-4xl tracking-[-0.01em]">自托管与 API 文档</h2>
    <p class="mt-4 text-[16px] leading-relaxed text-gray-600 dark:text-gray-400">
      url.vet 是<a href={REPO} class={LINK} target="_blank" rel="noopener noreferrer">开源的</a
      >。你可以把整套东西跑在自己的服务器上，也可以在自己的应用里用它的 API。
    </p>
    <ul class="mt-6">
      {#each DEV_LINKS as link}
        <li class="border-t border-gray-200 dark:border-gray-800">
          <a
            href={link.href}
            target="_blank"
            rel="noopener noreferrer"
            class="group flex items-center justify-between gap-4 py-4"
          >
            <span>
              <span class="block font-serif text-[1.35rem] leading-tight">{link.label}</span>
              <span class="mt-1 block text-[15px] text-gray-600 dark:text-gray-400"
                >{link.note}</span
              >
            </span>
            <span
              class="text-gray-400 group-hover:text-gray-900 dark:group-hover:text-gray-100 group-hover:translate-x-0.5 group-hover:-translate-y-0.5 transition-transform"
              aria-hidden="true">↗</span
            >
          </a>
        </li>
      {/each}
    </ul>
  </section>

  <PageCta />
</div>
