<script lang="ts">
  import StatusIcon from "../StatusIcon.svelte";
  import type { SSLInfo, TLSInfo } from "../../types";
  import TooltipIcon from "../TooltipIcon.svelte";
  export let sslInfo: SSLInfo | undefined;
  export let tlsInfo: TLSInfo | undefined;
</script>

{#if sslInfo || tlsInfo}
  <section class="p-4 sm:p-5">
    <div
      class="space-y-0 divide-y divide-gray-300 dark:divide-gray-800 text-sm text-gray-800 dark:text-gray-200 max-w-4xl w-full mx-auto"
    >
      {#if sslInfo}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>SSL 支持：</span>
            <TooltipIcon text="检查网站是否支持安全的 HTTPS 连接。" />
          </div>
          {#if sslInfo.HasTLS}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> 已启用</span
            >
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> 未启用</span
            >
          {/if}
        </div>

        {#if sslInfo.HasTLS}
          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>证书链：</span>
              <TooltipIcon text="验证 SSL 证书是否由受信任的机构签发，以及完整证书链是否有效。" />
            </div>
            {#if sslInfo.ChainValid}
              <span
                class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
                ><StatusIcon kind="ok" /> 有效</span
              >
            {:else}
              <span class="text-red-400 font-medium flex items-center gap-1.5"
                ><StatusIcon kind="bad" /> 无效 / 自签名</span
              >
            {/if}
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>证书颁发机构：</span>
              <TooltipIcon text="签发这张 SSL 证书的机构。" />
            </div>
            <span class="font-medium text-gray-800 dark:text-white">{sslInfo.Issuer || "-"}</span>
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>证书年龄：</span>
              <TooltipIcon text="证书是几天前签发的。新域名上刚签发的证书可能值得怀疑。" />
            </div>
            <span class="font-medium text-gray-800 dark:text-white">{sslInfo.AgeDays} 天</span>
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>生效时间：</span>
              <TooltipIcon text="这张证书开始生效的日期。" />
            </div>
            <span class="font-medium text-gray-800 dark:text-white">{sslInfo.NotBefore}</span>
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>到期时间：</span>
              <TooltipIcon text="当前 SSL 证书的到期时间。" />
            </div>
            <span class="font-medium text-gray-800 dark:text-white">{sslInfo.NotAfter}</span>
          </div>

          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>证书风险等级：</span>
              <TooltipIcon text="对证书技术完整性的总体评估。" />
            </div>
            {#if !sslInfo.IsSuspicious}
              <span
                class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
                ><StatusIcon kind="ok" /> 低风险</span
              >
            {:else}
              <span class="text-amber-700 dark:text-amber-400 font-medium flex items-center gap-1.5"
                ><StatusIcon kind="warn" /> 可疑</span
              >
            {/if}
          </div>

          {#if sslInfo.Reasons && sslInfo.Reasons.length > 0}
            <div
              class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2"
            >
              <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
                <span>技术警告：</span>
                <TooltipIcon text="这张证书被标记的具体技术原因。" />
              </div>
              <ul class="text-xs text-amber-700 dark:text-amber-400 list-disc list-inside">
                {#each sslInfo.Reasons as reason}
                  <li>{reason}</li>
                {/each}
              </ul>
            </div>
          {/if}

          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>证书指纹：</span>
              <TooltipIcon text="这张证书的唯一标识（SHA-256 哈希）。" />
            </div>
            <span class="font-mono text-[12px] text-gray-600 dark:text-gray-300 break-all"
              >{sslInfo.Fingerprint}</span
            >
          </div>
        {/if}
      {/if}

      {#if tlsInfo}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>TLS 颁发机构（连接）：</span>
            <TooltipIcon text="在实时连接中检测到的证书颁发机构。" />
          </div>
          <span class="font-medium text-gray-800 dark:text-white">{tlsInfo.Issuer || "-"}</span>
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>主机名匹配：</span>
            <TooltipIcon text="确保证书确实是为正在访问的域名签发的。" />
          </div>
          {#if !tlsInfo.HostnameMismatch}
            <span
              class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="ok" /> 匹配</span
            >
          {:else}
            <span class="text-red-400 font-medium flex items-center gap-1.5"
              ><StatusIcon kind="bad" /> 不匹配</span
            >
          {/if}
        </div>
      {/if}
    </div>
  </section>
{/if}
