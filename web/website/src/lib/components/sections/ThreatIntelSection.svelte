<script lang="ts">
  import StatusIcon from "../StatusIcon.svelte";
  import type { FeedMatch, PhishingResult, WebRiskResult } from "../../types";
  import TooltipIcon from "../TooltipIcon.svelte";
  export let phishing: PhishingResult | undefined;
  export let feeds: FeedMatch | undefined = undefined;
  export let webRisk: WebRiskResult | undefined = undefined;
  /** Which Google service gave the verdict in webRisk. */
  export let googleSource = "Google Web Risk";

  const THREAT_NAMES: Record<string, string> = {
    SOCIAL_ENGINEERING: "钓鱼",
    MALWARE: "恶意软件",
    UNWANTED_SOFTWARE: "不需要的软件",
    POTENTIALLY_HARMFUL_APPLICATION: "有害应用",
  };
  $: listedBy = feeds?.sources?.join("和") ?? "";
  $: googleThreats = (webRisk?.threat_types ?? []).map((t) => THREAT_NAMES[t] ?? t).join("和");
</script>

{#if phishing || feeds || webRisk}
  <section class="p-4 sm:p-5 space-y-4">
    {#if feeds || webRisk}
      <div
        class="{feeds?.listed || webRisk?.listed
          ? 'border-red-700/60 bg-red-950/20'
          : 'border-gray-300 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/30'} rounded-lg border p-4"
      >
        <div class="flex items-center gap-2 mb-3">
          <span class="text-sm font-semibold text-gray-900 dark:text-white">威胁情报列表</span>
          <TooltipIcon
            text="已上报的钓鱼和恶意软件链接列表，在我们的服务器上保持更新，每次扫描都会比对。"
          />
        </div>
        {#if feeds}
          {#if !feeds.listed}
            <p
              class="text-sm text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-2"
            >
              <StatusIcon kind="ok" />
              {feeds.checked?.length
                ? `未被 ${feeds.checked.join(" 或 ")} 收录。`
                : "不在我们持有的任何钓鱼或恶意软件列表中。"}
            </p>
          {:else if feeds.match === "url"}
            <p class="text-sm font-semibold text-red-400 flex items-center gap-2">
              <StatusIcon kind="bad" /> 此链接已被 {listedBy} 收录。
            </p>
          {:else}
            <p class="text-sm font-semibold text-red-400 flex items-center gap-2">
              <StatusIcon kind="bad" /> 此站点的其他页面已被 {listedBy} 收录。
            </p>
          {/if}
        {/if}
        {#if webRisk}
          <p
            class="text-sm font-medium flex items-center gap-2 mt-2 {webRisk.listed
              ? 'text-red-400 font-semibold'
              : 'text-emerald-700 dark:text-emerald-400'}"
          >
            {#if webRisk.listed}
              <StatusIcon kind="bad" /> Google 将此链接标记为 {googleThreats || "危险"}。
            {:else}
              <StatusIcon kind="ok" /> 未被 {googleSource} 收录。
            {/if}
          </p>
          {#if webRisk.listed}
            <!-- Required by the Safe Browsing and Web Risk terms wherever Google's verdict is shown. -->
            <p class="mt-2 text-xs leading-snug text-gray-500 dark:text-gray-400">
              <a
                href="https://developers.google.com/safe-browsing/v4/advisory"
                target="_blank"
                rel="noopener noreferrer"
                class="underline underline-offset-2">安全提示由 Google 提供</a
              >。Google 致力于提供关于不安全网络资源最准确、最新的信息。但 Google
              无法保证这些信息全面无误：有些危险网站可能没有被识别出来，也有些安全网站可能被误报。
            </p>
          {/if}
        {/if}
      </div>
    {/if}

    {#if phishing}
      <div
        class="{phishing.in_database && phishing.valid
          ? 'border-red-700/60 bg-red-950/20'
          : 'border-gray-300 dark:border-gray-700 bg-gray-50 dark:bg-gray-800/30'} rounded-lg border p-4"
      >
        <!-- Block header -->
        <div class="flex items-center gap-2 mb-3">
          <img src="https://phishtank.com/favicon.ico" alt="" class="w-4 h-4" />
          <span class="text-sm font-semibold text-gray-900 dark:text-white">PhishTank 查询</span>
          <TooltipIcon text="比对由社区核实的钓鱼数据库，每次扫描都会运行。" />
        </div>

        {#if !phishing.in_database}
          <!-- Not in database at all — cleanest result -->
          <p
            class="text-sm text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-2"
          >
            <StatusIcon kind="ok" /> 未在 PhishTank 数据库中找到。
          </p>
        {:else if phishing.verified && !phishing.valid}
          <!-- In database, reviewed, and confirmed NOT phishing -->
          <p
            class="text-sm text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-2"
          >
            <StatusIcon kind="ok" /> 已由 PhishTank 社区审核，确认不是钓鱼。
          </p>
          {#if phishing.phish_id}
            <p class="text-xs text-gray-500 mt-1">
              报告
              {#if phishing.phish_detail_page}
                <a
                  href={phishing.phish_detail_page}
                  target="_blank"
                  rel="noopener noreferrer"
                  class="text-gray-400 hover:text-gray-600 dark:text-gray-300 underline"
                  >#{phishing.phish_id}</a
                >
              {:else}
                #{phishing.phish_id}
              {/if}
            </p>
          {/if}
        {:else}
          <!-- In database and valid=true (phishing) or not yet reviewed -->
          <div
            class="space-y-0 divide-y divide-gray-100 dark:divide-gray-700/50 text-sm text-gray-800 dark:text-gray-200"
          >
            <!-- in_database -->
            <div
              class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2 first:pt-0"
            >
              <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                <span>在数据库中：</span>
                <TooltipIcon text="这个 URL 已被提交或举报到 PhishTank。" />
              </div>
              <span class="text-amber-700 dark:text-amber-400 font-medium">是，已有记录</span>
            </div>

            <!-- valid — the actual phishing signal -->
            <div
              class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2"
            >
              <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                <span>是否为钓鱼：</span>
                <TooltipIcon
                  text="这个 URL 是否被确认为钓鱼网站。true = 是钓鱼，false = 不是钓鱼。"
                />
              </div>
              {#if phishing.valid}
                <span class="font-semibold text-red-400">是，已确认钓鱼</span>
              {:else}
                <span class="font-semibold text-gray-400">否</span>
              {/if}
            </div>

            <!-- verified — reviewed status -->
            <div
              class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2"
            >
              <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                <span>社区已审核：</span>
                <TooltipIcon
                  text="PhishTank 社区是否已审核这份报告（true = 已审核，false = 仍待审核）。"
                />
              </div>
              {#if phishing.verified}
                <span class="text-gray-900 dark:text-white font-medium">是</span>
              {:else}
                <span class="text-amber-700 dark:text-amber-400 font-medium">否，等待审核</span>
              {/if}
            </div>

            <!-- verified_at -->
            {#if phishing.verified_at}
              <div
                class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2"
              >
                <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                  <span>审核时间：</span>
                  <TooltipIcon text="PhishTank 社区审核这份报告的时间。" />
                </div>
                <span class="font-medium text-gray-800 dark:text-white">{phishing.verified_at}</span
                >
              </div>
            {/if}

            <!-- target -->
            {#if phishing.target}
              <div
                class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2"
              >
                <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                  <span>冒充目标：</span>
                  <TooltipIcon text="这个钓鱼 URL 冒充的品牌或服务。" />
                </div>
                <span class="font-medium text-gray-800 dark:text-white">{phishing.target}</span>
              </div>
            {/if}

            <!-- phish_id + phish_detail_page -->
            {#if phishing.phish_id}
              <div
                class="flex flex-col md:grid md:grid-cols-[220px,1fr] md:items-center gap-1 md:gap-4 py-2 last:pb-0"
              >
                <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                  <span>PhishTank 报告：</span>
                  <TooltipIcon text="唯一的 PhishTank ID。点击可查看完整报告页面。" />
                </div>
                {#if phishing.phish_detail_page}
                  <a
                    href={phishing.phish_detail_page}
                    target="_blank"
                    rel="noopener noreferrer"
                    class="font-mono text-sm text-gray-900 dark:text-gray-100 underline underline-offset-4 decoration-gray-300 dark:decoration-gray-700 hover:decoration-current"
                    >#{phishing.phish_id}</a
                  >
                {:else}
                  <span class="font-mono text-sm text-gray-600 dark:text-gray-300"
                    >#{phishing.phish_id}</span
                  >
                {/if}
              </div>
            {/if}
          </div>
        {/if}
      </div>
    {/if}
  </section>
{/if}
