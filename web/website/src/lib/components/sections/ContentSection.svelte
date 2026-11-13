<script lang="ts">
  import { slide } from "svelte/transition";
  import StatusIcon from "../StatusIcon.svelte";
  import type { ContentData } from "../../types";
  export let contentData: ContentData | undefined;

  type Kind = "ok" | "warn" | "bad" | "info";
  type Check = { label: string; value: string; kind: Kind; hint: string };

  let showDetails = false;

  const VALUE_COLOR: Record<Kind, string> = {
    ok: "text-emerald-700 dark:text-emerald-400",
    warn: "text-yellow-700 dark:text-yellow-400",
    bad: "text-red-600 dark:text-red-400",
    info: "text-gray-700 dark:text-gray-300",
  };

  $: forms = contentData?.forms ?? [];
  $: iframes = contentData?.iframes ?? [];
  $: externalForms = forms.filter((f) => f.is_external).length;
  $: hiddenForm = forms.some((f) => f.is_hidden);
  $: brands = contentData?.brand_check?.detected_names ?? [];

  $: checks = contentData
    ? ([
        contentData.brand_check?.is_mismatch
          ? {
              label: "品牌",
              value: `冒充 ${contentData.brand_check.brand_found}`,
              kind: "bad",
              hint: "页面提到了知名品牌，但网站不在该品牌的官方域名上。",
            }
          : {
              label: "品牌",
              value: brands.length ? `已验证：${brands.join("、")}` : "未发现已知品牌",
              kind: brands.length ? "ok" : "info",
              hint: "检查页面上提到的品牌与网站所在的域名是否一致。",
            },
        {
          label: "登录表单",
          value: contentData.has_login_form ? "存在" : "无",
          kind: contentData.has_login_form ? "info" : "ok",
          hint: "带有密码或类似用户名输入框的表单。正规网站上有这种表单很正常。",
        },
        {
          label: "支付表单",
          value: contentData.has_payment_form ? "存在" : "无",
          kind: contentData.has_payment_form ? "warn" : "ok",
          hint: "要求填写卡号、CVV 或账单信息的表单。",
        },
        {
          label: "个人信息",
          value: contentData.has_personal_form ? "要求填写" : "无",
          kind: contentData.has_personal_form ? "info" : "ok",
          hint: "要求填写地址、电话号码等信息的表单。",
        },
        {
          label: "隐藏元素",
          value:
            contentData.has_hidden_iframe && hiddenForm
              ? "隐藏 iframe 和表单"
              : contentData.has_hidden_iframe
                ? "隐藏 iframe"
                : hiddenForm
                  ? "隐藏表单"
                  : "无",
          kind: contentData.has_hidden_iframe || hiddenForm ? "bad" : "ok",
          hint: "看不见的 iframe 或表单可以在后台运行，而你不会察觉。",
        },
        {
          label: "跟踪像素",
          value: contentData.has_tracking ? "存在" : "无",
          kind: contentData.has_tracking ? "info" : "ok",
          hint: "用来跟踪访问或邮件打开情况的 1×1 图片。",
        },
      ] as Check[])
    : [];

  function iconKind(k: Kind): "ok" | "warn" | "bad" {
    return k === "info" ? "ok" : k;
  }
</script>

