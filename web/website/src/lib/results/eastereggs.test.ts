import { describe, expect, it } from 'vitest';
import { REPO } from '../site';
import { easterEgg, localAddressMessage } from './eastereggs';

describe('easterEgg', () => {
  it.each([
    ['url.vet', '检查我自己'],
    ['https://www.url.vet/about', '检查我自己'],
    [REPO, '那是我住的地方'],
    ['https://www.youtube.com/watch?v=dQw4w9WgXcQ', 'Never gonna give you up'],
    ['https://youtu.be/dQw4w9WgXcQ', 'Never gonna give you up'],
    ['pаypal.com', '看出差别了吗'],
    ['imnotascam.lol', '首页的标语'],
    ['paypal.com', '这是真的那个'],
    ['https://www.apple.com/iphone', '这是真的那个'],
    ['example.com', '最无辜'],
    ['en.wikipedia.org/wiki/Phishing', '需要引用来源'],
    ['1.1.1.1', '电话簿'],
    ['https://invoice.zip', '不是文件'],
  ])('%s', (url, expected) => {
    expect(easterEgg(url)).toContain(expected);
  });

  it('stays quiet for ordinary links', () => {
    expect(easterEgg('github.com')).toBeNull();
    expect(easterEgg('https://www.youtube.com/watch?v=abc')).toBeNull();
    expect(easterEgg('paypal.com.evil.top')).toBeNull();
  });
});

describe('localAddressMessage', () => {
  it.each(['localhost', 'http://localhost:3000', '127.0.0.1', 'http://[::1]/'])(
    '%s 是你自己的电脑',
    (url) => {
      expect(localAddressMessage(url)).toContain('你自己的电脑');
    }
  );

  it.each(['192.168.1.1', '10.0.0.5', '172.20.1.1', 'printer.local'])('%s 是私有网络', (url) => {
    expect(localAddressMessage(url)).toContain('私有网络');
  });

  it('lets public addresses through', () => {
    expect(localAddressMessage('example.com')).toBeNull();
    expect(localAddressMessage('172.32.0.1')).toBeNull();
    expect(localAddressMessage('8.8.8.8')).toBeNull();
  });
});
