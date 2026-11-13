package analyzer

import (
	"fmt"
	"net"
	"strings"

	"github.com/klkl2164/urlvet/internal/service/checks"
	"github.com/klkl2164/urlvet/internal/service/threatfeeds"
)

// When a link leads to another site, the scan judges the page it lands on:
// trust comes only from there, since the site that sent the visitor on can't
// vouch for where it sent them. What the link itself gets wrong still counts,
// though, and so does anything bad it passes through on the way. These are
// the findings carried over from the link as given (the origin) and from the
// hops between it and the destination: risk only, never trust.

// OriginFindings is what the submitted link and the hops after it add to the
// destination's result.
type OriginFindings struct {
	URL        string   `json:"url"`
	BadReasons []string `json:"bad_reasons,omitempty"`
	Risk       int      `json:"risk"`
	// Confirmed is set when a list names the link or a hop: the result is
	// Risky, whatever the destination looks like.
	Confirmed bool `json:"confirmed,omitempty"`
}

func (o *OriginFindings) add(reason string, risk int) {
	o.BadReasons = append(o.BadReasons, reason)
	o.Risk += risk
}

// redirectTarget returns where the scanned page sends visitors on to another
// site, by HTTP redirect or straight away by script or meta refresh, or "".
// Moves within one site or brand were already cleared by applyExceptions.
func redirectTarget(resp Response) string {
	if rr := resp.Analysis.RedirectionResult; rr.HasDomainJump && rr.FinalURL != "" {
		return rr.FinalURL
	}
	if c := resp.ContentData; c != nil && c.ScriptRedirect != nil && c.ScriptRedirect.CrossDomain && c.ScriptRedirect.Target != "" {
		return c.ScriptRedirect.Target
	}
	return ""
}

// worthFollowing reports whether a redirect's destination should be judged
// in the link's place. A ranked site or one on a restricted registry sending
// visitors on, or anything leading to a site nobody knows, is judged where it
// lands. An unknown site landing on a well-known one is the shape of
// cloaking: kits send scanners to google.com or the brand's real login and
// victims to the phishing page, so the destination's good name can't be
// borrowed; the link is judged itself, and the jump counts against it.
func worthFollowing(src Response, dest string) bool {
	// A customer's subdomain on a hosting platform isn't the platform: its
	// rank, if any, is the platform's.
	ranked := src.Features.Rank > 0 && !src.Features.TLD.IsHostingPlatform
	return ranked || src.Features.TLD.IsTrusted || !isVouchedDestination(dest)
}

