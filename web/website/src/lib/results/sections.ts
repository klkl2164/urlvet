import type { Tone } from '../ui/tone';
import type { AnalyzeResult } from '../types';

export type SectionId =
  | 'domain'
  | 'analysis'
  | 'threatintel'
  | 'security'
  | 'content'
  | 'features'
  | 'infrastructure';

export type Status = { tone: Tone; label: string };

/** Chinese relative age, from the API's day count. */
export function ageLabel(days: number): string {
  if (days <= 0) return '今天';
  if (days < 30) return `${days} 天`;
  if (days < 60) return '1 个月';
  if (days < 365) return `${Math.floor(days / 30)} 个月`;
  const years = Math.floor(days / 365);
  const months = Math.floor((days % 365) / 30);
  return months > 0 ? `${years} 年 ${months} 个月` : `${years} 年`;
}

// Sections open by default only when they hold a real problem.
export function computeExpanded(d: AnalyzeResult | null): Record<string, boolean> {
  if (!d)
    return {
      domain: false,
      analysis: false,
      threatintel: false,
      security: false,
      content: false,
      features: false,
      infrastructure: false,
      performance: false,
    };
  return {
    // age_days is null when the registry publishes no creation date, and
    // null < 365 is true in JS, so check for it explicitly.
    domain: d.domain_info?.age_days != null && d.domain_info.age_days < 365,
    analysis:
      !!d.analysis?.redirection_result?.has_domain_jump ||
      (d.analysis?.redirection_result?.chain_length ?? 0) > 3,
    threatintel:
      !!d.phishing?.valid ||
      !!d.threat_feeds?.listed ||
      !!d.web_risk?.listed ||
      !!d.safe_browsing?.listed,
    security: !!(
      d.ssl_info?.IsSuspicious ||
      d.ssl_info?.KnownBadChain ||
      (d.ssl_info && !d.ssl_info.HasTLS) ||
      (d.ssl_info?.HasTLS && !d.ssl_info.ChainValid) ||
      d.tls_info?.HostnameMismatch
    ),
    content: !!(
      d.content_data?.brand_check?.is_mismatch ||
      d.content_data?.has_hidden_iframe ||
      d.content_data?.forms?.some((f) => f.is_external)
    ),
    features: !!(
      d.features?.url?.has_homoglyph ||
      d.features?.url?.contains_punycode ||
      d.features?.url?.uses_ip ||
      d.features?.url?.url_shortener ||
      d.typosquat_result?.is_suspicious ||
      (d.domain_randomness?.entropy ?? 0) > 3.8 ||
      d.features?.tld?.is_risky_tld ||
      (d.features?.tld && !d.features?.tld?.is_icann && !d.features?.tld?.is_hosting_platform)
    ),
    infrastructure: !!(
      d.infrastructure &&
      !d.infrastructure.nameservers_valid &&
      !d.features?.tld?.is_hosting_platform
    ),
    performance: false,
  };
}

