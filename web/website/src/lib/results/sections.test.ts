import { describe, expect, it } from 'vitest';
import { ageLabel } from './sections';

describe('ageLabel', () => {
  it.each([
    [0, '今天'],
    [1, '1 天'],
    [29, '29 天'],
    [45, '1 个月'],
    [200, '6 个月'],
    [400, '1 年 1 个月'],
    [730, '2 年'],
  ])('%i days reads as %s', (days, expected) => {
    expect(ageLabel(days)).toBe(expected);
  });
});
