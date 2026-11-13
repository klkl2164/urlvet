<script lang="ts">
  import StatusIcon from "../StatusIcon.svelte";
  import TooltipIcon from "../TooltipIcon.svelte";
  import { ageLabel } from "../../results/sections";
  export let domainInfo: any;
  export let rank: number | undefined;

  // Some registries (.de, .eu, …) don't publish these dates; the API then
  // sends Go's zero time, 0001-01-01.
  const NOT_PUBLISHED = "注册局未公布";
  const hasDate = (d: string | undefined) => !!d && !d.startsWith("0001-");
</script>

{#if domainInfo}
  <section class="p-4 sm:p-5">
    <div
      class="space-y-0 divide-y divide-gray-300 dark:divide-gray-800 text-sm text-gray-800 dark:text-gray-200 max-w-4xl w-full mx-auto"
    >
      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>域名：</span>
          <TooltipIcon text="网站注册的名称，也就是用户在浏览器里输入并访问的内容。" />
        </div>
        <span class="font-medium text-gray-800 dark:text-white">{domainInfo.domain}</span>
      </div>

      {#if rank !== undefined}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>全球流量排名：</span>
            <TooltipIcon
              text="对网站全球受欢迎程度的粗略估计，数字越小表示访客越多。数据来自流量和参与度统计。"
            />
          </div>
          <span class="font-medium text-gray-800 dark:text-white"
            >{rank === 0 ? "未排名" : rank}</span
          >
        </div>
      {/if}

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>注册商：</span>
          <TooltipIcon
            text="负责管理这个域名注册的公司或机构（例如 GoDaddy、Namecheap、Google Domains）。"
          />
        </div>
        <span class="font-medium text-gray-800 dark:text-white">{domainInfo.registrar || "-"}</span>
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>域名年龄：</span>
          <TooltipIcon
            text="域名首次注册距今有多久。注册时间越早的域名，通常说明网站更成熟、更可信。"
          />
        </div>
        <span class="font-medium text-gray-800 dark:text-white"
          >{domainInfo.age_days == null ? NOT_PUBLISHED : ageLabel(domainInfo.age_days)}</span
        >
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>已启用 DNSSEC：</span>
          <TooltipIcon
            text="一项安全功能，通过对 DNS 数据做数字签名，帮助防止 DNS 篡改和重定向攻击。"
          />
        </div>
        {#if domainInfo.dnssec}
          <span class="text-emerald-700 dark:text-emerald-400 font-medium flex items-center gap-1.5"
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
          <span>创建时间：</span>
          <TooltipIcon text="这个域名首次注册并在互联网上生效的日期。" />
        </div>
        <span class="font-medium text-gray-800 dark:text-white"
          >{hasDate(domainInfo.created) ? domainInfo.created : NOT_PUBLISHED}</span
        >
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>更新时间：</span>
          <TooltipIcon
            text="域名注册信息最后一次被修改的日期（例如更改联系人或更新名称服务器）。"
          />
        </div>
        <span class="font-medium text-gray-800 dark:text-white"
          >{hasDate(domainInfo.updated) ? domainInfo.updated : NOT_PUBLISHED}</span
        >
      </div>

      <div
        class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] md:items-center gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
      >
        <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
          <span>到期时间：</span>
          <TooltipIcon text="这个域名的注册到期日期，除非所有者续费。" />
        </div>
        <span class="font-medium text-gray-800 dark:text-white"
          >{hasDate(domainInfo.expiry) ? domainInfo.expiry : NOT_PUBLISHED}</span
        >
      </div>

      {#if domainInfo.nameservers?.length}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>名称服务器：</span>
            <TooltipIcon
              text="负责把互联网流量指向正确网站主机的服务器，它们把域名和托管服务关联起来。"
            />
          </div>
          <div class="flex flex-wrap gap-2">
            {#each domainInfo.nameservers as ns}
              <span
                class="px-2 py-1 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-white rounded-sm text-xs"
                >{ns}</span
              >
            {/each}
          </div>
        </div>
      {/if}

      {#if domainInfo.status?.length}
        <div
          class="flex flex-col md:grid md:grid-cols-[minmax(0,280px),1fr] gap-2 md:gap-4 py-2 first:pt-0 last:pb-0"
        >
          <div class="flex items-center gap-1 text-gray-600 dark:text-gray-400">
            <span>状态：</span>
            <TooltipIcon
              text="域名的生命周期或控制状态，例如 'active'、'clientTransferProhibited' 或 'pendingDelete'。它们表示管理或运营上的限制。"
            />
          </div>
          <div class="flex flex-wrap gap-2">
            {#each domainInfo.status as st}
              <span
                class="px-2 py-1 bg-gray-100 dark:bg-gray-800 text-gray-700 dark:text-white rounded-sm text-xs"
                >{st}</span
              >
            {/each}
          </div>
        </div>
      {/if}
    </div>
  </section>
{/if}