// One-glance summary shown on each collapsed section header.
export function sectionStatuses(d: AnalyzeResult): Record<SectionId, Status> {
  const days: number | undefined = d.domain_info?.age_days ?? undefined;
  const redir = d.analysis?.redirection_result;
  const ssl = d.ssl_info;
  const c = d.content_data;
  const url = d.features?.url;
  const tld = d.features?.tld;
  const ips: string[] = d.infrastructure?.ip_addresses ?? [];

  return {
    domain:
      days === undefined
        ? { tone: 'neutral', label: '年龄未知' }
        : days <= 30
          ? { tone: 'bad', label: `新建 · ${ageLabel(days)}` }
          : days < 365
            ? { tone: 'warn', label: `较新 · ${ageLabel(days)}` }
            : { tone: 'good', label: ageLabel(days) },
    analysis: redir?.has_domain_jump
      ? { tone: 'bad', label: '离开原域名' }
      : (redir?.chain_length ?? 0) > 3
        ? { tone: 'warn', label: `${redir.chain_length} 次跳转` }
        : redir?.is_redirected
          ? { tone: 'neutral', label: '站内重定向' }
          : { tone: 'good', label: '无重定向' },
    threatintel:
      d.web_risk?.listed || d.safe_browsing?.listed
        ? { tone: 'bad', label: 'Google 已标记' }
        : d.threat_feeds?.listed
          ? d.threat_feeds.match === 'url'
            ? { tone: 'bad', label: '已列为威胁' }
            : { tone: 'bad', label: '站内有页面被列入' }
          : d.phishing?.valid
            ? d.phishing.verified
              ? { tone: 'bad', label: '已确认钓鱼' }
              : { tone: 'bad', label: '被举报为钓鱼' }
            : { tone: 'good', label: '未列入' },
    security:
      ssl && !ssl.HasTLS
        ? { tone: 'bad', label: '无 HTTPS' }
        : ssl?.KnownBadChain || ssl?.IsSuspicious
          ? { tone: 'bad', label: '证书可疑' }
          : (ssl?.HasTLS && !ssl.ChainValid) || d.tls_info?.HostnameMismatch
            ? { tone: 'bad', label: '证书无效' }
            : { tone: 'good', label: '证书有效' },
    content: !c
      ? { tone: 'neutral', label: '未抓取' }
      : c.brand_check?.is_mismatch
        ? { tone: 'bad', label: '品牌不匹配' }
        : c.forms?.some((f) => f.is_external)
          ? { tone: 'bad', label: '外部表单' }
          : c.has_hidden_iframe
            ? { tone: 'bad', label: '隐藏 iframe' }
            : c.has_payment_form
              ? { tone: 'warn', label: '支付表单' }
              : c.has_login_form
                ? { tone: 'neutral', label: '登录表单' }
                : { tone: 'good', label: '无异常' },
    features: url?.has_homoglyph
      ? { tone: 'bad', label: '同形字符' }
      : url?.contains_punycode
        ? { tone: 'bad', label: 'Punycode' }
        : url?.uses_ip
          ? { tone: 'bad', label: '裸 IP' }
          : d.typosquat_result?.is_suspicious
            ? { tone: 'bad', label: '仿冒域名' }
            : url?.url_shortener
              ? { tone: 'bad', label: '短链接服务' }
              : (d.domain_randomness?.entropy ?? 0) > 3.8
                ? { tone: 'bad', label: '看起来随机' }
                : tld?.is_risky_tld
                  ? { tone: 'warn', label: '高风险 TLD' }
                  : tld && !tld.is_icann && !tld.is_hosting_platform
                    ? { tone: 'warn', label: '未受监管的 TLD' }
                    : { tone: 'good', label: '干净' },
    infrastructure:
      d.infrastructure && !d.infrastructure.nameservers_valid && !tld?.is_hosting_platform
        ? { tone: 'bad', label: 'DNS 问题' }
        : ips.length
          ? { tone: 'good', label: `${ips.length} 个 IP` }
          : { tone: 'neutral', label: 'DNS 正常' },
  };
}

// Outline icons (24x24) for each section header.
export const SECTION_ICONS: Record<SectionId, string> = {
  domain:
    'M9.568 3H5.25A2.25 2.25 0 003 5.25v4.318c0 .597.237 1.17.659 1.591l9.581 9.581c.699.699 1.78.872 2.607.33a18.095 18.095 0 005.223-5.223c.542-.827.369-1.908-.33-2.607L11.16 3.66A2.25 2.25 0 009.568 3zM6 6h.008v.008H6V6z',
  analysis: 'M7.5 21L3 16.5m0 0L7.5 12M3 16.5h13.5m0-13.5L21 7.5m0 0L16.5 12M21 7.5H7.5',
  threatintel:
    'M9 12.75L11.25 15 15 9.75m-3-7.036A11.959 11.959 0 013.598 6 11.99 11.99 0 003 9.749c0 5.592 3.824 10.29 9 11.623 5.176-1.332 9-6.03 9-11.622 0-1.31-.21-2.571-.598-3.751h-.152c-3.196 0-6.1-1.248-8.25-3.285z',
  security:
    'M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z',
  content:
    'M19.5 14.25v-2.625a3.375 3.375 0 00-3.375-3.375h-1.5A1.125 1.125 0 0113.5 7.125v-1.5a3.375 3.375 0 00-3.375-3.375H8.25m0 12.75h7.5m-7.5 3H12M10.5 2.25H5.625c-.621 0-1.125.504-1.125 1.125v17.25c0 .621.504 1.125 1.125 1.125h12.75c.621 0 1.125-.504 1.125-1.125V11.25a9 9 0 00-9-9z',
  features:
    'M13.19 8.688a4.5 4.5 0 011.242 7.244l-4.5 4.5a4.5 4.5 0 01-6.364-6.364l1.757-1.757m13.35-.622l1.757-1.757a4.5 4.5 0 00-6.364-6.364l-4.5 4.5a4.5 4.5 0 001.242 7.244',
  infrastructure:
    'M5.25 14.25h13.5m-13.5 0a3 3 0 01-3-3m3 3a3 3 0 100 6h13.5a3 3 0 100-6m-16.5-3a3 3 0 013-3h13.5a3 3 0 013 3m-19.5 0a4.5 4.5 0 01.9-2.7L5.737 5.1a3.375 3.375 0 012.7-1.35h7.126c1.062 0 2.062.5 2.7 1.35l2.587 3.45a4.5 4.5 0 01.9 2.7m0 0a3 3 0 01-3 3m0 3h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008zm-3 6h.008v.008h-.008v-.008zm0-6h.008v.008h-.008v-.008z',
};