{#if contentData}
  <section class="p-4 sm:p-5 space-y-4">
    <!-- Title -->
    <div>
      <p
        class="text-[10px] sm:text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-widest"
      >
        页面标题
      </p>
      <p class="mt-1 text-sm font-medium text-gray-900 dark:text-white break-words">
        {contentData.title || "（无标题）"}
      </p>
    </div>

    <!-- Checks -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-2">
      {#each checks as check}
        <div
          class="flex items-start gap-2.5 rounded-lg border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 px-3 py-2.5"
          title={check.hint}
        >
          <span
            class="mt-0.5 {check.kind === 'info'
              ? 'text-gray-400 dark:text-gray-500'
              : VALUE_COLOR[check.kind]}"
          >
            <StatusIcon kind={iconKind(check.kind)} />
          </span>
          <div class="min-w-0">
            <p class="text-xs text-gray-500 dark:text-gray-400">{check.label}</p>
            <p class="text-sm font-medium break-words {VALUE_COLOR[check.kind]}">{check.value}</p>
          </div>
        </div>
      {/each}
    </div>

    <!-- Forms summary -->
    <div class="flex flex-wrap items-center gap-x-3 gap-y-2 text-sm">
      {#if forms.length === 0}
        <span class="text-gray-600 dark:text-gray-400">此页面没有表单。</span>
      {:else if externalForms > 0}
        <span class="inline-flex items-center gap-1.5 text-red-600 dark:text-red-400 font-medium">
          <StatusIcon kind="bad" />
          {forms.length} 个表单中有 {externalForms} 个会把数据发送到其他域名
        </span>
      {:else}
        <span
          class="inline-flex items-center gap-1.5 text-emerald-700 dark:text-emerald-400 font-medium"
        >
          <StatusIcon kind="ok" />
          {forms.length} 个表单，全部提交到本站
        </span>
      {/if}

      {#if forms.length || iframes.length}
        <button
          type="button"
          class="ml-auto inline-flex items-center gap-1 font-mono text-xs text-gray-500 hover:text-gray-900 dark:hover:text-gray-100 focus:outline-hidden focus-visible:underline"
          aria-expanded={showDetails}
          on:click={() => (showDetails = !showDetails)}
        >
          {showDetails ? "隐藏" : "显示"}技术细节
          <svg
            class="w-3.5 h-3.5 transition-transform duration-200 {showDetails ? 'rotate-180' : ''}"
            viewBox="0 0 20 20"
            fill="currentColor"
            aria-hidden="true"
          >
            <path
              fill-rule="evenodd"
              d="M5.23 7.21a.75.75 0 011.06.02L10 11.188l3.71-3.958a.75.75 0 111.08 1.04l-4.25 4.53a.75.75 0 01-1.08 0l-4.25-4.53a.75.75 0 01.02-1.06z"
              clip-rule="evenodd"
            />
          </svg>
        </button>
      {/if}
    </div>

    {#if showDetails}
      <div transition:slide={{ duration: 200 }} class="space-y-3">
        {#each forms as form, i}
          <div
            class="rounded-lg border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 overflow-hidden"
          >
            <div
              class="flex flex-wrap items-center gap-2 px-3 py-2 border-b border-gray-200 dark:border-gray-800 bg-gray-50 dark:bg-gray-800/40"
            >
              <span
                class="text-[10px] font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-widest"
                >表单 {i + 1}</span
              >
              <span
                class="px-1.5 py-0.5 rounded-sm text-[10px] font-mono font-semibold uppercase bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300"
                >{form.method}</span
              >
              {#if form.is_hidden}
                <span
                  class="px-2 py-0.5 rounded-full border text-[10px] font-semibold bg-red-100 dark:bg-red-500/20 text-red-700 dark:text-red-300 border-red-300 dark:border-red-500/30"
                  >隐藏</span
                >
              {/if}
              {#if form.has_password}
                <span
                  class="px-2 py-0.5 rounded-full border text-[10px] font-semibold bg-yellow-100 dark:bg-yellow-500/20 text-yellow-700 dark:text-yellow-300 border-yellow-300 dark:border-yellow-500/30"
                  >密码</span
                >
              {/if}
              {#if form.has_payment}
                <span
                  class="px-2 py-0.5 rounded-full border text-[10px] font-semibold bg-red-100 dark:bg-red-500/20 text-red-700 dark:text-red-300 border-red-300 dark:border-red-500/30"
                  >支付</span
                >
              {/if}
              {#if form.has_user_like}
                <span
                  class="px-2 py-0.5 rounded-full border text-[10px] font-semibold bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300 border-gray-300 dark:border-gray-700"
                  >身份</span
                >
              {/if}
              {#if form.has_personal}
                <span
                  class="px-2 py-0.5 rounded-full border text-[10px] font-semibold bg-gray-100 dark:bg-gray-800 text-gray-600 dark:text-gray-300 border-gray-300 dark:border-gray-700"
                  >个人信息</span
                >
              {/if}
            </div>
            <dl class="px-3 py-2.5 space-y-2 text-sm">
              <div class="flex flex-col sm:flex-row sm:gap-3">
                <dt class="sm:w-28 shrink-0 text-xs text-gray-500 dark:text-gray-400 sm:pt-0.5">
                  提交到
                </dt>
                <dd class="min-w-0">
                  <span class="font-mono text-xs text-gray-800 dark:text-gray-200 break-all"
                    >{form.action || "（本页）"}</span
                  >
                  <span
                    class="ml-1 inline-flex items-center gap-1 text-xs font-medium align-middle {form.is_external
                      ? 'text-red-600 dark:text-red-400'
                      : 'text-emerald-700 dark:text-emerald-400'}"
                  >
                    <StatusIcon kind={form.is_external ? "bad" : "ok"} />
                    {form.is_external ? "外部域名" : "同站"}
                  </span>
                </dd>
              </div>
              {#if form.submit_texts?.length}
                <div class="flex flex-col sm:flex-row sm:gap-3">
                  <dt class="sm:w-28 shrink-0 text-xs text-gray-500 dark:text-gray-400 sm:pt-0.5">
                    按钮
                  </dt>
                  <dd class="flex flex-wrap gap-1 min-w-0">
                    {#each form.submit_texts as text}
                      <span
                        class="px-2 py-0.5 rounded-sm border border-gray-300 dark:border-gray-700 text-xs text-gray-700 dark:text-gray-300 break-all"
                        >{text}</span
                      >
                    {/each}
                  </dd>
                </div>
              {/if}
              {#if form.inputs?.length}
                <div class="flex flex-col sm:flex-row sm:gap-3">
                  <dt class="sm:w-28 shrink-0 text-xs text-gray-500 dark:text-gray-400 sm:pt-0.5">
                    字段
                  </dt>
                  <dd class="flex flex-col gap-1 min-w-0">
                    {#each form.inputs as input}
                      <span class="font-mono text-[11px] text-gray-600 dark:text-gray-400 break-all"
                        >{input}</span
                      >
                    {/each}
                  </dd>
                </div>
              {/if}
            </dl>
          </div>
        {/each}

        {#each iframes as iframe, i}
          <div
            class="rounded-lg border border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900 overflow-hidden"
          >
            <div
              class="flex items-center gap-2 px-3 py-2 border-b border-gray-200 dark:border-gray-800 bg-gray-50 dark:bg-gray-800/40"
            >
              <span
                class="text-[10px] font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-widest"
                >iframe {i + 1}</span
              >
              {#if iframe.is_hidden}
                <span
                  class="px-2 py-0.5 rounded-full border text-[10px] font-semibold bg-red-100 dark:bg-red-500/20 text-red-700 dark:text-red-300 border-red-300 dark:border-red-500/30"
                  >隐藏</span
                >
              {/if}
              <span class="ml-auto font-mono text-[11px] text-gray-500 dark:text-gray-400"
                >{iframe.width || "自动"} × {iframe.height || "自动"}</span
              >
            </div>
            <p class="px-3 py-2.5 font-mono text-xs text-gray-800 dark:text-gray-200 break-all">
              {iframe.src || "（无来源）"}
            </p>
          </div>
        {/each}
      </div>
    {/if}
  </section>
{/if}
