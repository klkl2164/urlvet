<script lang="ts">
  import { vettyMemory } from "$lib/components/guide/store";
  import PageCta from "$lib/components/PageCta.svelte";
  import { REPO } from "$lib/site";
  import { PILL_OUTLINE } from "$lib/ui/buttons";
  import { GITHUB_LOGO } from "$lib/ui/icons";
  import { LINK } from "$lib/ui/text";
  import { tick } from "svelte";

  const VETTY = [
    {
      title: "Vetty 会做什么",
      dot: "bg-emerald-500",
      items: [
        "告诉你怎么使用 url.vet",
        "帮你检查一个链接",
        "用简单的话解释检查结果",
        "如果页面是假的，带你去找真正的网站",
      ],
    },
    {
      title: "Vetty 绝不会做什么",
      dot: "bg-red-500",
      items: [
        "收集你的数据",
        "用你输入的内容训练 AI",
        "要求你注册或付费",
        "展示广告，或者向你推销任何东西",
      ],
    },
  ];

  // Opens Vetty's bubble, bringing him back first if he was hidden.
  async function sayHi() {
    vettyMemory.set({ hidden: false });
    await tick();
    if (!document.querySelector(".vetty[aria-expanded='true']")) {
      document.querySelector<HTMLElement>("button.vetty")?.click();
    }
  }

  const PRINCIPLES = [
    {
      title: "结果可解释",
      desc: "一个解释不了的判定结果，就没法让人相信。每一个信号都会展示出来，还会说明它是怎么引出最终判定结果的。",
    },
    {
      title: "为所有人而做",
      desc: "对安全从业者来说足够详细，对你的父母来说也足够通俗。同一次扫描，你想看多深就看多深。",
    },
    {
      title: "完全开源",
      desc: "安全工具就应该让人检查。检测引擎、评分逻辑和网站全部放在 GitHub 上，可以随便阅读、审计，也可以自己托管。",
    },
    {
      title: "实时检查",
      desc: "钓鱼页面几小时内就会出现又消失。url.vet 在你发起请求的那一刻才检查这个链接，所以能抓到黑名单还没来得及收录的页面。",
    },
    {
      title: "隐私优先",
      desc: "不需要账号，没有广告，没有追踪器，不卖数据。我们不记录你是谁，也不记录你扫描了什么。",
    },
  ];
</script>

<svelte:head>
  <title>关于 url.vet</title>
  <meta
    name="description"
    content="url.vet 告诉你一个链接是否安全，以及为什么。免费、开源，并且从设计上保护隐私。"
  />
  <meta property="og:title" content="关于 url.vet" />
  <meta
    property="og:description"
    content="url.vet 告诉你一个链接是否安全，以及为什么。免费、开源，并且从设计上保护隐私。"
  />
  <meta property="og:type" content="website" />
  <meta property="og:url" content="https://url.vet/about" />
  <meta property="og:image" content="https://url.vet/og" />
  <meta name="twitter:card" content="summary_large_image" />
  <meta name="twitter:title" content="关于 url.vet" />
  <meta
    name="twitter:description"
    content="url.vet 告诉你一个链接是否安全，以及为什么。免费、开源，并且从设计上保护隐私。"
  />
  <meta name="twitter:image" content="https://url.vet/og" />
  <link rel="canonical" href="https://url.vet/about" />
</svelte:head>

