import { describe, expect, it } from 'vitest';
import { incompleteNote } from './incomplete';

describe('incompleteNote', () => {
  it('says nothing when every check ran', () => {
    expect(incompleteNote(undefined, false)).toBeNull();
    expect(incompleteNote([], false)).toBeNull();
  });

  it('names a PhishTank miss without doubting the result', () => {
    expect(incompleteNote(['phishtank_check'], false)).toBe(
      '这次没能检查「钓鱼数据库（PhishTank）」。其余检查都正常完成了。'
    );
  });

  it('lists what the site failed to answer, and suggests a retry', () => {
    expect(incompleteNote(['content_check', 'tls_combined_check'], true)).toBe(
      '这次没能检查「页面本身和安全证书」，所以这个结果可能缺少一些信号。请过一分钟再试。'
    );
  });

  it('merges checks that share a description', () => {
    expect(incompleteNote(['dns_validity_check', 'ip_resolution', 'domain_rank'], true)).toContain(
      'DNS 记录和网站的知名度'
    );
  });

  it('does not stack conjunctions when a name already contains one', () => {
    expect(incompleteNote(['content_check', 'http_combined_check'], true)).toContain(
      '页面本身和重定向与 HTTPS'
    );
  });

  it('keeps a space before 和 when the list ends in Latin text', () => {
    expect(incompleteNote(['safe_browsing_check', 'content_check'], false)).toContain(
      'Google Safe Browsing 和页面本身'
    );
  });

  it('names the Google checks', () => {
    expect(incompleteNote(['safe_browsing_check'], false)).toBe(
      '这次没能检查「Google Safe Browsing」。其余检查都正常完成了。'
    );
  });
});
