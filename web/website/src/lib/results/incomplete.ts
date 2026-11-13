// Plain words for checks that didn't finish, keyed by backend task name.
const CHECK_NAMES: Record<string, string> = {
  phishtank_check: '钓鱼数据库（PhishTank）',
  safe_browsing_check: 'Google Safe Browsing',
  webrisk_check: 'Google Web Risk',
  content_check: '页面本身',
  http_combined_check: '重定向与 HTTPS',
  tls_combined_check: '安全证书',
  whois_lookup: '域名的注册信息',
  dns_validity_check: 'DNS 记录',
  ip_resolution: 'DNS 记录',
  domain_rank: '网站的知名度',
};

/** "、" between items and 和 before the last; a space keeps Latin text from gluing to it. */
function joinNames(names: string[]): string {
  if (names.length === 1) return names[0];
  const head = names.slice(0, -1).join('、');
  const glue = /[A-Za-z]$/.test(head) ? ' 和' : '和';
  return `${head}${glue}${names[names.length - 1]}`;
}

/** One sentence saying what wasn't checked, or null if everything ran. */
export function incompleteNote(
  checks: string[] | undefined,
  incomplete: boolean | undefined
): string | null {
  if (!checks?.length)
    return incomplete ? '有些检查没能完成，所以这个结果可能缺少一些信号。请过一分钟再试。' : null;

  const names = [...new Set(checks.map((c) => CHECK_NAMES[c] ?? '其他几项检查'))];
  const list = joinNames(names);
  return incomplete
    ? `这次没能检查「${list}」，所以这个结果可能缺少一些信号。请过一分钟再试。`
    : `这次没能检查「${list}」。其余检查都正常完成了。`;
}
