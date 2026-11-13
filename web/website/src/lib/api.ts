import { env } from '$env/dynamic/public';
import type {
  ApiResponse,
  DepthResponse,
  HstsResponse,
  IpCheckResponse,
  IpResolveResponse,
  LengthResponse,
  PunycodeResponse,
  RankResponse,
  RedirectsResponse,
  RiskyTldResponse,
  StatusCodeResponse,
  TrustedTldResponse,
  UrlShortenerResponse,
  WhoisResponse,
} from './types';

const PUBLIC_BASE_URL = env.PUBLIC_BASE_URL || 'http://localhost:8080/api/v1';

/** The backend writes errors in Chinese; the English cases keep older self-hosted backends working. */
function errorMessage(status: number, serverMessage?: string): string {
  if (serverMessage?.includes('链接太长') || serverMessage?.includes('maximum allowed length'))
    return '链接太长了，请缩短后再试。';
  switch (serverMessage) {
    case '请求太频繁':
    case 'too many requests':
      return '请求太频繁了，请稍后再试。';
    case '链接格式不正确':
    case '缺少链接参数':
    case 'invalid url':
    case 'invalid URL':
    case 'url query param is required':
    case 'url query parameter is required':
      return '这个地址看起来不对，检查一下再试。';
    case '无法分析这个链接':
    case 'could not analyze this URL':
      return '无法分析这个链接，请过一会儿再试。';
    case '报告内容不完整或格式不对':
    case 'invalid report':
      return '报告内容不完整或格式不对，检查一下再提交。';
    case '无法保存报告':
    case 'could not save report':
      return '报告没保存成功，请稍后再试。';
    case '截图失败':
    case 'failed to capture screenshot':
      return '截图没成功，请稍后再试。';
  }
  if (status === 429) return '请求太频繁了，请稍后再试。';
  if (status >= 500) return '服务器出错了，请稍后再试。';
  return `请求失败（HTTP ${status}），请稍后再试。`;
}

async function makeRequest<T>(endpoint: string, url: string): Promise<ApiResponse<T>> {
  try {
    // Ensure the URL is properly formatted before encoding
    const formattedUrl = url.trim();
    const encodedUrl = encodeURIComponent(formattedUrl);
    const fullUrl = `${PUBLIC_BASE_URL}${endpoint}?url=${encodedUrl}`;

    const response = await fetch(fullUrl);

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ error: '' }));
      console.error(`API Error:`, errorData);
      return { error: errorMessage(response.status, errorData.error) };
    }

    const data = await response.json();
    return { data };
  } catch (error) {
    console.error(`API Error for ${endpoint}:`, error);
    return { error: '网络错误，请检查网络连接后重试。' };
  }
}

async function makeScreenshotRequest(url: string): Promise<ApiResponse<string | null>> {
  try {
    const formattedUrl = url.trim();
    const encodedUrl = encodeURIComponent(formattedUrl);
    const fullUrl = `${PUBLIC_BASE_URL}/screenshot?url=${encodedUrl}`;

    const response = await fetch(fullUrl);

    if (!response.ok) {
      console.error(`Screenshot API Error: HTTP ${response.status}`);
      const errorData = await response.json().catch(() => ({ error: '' }));
      return { error: errorMessage(response.status, errorData.error) };
    }

    // Check if response is an image
    const contentType = response.headers.get('content-type');
    if (contentType && contentType.startsWith('image/')) {
      // Convert blob to object URL
      const blob = await response.blob();
      const blobUrl = URL.createObjectURL(blob);
      return { data: blobUrl };
    }

    // If not an image, return null
    return { data: null };
  } catch (error) {
    console.error(`Screenshot API Error:`, error);
    return { error: '网络错误，请检查网络连接后重试。' };
  }
}

export interface ReportPayload {
  url: string;
  verdict?: string;
  score?: number;
  expected_verdict?: 'Safe' | 'Suspicious' | 'Risky';
  comment?: string;
}

async function submitReport(payload: ReportPayload): Promise<ApiResponse<{ status: string }>> {
  try {
    const response = await fetch(`${PUBLIC_BASE_URL}/report`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(payload),
    });
    const data = await response.json().catch(() => ({}));
    if (!response.ok) {
      return { error: errorMessage(response.status, data.error) };
    }
    return { data };
  } catch {
    return { error: '网络错误，请检查网络连接后重试。' };
  }
}

export const api = {
  async report(payload: ReportPayload): Promise<ApiResponse<{ status: string }>> {
    return submitReport(payload);
  },
  async analyze(url: string): Promise<ApiResponse<any>> {
    return makeRequest<any>('/analyze', url);
  },
  async screenshot(url: string): Promise<ApiResponse<string | null>> {
    return makeScreenshotRequest(url);
  },
  async getRank(url: string): Promise<ApiResponse<RankResponse>> {
    return makeRequest<RankResponse>('/rank', url);
  },

  async getIpCheck(url: string): Promise<ApiResponse<IpCheckResponse>> {
    return makeRequest<IpCheckResponse>('/ip/check', url);
  },

  async getIpResolve(url: string): Promise<ApiResponse<IpResolveResponse>> {
    return makeRequest<IpResolveResponse>('/ip/resolve', url);
  },

  async getLength(url: string): Promise<ApiResponse<LengthResponse>> {
    return makeRequest<LengthResponse>('/length', url);
  },

  async getDepth(url: string): Promise<ApiResponse<DepthResponse>> {
    return makeRequest<DepthResponse>('/depth', url);
  },

  async getHsts(url: string): Promise<ApiResponse<HstsResponse>> {
    return makeRequest<HstsResponse>('/hsts', url);
  },

  async getRedirects(url: string): Promise<ApiResponse<RedirectsResponse>> {
    return makeRequest<RedirectsResponse>('/redirects', url);
  },

  async getPunycode(url: string): Promise<ApiResponse<PunycodeResponse>> {
    return makeRequest<PunycodeResponse>('/punycode', url);
  },

  async getTrustedTld(url: string): Promise<ApiResponse<TrustedTldResponse>> {
    return makeRequest<TrustedTldResponse>('/trusted-tld', url);
  },

  async getRiskyTld(url: string): Promise<ApiResponse<RiskyTldResponse>> {
    return makeRequest<RiskyTldResponse>('/risky-tld', url);
  },

  async getUrlShortener(url: string): Promise<ApiResponse<UrlShortenerResponse>> {
    return makeRequest<UrlShortenerResponse>('/url-shortener', url);
  },

  async getStatusCode(url: string): Promise<ApiResponse<StatusCodeResponse>> {
    return makeRequest<StatusCodeResponse>('/status-code', url);
  },

  async getWhois(url: string): Promise<ApiResponse<WhoisResponse>> {
    return makeRequest<WhoisResponse>('/whois', url);
  },
};
