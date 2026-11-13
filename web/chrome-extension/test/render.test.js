// @vitest-environment jsdom
import { beforeAll, describe, expect, it } from "vitest";
import * as R from "../src/render.js";

beforeAll(() => {
  globalThis.requestAnimationFrame ??= (cb) =>
    setTimeout(() => cb(performance.now()), 0);
});

const result = (over = {}) => ({
  url: "https://evil.example/login",
  verdict: "Risky",
  score: 12,
  label: "高风险",
  quip: "很可能不安全，最好别点。 🚫",
  bad: ["Very low traffic volume."],
  good: [],
  shortLink: null,
  incomplete: false,
  incompleteChecks: [],
  ...over,
});

describe("verdictCard", () => {
  it("shows the verdict, its badge and quip", () => {
    const card = R.verdictCard(result());
    expect(card.dataset.verdict).toBe("Risky");
    expect(card.textContent).toContain("危险");
    expect(card.textContent).toContain("高风险");
    expect(card.querySelector(".ring").getAttribute("aria-label")).toBe(
      "信任分 12，满分 100",
    );
  });
});

describe("flags", () => {
  it("never parses text from the network as HTML", () => {
    const evil = '<img src=x onerror="alert(1)">';
    const [section] = R.flags(result({ bad: [evil] }));
    expect(section.querySelector("img")).toBeNull();
    expect(section.textContent).toContain(evil);
  });
  it("shows green flags for a safe result, and how many more there are", () => {
    const [section] = R.flags(
      result({ verdict: "Safe", bad: [], good: ["a", "b", "c", "d", "e"] }),
      3,
    );
    expect(section.querySelector("h3").textContent).toBe("正面标记");
    expect(section.querySelectorAll("li")).toHaveLength(3);
    expect(section.textContent).toContain("完整报告中还有 2 条");
  });
  it("shows nothing when there are no flags", () => {
    expect(R.flags(result({ bad: [], good: [] }))).toEqual([]);
  });
});

describe("notes", () => {
  it("uses the website's short-link wording", () => {
    const [note] = R.notes(
      result({
        shortLink: {
          url: "https://goo.su/TIcsAp",
          chain: ["https://goo.su/TIcsAp", "https://www.roblox.com.bn/x"],
          target: "https://www.roblox.com.bn/x",
          resolved: true,
        },
      }),
    );
    expect(note.textContent).toBe(
      "这是 goo.su 上的一个短链接。它指向 www.roblox.com.bn，这个结果是针对那一页的。",
    );
  });
  it("says which checks didn't finish", () => {
    const [note] = R.notes(
      result({ incomplete: true, incompleteChecks: ["content_check"] }),
    );
    expect(note.textContent).toContain("页面本身");
  });
});

describe("wordmark", () => {
  it("takes the verdict's state on its dot", () => {
    expect(R.wordmark("Risky").querySelector(".dot").dataset.state).toBe(
      "Risky",
    );
  });
});

describe("redirect note", () => {
  it("says where a well-known site's redirect led", () => {
    const [note] = R.notes(
      result({
        redirectedFrom: {
          url: "https://sba.yandex.ru/redirect?url=x",
          chain: [],
          target: "https://nus.nbt.mybluehost.me/x",
        },
      }),
    );
    expect(note.textContent).toContain("会把访客转发到 nus.nbt.mybluehost.me");
  });
});
