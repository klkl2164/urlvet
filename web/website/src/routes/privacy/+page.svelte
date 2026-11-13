<script lang="ts">
  import PageCta from "$lib/components/PageCta.svelte";
  import { REPO } from "$lib/site";
  import { LINK } from "$lib/ui/text";

  // Plain-language privacy notice. Keep it in sync with what the server actually does:
  // logging (server/internal/logger, handler/middleware/request_log.go), rate limiting
  // (handler/middleware/rate_limit.go), caching (server/internal/constants/ttls.go),
  // reports (handler/report.go).

  // The short version: what we don't have, and the few things we keep for a while.
  const NONE = [
    "账号",
    "Cookie",
    "广告",
    "追踪器",
    "关于你的用户画像",
    "日志里的 IP 地址",
    "出售的数据",
  ];
  // What we keep, how long for, and what's in it.
  const KEEP = [
    {
      t: "扫描结果",
      for: "24 小时",
      d: "按链接缓存，所以重复检查或者分享出去的检查都能立刻加载出来。部分单项检查会更早过期。",
    },
    {
      t: "页面截图",
      for: "24 小时",
      d: "保存在我们的服务器上，之后自动删除。",
    },
    {
      t: "最近的搜索",
      for: "直到重启",
      d: "最近 100 次扫描（链接、判定结果、分数、时间），保存在服务器内存里供管理后台查看。从不写入磁盘。",
    },
    {
      t: "你发来的报告",
      for: "直到问题修复",
      d: "链接、我们的判定结果和分数、你的建议、你的备注以及时间。没有 IP 地址，没有账号。",
    },
  ];

  const THIRD_PARTIES = [
    {
      who: "Vercel，我们的托管服务商",
      what: "托管这个网站，所以它会看到你的 IP 地址和你打开的页面（如果是分享链接，这也包括被扫描的链接）。它会短暂保留请求日志。",
      href: "https://vercel.com/legal/privacy-policy",
      linkText: "查看 Vercel 的隐私政策",
    },
    {
      who: "你扫描的网站",
      what: "我们的服务器会访问它并截图。这个网站以及它上面的任何追踪器看到的都是我们的服务器，而不是你。",
    },
    {
      who: "PhishTank",
      what: "通常什么都收不到。我们每隔几小时下载一次它收录的钓鱼举报列表，并在自己的服务器上拿你的链接去比对。只有在这份列表还没加载成功时，我们才会把链接发给 PhishTank。",
    },
    {
      who: "Google Safe Browsing",
      what: "只会收到由链接算出的一小段短码，绝不会收到链接本身。很多不同的链接会共用同一段短码，所以 Google 分辨不出你检查的是哪一个。它会从自己的危险网站列表里返回匹配结果，检查在我们服务器上完成。",
      href: "https://policies.google.com/privacy",
      linkText: "查看 Google 的隐私政策",
    },
    {
      who: "DNS 解析器",
      what: "只会收到域名，用来找到域名背后的服务器。",
    },
    {
      who: "域名注册局",
      what: "只会收到域名，用来查询它的注册时间。",
    },
  ];
</script>

<svelte:head>
  <title>隐私 | url.vet</title>
  <meta
    name="description"
    content="我们检查链接，不检查你。没有账号、Cookie 或分析工具，不出售任何数据，日志里也没有 IP 地址。"
  />
  <meta property="og:title" content="隐私 | url.vet" />
  <meta
    property="og:description"
    content="我们检查链接，不检查你。没有账号、Cookie 或分析工具，不出售任何数据，日志里也没有 IP 地址。"
  />
  <meta property="og:type" content="website" />
  <meta property="og:url" content="https://url.vet/privacy" />
  <meta property="og:image" content="https://url.vet/og" />
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:title" content="隐私 | url.vet" />
  <meta
    name="twitter:description"
    content="我们检查链接，不检查你。没有账号、Cookie 或分析工具，不出售任何数据，日志里也没有 IP 地址。"
  />
  <meta name="twitter:image" content="https://url.vet/og" />
  <link rel="canonical" href="https://url.vet/privacy" />
</svelte:head>

