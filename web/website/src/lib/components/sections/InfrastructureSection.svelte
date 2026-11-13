<script lang="ts">
  import StatusIcon from "../StatusIcon.svelte";
  import TooltipIcon from "../TooltipIcon.svelte";
  export let infrastructure: any;
  export let isHostingPlatform: boolean = false;
</script>

{#if infrastructure}
  <section class="p-4 sm:p-5">
    <div
      class="space-y-0 divide-y divide-gray-300 dark:divide-gray-800 text-sm text-gray-800 dark:text-gray-200 max-w-4xl w-full mx-auto"
    >
      {#if infrastructure.ip_addresses?.length}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>服务器 IP 地址：</span>
            <TooltipIcon
              text="网站所在的实际网络地址。每个 IP 指向一台具体的物理服务器或云服务器。"
            />
          </div>
          <div class="flex flex-wrap gap-2">
            {#each infrastructure.ip_addresses as ip}
              <span
                class="px-2 py-1 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-white rounded-sm text-xs break-all"
                >{ip}</span
              >
            {/each}
          </div>
        </div>
      {/if}

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>DNS 配置：</span>
          <TooltipIcon
            text="名称服务器决定域名指向哪里，可以理解为互联网的「通讯录」，把域名和正确的托管服务商对应起来。"
          />
        </div>
        {#if infrastructure.nameservers_valid}
          <span class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
            ><StatusIcon kind="ok" /> 已检测到</span
          >
        {:else if isHostingPlatform}
          <span class="text-gray-500 dark:text-gray-400 font-medium flex items-center gap-1.5"
            >由托管平台管理</span
          >
        {:else}
          <span class="text-red-400 font-medium flex items-center gap-1.5"
            ><StatusIcon kind="bad" /> 未检测到</span
          >
        {/if}
      </div>

      {#if infrastructure.ns_hosts?.length > 0}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>名称服务器主机：</span>
            <TooltipIcon text="负责管理域名 DNS 设置的服务器，通常属于你的注册商或托管服务商。" />
          </div>
          <div class="flex flex-wrap gap-2">
            {#each infrastructure.ns_hosts as ns_host}
              <span
                class="px-2 py-1 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-white rounded-sm text-xs break-all"
                >{ns_host}</span
              >
            {/each}
          </div>
        </div>
      {/if}

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>邮件服务器记录：</span>
          <TooltipIcon text="MX 记录决定发往这个域名的邮件投递到哪里，是安全收发邮件的基础。" />
        </div>
        {#if infrastructure.mx_records_valid}
          <span class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
            ><StatusIcon kind="ok" /> 已检测到</span
          >
        {:else if isHostingPlatform}
          <span class="text-gray-500 dark:text-gray-400 font-medium flex items-center gap-1.5"
            >不适用</span
          >
        {:else}
          <span class="text-red-400 font-medium flex items-center gap-1.5"
            ><StatusIcon kind="bad" /> 未检测到</span
          >
        {/if}
      </div>

      {#if infrastructure.mx_hosts?.length > 0}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>MX 主机：</span>
            <TooltipIcon text="负责处理这个域名邮件往来的邮件服务器。" />
          </div>
          <div class="flex flex-wrap gap-2">
            {#each infrastructure.mx_hosts as mx_host}
              <span
                class="px-2 py-1 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-white rounded-sm text-xs break-all"
                >{mx_host}</span
              >
            {/each}
          </div>
        </div>
      {/if}
    </div>
  </section>
{/if}
