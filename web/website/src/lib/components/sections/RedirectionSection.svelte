<script lang="ts">
  import StatusIcon from "../StatusIcon.svelte";
  import TooltipIcon from "../TooltipIcon.svelte";
  import RedirectPath from "./RedirectPath.svelte";
  export let analysis: any;
  export let domain: string;
  export let verdict: string | undefined = undefined;
</script>

{#if analysis}
  <section class="p-4 sm:p-5">
    <div
      class="space-y-0 divide-y divide-gray-300 dark:divide-gray-800 text-sm text-gray-800 dark:text-gray-200 max-w-4xl w-full mx-auto"
    >
      {#if analysis.redirection_result?.is_redirected && analysis.redirection_result.chain?.length > 1}
        <div class="pb-5">
          <RedirectPath chain={analysis.redirection_result.chain} {domain} {verdict} />
        </div>
      {/if}

      {#if analysis.redirection_result}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>已重定向：</span>
            <TooltipIcon text="表示访问该 URL 时是否会自动转发到另一个 URL。" />
          </div>
          {#if analysis.redirection_result.is_redirected}
            <span class="font-medium text-gray-900 dark:text-white break-all">是</span>
          {:else}
            <span class="font-medium text-gray-900 dark:text-white break-all">否</span>
          {/if}
        </div>

        {#if analysis.redirection_result.final_url}
          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>最终 URL 域名：</span>
              <TooltipIcon
                text="访客经过所有重定向后最终到达的域名。可以用来发现域名变化或钓鱼重定向。"
              />
            </div>
            <span class="font-medium text-gray-900 dark:text-white break-all"
              >{analysis.redirection_result.final_url_domain}</span
            >
          </div>
          <div
            class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
          >
            <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
              <span>最终 URL：</span>
              <TooltipIcon text="用户经过所有重定向后最终到达的完整 URL。" />
            </div>
            <span class="font-medium text-gray-900 dark:text-white break-all"
              >{analysis.redirection_result.final_url}</span
            >
          </div>
        {/if}

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>域名是否跳到其他域名：</span>
            <TooltipIcon text="检查网站是否重定向到完全不同的域名，这可能意味着钓鱼或跟踪。" />
          </div>
          {#if analysis.redirection_result.has_domain_jump}
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
            <span>重定向链长度：</span>
            <TooltipIcon text="显示网站到达最终目的地之前经过了多少次重定向。" />
          </div>
          <span class="font-medium text-gray-800 dark:text-white"
            >{analysis.redirection_result.chain_length}</span
          >
        </div>
      {/if}

      {#if analysis.http_status}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>HTTP 状态码：</span>
            <TooltipIcon
              text="访问该 URL 时服务器返回的响应码（例如 200 = OK，404 = Not Found）。"
            />
          </div>
          <span class="font-medium text-gray-800 dark:text-white"
            >{analysis.http_status.code} {analysis.http_status.text}</span
          >
        </div>

        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>重定向状态码（3xx）：</span>
            <TooltipIcon
              text="表示 HTTP 状态码是否为重定向（3xx）码，这类状态码会直接把访客送到另一个 URL。"
            />
          </div>
          {#if analysis.http_status.is_redirect}
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
      {/if}

      {#if analysis.is_hsts_supported !== undefined}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>支持 HSTS（仅 HTTPS）：</span>
            <TooltipIcon text="显示网站是否自动强制使用 HTTPS 连接，以提升安全性并防止攻击。" />
          </div>
          {#if analysis.is_hsts_supported}
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
      {/if}
    </div>
  </section>
{/if}
