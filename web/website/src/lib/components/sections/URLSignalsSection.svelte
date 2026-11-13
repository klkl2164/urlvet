<script lang="ts">
  import StatusIcon from "../StatusIcon.svelte";
  import type { DomainRandomness, TyposquatResult } from "../../types";
  import TooltipIcon from "../TooltipIcon.svelte";
  export let features: any;
  export let domainRandomness: DomainRandomness | undefined;
  export let typosquatResult: TyposquatResult | undefined;
</script>

{#if features}
  <section class="p-4 sm:p-5">
    <div
      class="space-y-0 divide-y divide-gray-300 dark:divide-gray-800 text-sm text-gray-800 dark:text-gray-200 max-w-4xl w-full mx-auto"
    >
      {#if features.tld}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>域名后缀（TLD）：</span>
            <TooltipIcon
              text="域名最后的部分（例如 .com、.org、.io），可以暗示网站的信任程度或用途。"
            />
          </div>
          <span class="font-medium text-gray-800 dark:text-white">.{features.tld.tld}</span>
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>高信任域名后缀：</span>
            <TooltipIcon
              text="表示这个域名后缀（TLD）是否被广泛认可，并为政府等高度可信的机构所常用。"
            />
          </div>
          {#if features.tld.is_trusted_tld}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> 是</span
            >
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> 否</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>是否为高风险 TLD：</span>
            <TooltipIcon text="有些 TLD 经常被诈骗者或恶意网站滥用。「是」表示需要谨慎。" />
          </div>
          {#if features.tld.is_risky_tld}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> 是</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> 否</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>TLD 是否被 ICANN 认可：</span>
            <TooltipIcon text="ICANN 负责管理全球域名。被认可意味着这个 TLD 由官方管理和监督。" />
          </div>
          {#if features.tld.is_icann}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> 是</span
            >
          {:else if features.tld.is_hosting_platform}
            <span class="text-yellow-600 dark:text-yellow-400 font-medium flex items-center gap-1.5"
              >托管平台</span
            >
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> 否</span
            >
          {/if}
        </div>
      {/if}

      {#if features.url}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>使用短链接服务：</span>
            <TooltipIcon
              text="显示 URL 是否使用了缩短服务（例如 bit.ly）。短链接会隐藏网站的真实去向。"
            />
          </div>
          {#if features.url.url_shortener}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> 是</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> 否</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>直接使用 IP 地址：</span>
            <TooltipIcon text="有些恶意网站用 IP 地址代替域名来躲避检测。正规网站很少这样做。" />
          </div>
          {#if features.url.uses_ip}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> 是</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> 否</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>包含 Punycode 字符：</span>
            <TooltipIcon
              text="Punycode 让域名可以使用特殊字符或非拉丁字符（例如 xn--example）。诈骗者有时利用它做相似域名的攻击。"
            />
          </div>
          {#if features.url.contains_punycode}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> 是</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> 否</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>URL 过长：</span>
            <TooltipIcon text="超长的 URL 可以用来掩盖恶意内容，或诱使用户信任一个假网站。" />
          </div>
          {#if features.url.too_long}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> 是</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> 否</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>URL 层级过深（斜杠过多）：</span>
            <TooltipIcon text="多层嵌套路径的 URL 可能意味着重定向或隐藏内容，常见于钓鱼尝试。" />
          </div>
          {#if features.url.too_deep}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> 是</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> 否</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>包含相似字母（同形字符）：</span>
            <TooltipIcon
              text="检测域名中是否包含看起来像其他字符的字符（例如 'go0gle.com'），这种手法常用于冒充。"
            />
          </div>
          {#if features.url.has_homoglyph}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> 是</span
            >
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> 否</span
            >
          {/if}
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>子域名数量：</span>
            <TooltipIcon
              text="显示使用了多少个子域名（例如 shop.example.com）。过多可能意味着可疑或临时的搭建。"
            />
          </div>
          <span class="font-medium text-gray-800 dark:text-white"
            >{features.url.subdomain_count}</span
          >
        </div>
      {/if}

      {#if domainRandomness && domainRandomness.entropy !== undefined}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>域名随机度（熵）：</span>
            <TooltipIcon
              text="衡量域名的不可预测程度。熵值高通常说明是恶意软件（DGA）随机生成的域名。"
            />
          </div>
          <div class="flex items-center gap-2">
            <span class="font-medium text-gray-800 dark:text-white"
              >{domainRandomness.entropy.toFixed(2)}</span
            >
            {#if domainRandomness.entropy > 3.8}
              <span
                class="text-xs px-1.5 py-0.5 bg-red-100 dark:bg-red-900/30 text-red-600 dark:text-red-400 rounded-sm"
                >高熵</span
              >
            {/if}
          </div>
        </div>
      {/if}

      {#if typosquatResult}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>仿冒域名检查：</span>
            <TooltipIcon
              text="检查这个域名是否通过替换、增加或删除字符而酷似知名品牌（例如 paypa1.com 与 paypal.com），或者品牌名是否被嵌在更长的域名里。"
            />
          </div>
          {#if typosquatResult.is_suspicious}
            {#if typosquatResult.is_combo_squat}
              <span class="text-red-400 font-medium flex flex-col gap-0.5">
                <span class="inline-flex items-center gap-1"
                  ><StatusIcon kind="bad" /> 组合仿冒，包含品牌名</span
                >
                <span class="text-xs text-red-300 font-normal">
                  包含「<span class="font-mono">{typosquatResult.matched_brand}</span
                  >」（官方：{typosquatResult.matched_domain}）
                </span>
              </span>
            {:else}
              <span class="text-red-400 font-medium flex flex-col gap-0.5">
                <span class="inline-flex items-center gap-1"
                  ><StatusIcon kind="bad" /> 仿冒域名，与已知品牌相差 {typosquatResult.distance} 个字符</span
                >
                <span class="text-xs text-red-300 font-normal">
                  <span class="font-mono">{typosquatResult.matched_domain}</span>
                </span>
              </span>
            {/if}
          {:else}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> 未找到匹配</span
            >
          {/if}
        </div>
      {/if}
    </div>
  </section>
{/if}