<div class="max-w-3xl mx-auto px-6 pt-16 md:pt-24 pb-20 text-gray-900 dark:text-gray-100">
  <p class="font-mono text-xs uppercase tracking-wider text-gray-500">隐私</p>
  <h1 class="mt-4 font-serif text-5xl md:text-6xl leading-[1.02] tracking-[-0.015em]">
    我们检查链接，<span class="italic">不检查你。</span>
  </h1>
  <p class="mt-6 text-xl leading-relaxed text-gray-700 dark:text-gray-300">
    url.vet 只需要你粘贴的链接，别的都不要。下面就是它具体会经历什么。
  </p>

  <!-- At a glance: what we don't have, as one statement; what we keep, as three figures -->
  <section class="mt-14" aria-label="概述">
    <p class="font-mono text-[11px] uppercase tracking-wider text-gray-500">简单说</p>
    <p class="mt-4 font-serif text-[1.75rem] md:text-[2.1rem] leading-[1.3] tracking-[-0.01em]">
      {#each NONE as item}
        <span class="whitespace-nowrap"
          ><span class="text-gray-400 dark:text-gray-500">没有</span>{item}。</span
        >{" "}
      {/each}
    </p>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">当你扫描一个链接</h2>
    <div class="mt-5 space-y-4 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      <p>
        链接会发到我们的服务器，服务器执行各项检查，再把结果发回来。在响应请求的过程中，我们的服务器会
        看到你的 IP 地址，只把它用于一件事：限制同一个地址每分钟能发多少请求。它只保留一个无法还原的
        加密摘要，并且一分钟后就把这个摘要删掉。
      </p>
      <p>
        <span class="text-gray-900 dark:text-gray-100 font-medium"
          >我们的日志里永远不会包含你的 IP 地址，也不会包含你扫描的链接。</span
        >
        每个请求只记录一行：访问的是哪个页面、是否成功、用了多长时间。如果某条错误信息里带有网址，写入之前会把
        <span class="font-mono text-[15px]">?</span> 之后的内容全部截掉。
      </p>
    </div>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">我们保留什么，保留多久</h2>
    <dl class="mt-6">
      {#each KEEP as item}
        <div
          class="flex flex-col sm:flex-row gap-1 sm:gap-6 py-5 border-t border-gray-200 dark:border-gray-800"
        >
          <dt class="sm:w-44 shrink-0">
            <span class="block font-serif text-[1.3rem] leading-tight">{item.t}</span>
            <span class="mt-1 block font-mono text-[11px] uppercase tracking-wider text-gray-500"
              >{item.for}</span
            >
          </dt>
          <dd class="text-[16px] leading-relaxed text-gray-700 dark:text-gray-300">{item.d}</dd>
        </div>
      {/each}
    </dl>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">还有谁会看到链接</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      有些检查需要外部协助。以下是仅有的其他参与方：
    </p>
    <dl class="mt-6">
      {#each THIRD_PARTIES as p}
        <div
          class="flex flex-col sm:flex-row gap-1 sm:gap-6 py-5 border-t border-gray-200 dark:border-gray-800"
        >
          <dt class="sm:w-44 shrink-0 font-serif text-[1.3rem] leading-tight">{p.who}</dt>
          <dd class="text-[16px] leading-relaxed text-gray-700 dark:text-gray-300">
            {p.what}
            {#if p.href}
              <a href={p.href} target="_blank" rel="noopener noreferrer" class={LINK}
                >{p.linkText}</a
              >.
            {/if}
          </dd>
        </div>
      {/each}
    </dl>
    <p class="mt-4 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      这个网站加载的所有内容（字体、图片、代码）都来自 url.vet 自己。
    </p>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">分享结果</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      分享链接里包含被扫描的链接、判定结果和分数，所以你发给的人（以及任何显示它的预览的应用）都能看到
      这些内容。里面没有关于你的任何信息。
    </p>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">在你的浏览器里</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      你的浏览器会记住你选的是浅色还是深色主题，以及你是否见过或者隐藏了
      Vetty。这些只保存在你的设备上，清除网站数据就会把它们删掉。
    </p>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">浏览器扩展</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      扩展会在你打开页面时检查它，只发送页面的地址，和在这里扫描的方式一样。最知名的 10000
      个网站永远不会
      被发送：这份列表保存在你的设备上。你的设置和信任的网站由浏览器保存，如果你用了浏览器同步，也会由浏览器
      同步；url.vet 永远收不到这些信息。你可以在扩展的设置里关掉自动检查。
    </p>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">你的权利</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      url.vet 由 <a
        href="https://abhizaik.com"
        target="_blank"
        rel="noopener noreferrer"
        class={LINK}>abhizaik</a
      >运营。我们不保存任何与你姓名绑定的信息。想了解我们保存了哪些关于你扫描或报告过的链接的信息，或者想
      删除某条报告，请发邮件到
      <a href="mailto:hi@url.vet" class="font-mono text-base {LINK}">hi@url.vet</a>。
    </p>
  </section>

  <section class="mt-20">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">自己验证</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      上面说的每一条都在<a href={REPO} target="_blank" rel="noopener noreferrer" class={LINK}
        >源代码</a
      >里，如果你不想用我们的版本，也可以自己跑一份。这一页的每一次改动也都是公开的，在
      <a
        href="{REPO}/commits/main/web/website/src/routes/privacy/+page.svelte"
        target="_blank"
        rel="noopener noreferrer"
        class={LINK}>GitHub 上的提交历史</a
      >里可以查到。
    </p>
  </section>

  <PageCta />
</div>
