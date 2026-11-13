import { describe, expect, it } from 'vitest';
import { redirectNote, shortLinkNote } from './shortlink';

describe('shortLinkNote', () => {
  it('says nothing for an ordinary link', () => {
    expect(shortLinkNote(undefined)).toBeNull();
  });

  it('names where a short link leads', () => {
    expect(
      shortLinkNote({
        url: 'https://goo.su/TIcsAp',
        chain: ['https://goo.su/TIcsAp', 'https://www.roblox.com.bn/users/1/profile'],
        target: 'https://www.roblox.com.bn/users/1/profile',
        resolved: true,
      })
    ).toBe('这是 goo.su 上的一个短链接。它指向 www.roblox.com.bn，这个结果是针对那一页的。');
  });

  it('names short links passed through on the way', () => {
    expect(
      shortLinkNote({
        url: 'https://bit.ly/x',
        chain: ['https://bit.ly/x', 'https://tinyurl.com/y', 'https://evil.example/login'],
        target: 'https://evil.example/login',
        resolved: true,
      })
    ).toBe(
      '这是 bit.ly 上的一个短链接。它指向 evil.example（途经 tinyurl.com），这个结果是针对那一页的。'
    );
  });

  it('says when the destination could not be found', () => {
    expect(
      shortLinkNote({ url: 'https://bom.so/abc', chain: ['https://bom.so/abc'], resolved: false })
    ).toBe(
      '这是 bom.so 上的一个短链接。我们看不到它指向哪里，所以这个结果只针对短链接本身，最高也只能判定为可疑。'
    );
  });
});

describe('redirectNote', () => {
  it('says nothing without a redirect', () => {
    expect(redirectNote(undefined)).toBeNull();
  });
  it('names the site that redirected and where to', () => {
    expect(
      redirectNote({
        url: 'https://sba.yandex.ru/redirect?url=x',
        chain: ['https://sba.yandex.ru/redirect?url=x', 'https://nus.nbt.mybluehost.me/x'],
        target: 'https://nus.nbt.mybluehost.me/x',
      })
    ).toBe(
      'sba.yandex.ru 上的这个链接会把访客转发到 nus.nbt.mybluehost.me。这个结果是针对那一页的，也包括链接本身存在的问题。'
    );
  });
});