<div class="max-w-3xl mx-auto px-6 pt-16 md:pt-24 pb-20 text-gray-900 dark:text-gray-100">
  <p class="font-mono text-xs uppercase tracking-wider text-gray-500">关于</p>
  <h1 class="mt-4 font-serif text-5xl md:text-6xl leading-[1.02] tracking-[-0.015em]">
    为我爱的人而做。<span class="italic sm:block">对所有人开放。</span>
  </h1>
  <p class="mt-6 text-xl leading-relaxed text-gray-700 dark:text-gray-300">
    url.vet 是一个免费的开源工具，告诉你一个链接是否安全，以及为什么。这一页讲的是它为什么存在，
    以及它坚持什么。
  </p>

  <!-- Maker: the why, in first person -->
  <section class="mt-20">
    <h2 class="mt-2 font-serif text-3xl md:text-4xl tracking-[-0.01em]">我为什么做它</h2>
    <!-- A letter: straight into the words, signed once at the end. -->
    <figure
      class="relative mt-8 rounded-2xl border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 px-6 py-7 md:px-10 md:py-9 shadow-xs shadow-black/[0.03]"
    >
      <span
        class="block h-9 font-serif text-[4.5rem] leading-none text-accent-light/40 dark:text-accent-dark/40 select-none"
        aria-hidden="true">&ldquo;</span
      >
      <blockquote
        class="mt-3 max-w-[36rem] space-y-5 text-[17px] md:text-lg leading-relaxed text-gray-800 dark:text-gray-200"
      >
        <p>
          我经常收到一个链接，然后琢磨它打开之后安不安全。我试过几个在线链接扫描工具。有的只给一句
          「安全」或者「危险」，别的什么都没有。有的则丢给我好几页技术数据，让我自己去翻。这两种我都
          不太满意。
        </p>
        <p>
          后来我发现，同样的问题在家里也有。我父母会打开、转发 WhatsApp 上发来的链接，其实根本不知道
          真假（而且常常是假的）。朋友会转给我一些「超值优惠」，点进去根本不是他们说的那家店。
        </p>
        <p>
          另外，我是个挺在意隐私的人，也不喜欢现在大多数应用对待用户的方式：什么都追踪，把你的数据全
          收集起来，再卖给广告商。这些我都觉得很烦，所以想做一个完全不这么干的东西。
        </p>
        <p>
          我想要一个家里人、朋友和我自己都能用的东西。给他们一个直接的答案，也给想看细节的人一份背后
          的理由。
        </p>
        <p class="font-medium text-gray-900 dark:text-gray-100">url.vet 就是这么开始的。</p>
      </blockquote>

      <!-- Sign-off: photo and name, both linking to the creator's site. -->
      <figcaption class="mt-8">
        <a
          href="https://abhizaik.com"
          target="_blank"
          rel="noopener noreferrer"
          class="group inline-flex items-center gap-4"
        >
          <img
            src="/team/abhizaik.jpg"
            alt=""
            class="w-12 h-12 rounded-full shrink-0 grayscale-[30%] group-hover:grayscale-0 transition"
            loading="lazy"
          />
          <span>
            <span
              class="block font-signature text-[1.75rem] leading-none text-gray-900 dark:text-gray-100"
              >abhizaik</span
            >
            <span class="mt-1.5 block text-[13px] text-gray-500">url.vet 的作者</span>
          </span>
        </a>
      </figcaption>
    </figure>
  </section>

  <!-- Principles -->
  <section class="mt-24">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">理念</h2>
    <ul class="mt-8">
      {#each PRINCIPLES as item}
        <li
          class="py-5 border-t border-gray-200 dark:border-gray-800 sm:grid sm:grid-cols-[13rem_1fr] sm:gap-8"
        >
          <p class="font-serif text-[1.45rem] leading-tight">{item.title}</p>
          <p class="mt-2 sm:mt-1 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">
            {item.desc}
          </p>
        </li>
      {/each}
    </ul>
  </section>

  <!-- Vetty -->
  <section class="mt-24" id="vetty">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">认识 Vetty</h2>
    <p class="mt-5 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      Vetty 是角落里那个回形针。他帮你使用 url.vet。他的灵感来自
      Clippy，旧版办公软件里的回形针助手。
    </p>

    <div class="mt-8 grid sm:grid-cols-2 gap-x-10 gap-y-8">
      {#each VETTY as list}
        <div class="pt-4 border-t border-gray-200 dark:border-gray-800">
          <h3 class="font-serif text-[1.45rem] leading-tight">{list.title}</h3>
          <ul class="mt-3 space-y-2 text-[15px] leading-relaxed text-gray-600 dark:text-gray-400">
            {#each list.items as item}
              <li class="flex gap-2.5">
                <span
                  class="mt-[0.55em] w-1.5 h-1.5 rounded-full shrink-0 {list.dot}"
                  aria-hidden="true"
                ></span>
                <span>{item}</span>
              </li>
            {/each}
          </ul>
        </div>
      {/each}
    </div>

    <p class="mt-8 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      不想让他在旁边？从他的菜单里把他藏起来就好。他会在页面底部等着。
    </p>
    <button type="button" class="mt-5 {PILL_OUTLINE} px-5 py-2.5" on:click={sayHi}>
      跟 Vetty 打个招呼
    </button>
  </section>

  <!-- Open source + help -->
  <section class="mt-24">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">
      全部，<span class="italic">都是开放的。</span>
    </h2>
    <div class="mt-5 space-y-4 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      <p>
        很多项目自称开源，其实只公开了网站。在这里，检测引擎、评分逻辑、API 和网站全部以 <a
          href="{REPO}/blob/main/LICENSE"
          target="_blank"
          rel="noopener noreferrer"
          class={LINK}>AGPL-3.0</a
        > 协议公开。你可以阅读、审计，或者自己托管。
      </p>
      <p>
        想帮忙？可以从一个 <a
          href="{REPO}/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22"
          target="_blank"
          rel="noopener noreferrer"
          class={LINK}>适合新手的 issue</a
        >
        开始，或者先读一下
        <a
          href="{REPO}/blob/main/.github/CONTRIBUTING.md"
          target="_blank"
          rel="noopener noreferrer"
          class={LINK}>贡献指南</a
        >。
      </p>
    </div>
    <a
      href={REPO}
      target="_blank"
      rel="noopener noreferrer"
      class="mt-7 {PILL_OUTLINE} px-5 py-2.5"
    >
      <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 24 24" aria-hidden="true">
        <path d={GITHUB_LOGO} />
      </svg>
      在 GitHub 上查看
    </a>
  </section>

  <!-- Contact -->
  <section class="mt-24">
    <h2 class="font-serif text-3xl md:text-4xl tracking-[-0.01em]">联系我们</h2>
    <div class="mt-5 space-y-4 text-[17px] leading-relaxed text-gray-700 dark:text-gray-300">
      <p>
        有反馈、有想法，或者只是想打个招呼？发邮件到 <a
          href="mailto:hi@url.vet"
          class="font-mono text-base {LINK}">hi@url.vet</a
        >。
      </p>
      <p>
        觉得某个判定结果不对？在任何结果上点「报告」是最快的办法。每一份报告都会有人看。发现了安全问题？
        请按照<a
          href="{REPO}/blob/main/.github/SECURITY.md"
          target="_blank"
          rel="noopener noreferrer"
          class={LINK}>安全策略</a
        >来处理。
      </p>
      <p class="text-[15px] text-gray-500">url.vet 也写作 URLvet。它以前叫 SafeSurf。</p>
    </div>
  </section>

  <PageCta />
</div>