// originFindings collects the risks of the submitted link itself: how it's
// written, what lists say about it, and how its own domain looks, but not
// anything that would vouch for it.
func originFindings(src Response) *OriginFindings {
	o := &OriginFindings{URL: src.URL}
	host := hostOf(src.URL)
	link := fmt.Sprintf("你拿到的这条链接（%s）", host)

	// How it's written.
	if src.Features.URL.UsesIP {
		o.add(link+"是一个裸 IP 地址。", 100)
	}
	if service, ip, ok := checks.IPHostname(host); ok {
		what := "一个裸服务器"
		if ip != "" {
			what = fmt.Sprintf("一个裸服务器（%s）", ip)
		}
		o.add(fmt.Sprintf("%s是%s，被 %s 写成了域名。", link, what, service), 100)
	}
	if src.Features.URL.ContainsPunycode {
		o.add(link+"使用了 punycode，可以让一个地址看起来像另一个。", 100)
	}
	if src.Features.URL.HasHomoglyph {
		o.add(link+"使用了形近字母。", 60)
	}
	if src.Features.URL.Keywords.HasKeywords {
		o.add(fmt.Sprintf("%s中含有 %s 这样的词。", link, strings.Join(src.Features.URL.Keywords.Found, "、")), 10)
	}
	if victimAddressInLink(src.URL) {
		o.add(link+"里带有供页面自动填写的邮箱地址，这是钓鱼页面为每个目标定制登录页的常用方式。", 20)
	}

	// Brands it borrows.
	ts := src.TyposquatResult
	switch {
	case ts.IsSuspicious && ts.IsComboSquat:
		risk := 20
		if src.Features.TLD.IsHostingPlatform {
			risk = 40
		}
		o.add(fmt.Sprintf("%s使用了品牌名「%s」，但该网站并不属于这个品牌。", link, ts.MatchedBrand), risk)
	case ts.IsSuspicious:
		o.add(fmt.Sprintf("%s 与 %s 非常相似。", link, ts.MatchedDomain), 40)
	}
	switch {
	case ts.EmbeddedDomain != "":
		o.add(fmt.Sprintf("%s 的开头看起来像 %s，实际并不是。", link, ts.EmbeddedDomain), 60)
	case ts.SubdomainBrand != "" && ts.SubdomainBrandExact:
		o.add(fmt.Sprintf("%s 有一个以「%s」命名的子域名，但该网站并不属于这个品牌。", link, ts.SubdomainBrand), 20)
	case ts.SubdomainBrand != "":
		o.add(fmt.Sprintf("%s 有一个使用品牌名「%s」的子域名，但该网站并不属于这个品牌。", link, ts.SubdomainBrand), 40)
	}
	if ts.PathBrand != "" {
		if src.Features.TLD.IsTrusted {
			o.add(fmt.Sprintf("%s 是政府或大学域名下「%s」的登录路径，说明该网站可能已被入侵。", link, ts.PathBrand), 60)
		} else {
			o.add(fmt.Sprintf("%s 的路径中包含「%s」，但该网站并不属于这个品牌。", link, ts.PathBrand), 20)
		}
	}

	// What lists say about it.
	if fm := src.ThreatFeeds; fm != nil && fm.Listed {
		sources := strings.Join(fm.Sources, " 和 ")
		if fm.Match == "url" {
			o.add(fmt.Sprintf("已确认威胁：%s 被 %s 列为钓鱼或恶意软件。", link, sources), 200)
			o.Confirmed = true
		} else {
			o.add(fmt.Sprintf("%s 的其他页面被 %s 列为钓鱼或恶意软件。", host, sources), 80)
		}
	}
	for _, g := range []*threatfeeds.GoogleThreatResult{src.WebRisk, src.SafeBrowsing} {
		if g != nil && g.Listed {
			o.add(fmt.Sprintf("已确认威胁：Google 将 %s 列为%s（由 Google 提供提示）。", host, webRiskThreatNames(g.ThreatTypes)), 200)
			o.Confirmed = true
			break
		}
	}
	if p := src.Phishing; p != nil && p.InDatabase && p.Valid && p.Verified {
		o.add(fmt.Sprintf("已确认钓鱼：%s 是经过核实的钓鱼链接。", link), 200)
		o.Confirmed = true
	}

	// How its own domain looks.
	if age, known := src.DomainInfo.Age(); known && age <= newDomainDays &&
		!(src.Features.TLD.IsHostingPlatform && hostingRegistration(src)) {
		o.add(fmt.Sprintf("%s 所在的域名比较新（%s注册）。", link, ago(age)), 25)
	}
	if src.Features.TLD.IsRisky {
		o.add(link+"使用了高风险域名后缀。", 20)
	}
	if isWordPressSystemPage(src.URL) {
		o.add(link+"指向 WordPress 系统目录，钓鱼工具包常藏在被入侵网站的这类目录里。", 30)
	}
	if store := userContentHost(host); store != "" && src.ContentData != nil && src.ContentData.ScriptRedirect != nil {
		o.add(fmt.Sprintf("%s 是任何人都能上传到 %s 的文件，它只负责把访客转发走。", link, store), 60)
	}
	return o
}

// addHopFindings adds what the hops between the link and its destination
// pass through. Only local checks: the hops aren't fetched or looked up.
func addHopFindings(o *OriginFindings, chain []string, from, to string) {
	if len(chain) < 3 {
		return
	}
	fromDomain, _ := checks.GetDomain(from)
	toDomain, _ := checks.GetDomain(to)
	seen := map[string]bool{}
	for _, hop := range chain[1 : len(chain)-1] {
		host := hostOf(hop)
		domain, _ := checks.GetDomain(hop)
		if host == "" || seen[host] || sameSite(domain, fromDomain) || sameSite(domain, toDomain) {
			continue
		}
		seen[host] = true
		switch {
		case net.ParseIP(strings.Trim(host, "[]")) != nil:
			o.add(fmt.Sprintf("途中经过了裸 IP 地址（%s）。", host), 60)
		default:
			if service, ip, ok := checks.IPHostname(host); ok {
				o.add(fmt.Sprintf("途中经过了被 %[2]s 写成域名的裸服务器（%[1]s）。", ip, service), 60)
			}
		}
		if risky, _, _ := checks.IsRiskyTld(domain); risky {
			o.add(fmt.Sprintf("途中经过了 %s，该域名使用高风险后缀。", host), 10)
		}
		if fm := threatfeeds.LookupLocal(hop, false); fm.Listed {
			sources := strings.Join(fm.Sources, " 和 ")
			if fm.Match == "url" {
				o.add(fmt.Sprintf("已确认威胁：途中经过了 %s，它被 %s 列为钓鱼或恶意软件。", host, sources), 200)
				o.Confirmed = true
			} else {
				o.add(fmt.Sprintf("途中经过了 %s，该网站的其他页面被 %s 列为钓鱼或恶意软件。", host, sources), 80)
			}
		}
	}
}

// sameSite reports whether two registrable domains are the same site.
func sameSite(a, b string) bool {
	return a != "" && b != "" && (a == b || strings.HasSuffix(a, "."+b) || strings.HasSuffix(b, "."+a))
}
