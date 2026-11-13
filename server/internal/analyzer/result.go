package analyzer

import (
	"encoding/base64"
	"fmt"
	"math"
	neturl "net/url"
	"path"
	"regexp"
	"slices"
	"strings"

	"github.com/klkl2164/urlvet/internal/constants"
	"github.com/klkl2164/urlvet/internal/service/checks"
	"github.com/klkl2164/urlvet/internal/service/threatfeeds"
)

func GenerateResult(resp Response) Result {
	var neutralReasons []string
	var goodReasons []string
	var badReasons []string
	trustScore := 0
	riskScore := 0

	// --- 1. Popularity & Rank ---
	if resp.Features.Rank == 0 {
		if service, ok := constants.ProviderServiceHosts[hostOf(resp.URL)]; ok {
			neutralReasons = append(neutralReasons, fmt.Sprintf("%s 上托管的文件（任何人都能上传）。", service))
		} else if service, ok := constants.UserPageHosts[strings.TrimPrefix(hostOf(resp.URL), "www.")]; ok {
			neutralReasons = append(neutralReasons, fmt.Sprintf("%s 上的一个页面（任何人都能在这里创建）。", service))
		} else if store := userContentHost(hostOf(resp.URL)); store != "" && resp.Features.TLD.IsHostingPlatform {
			neutralReasons = append(neutralReasons, fmt.Sprintf("存放在 %s（任何人都能上传文件）。", store))
		} else if resp.Features.TLD.IsHostingPlatform {
			// Personal/project pages on known hosting platforms are inherently unranked — not a risk signal.
			neutralReasons = append(neutralReasons, fmt.Sprintf("%s 下没有排名的子域名（个人或项目页面很常见）。", resp.Features.TLD.TLD))
		} else if age, ok := resp.DomainInfo.Age(); ok && age > establishedAfterDays {
			// Most real sites never make the top million. Being unranked only
			// points at phishing alongside a domain registered weeks ago.
			neutralReasons = append(neutralReasons, "不在全球前 100 万网站之列（小网站很常见）。")
		} else {
			badReasons = append(badReasons, "流量非常低。")
			riskScore += 10
		}
	} else if resp.Features.Rank > 0 && resp.Features.Rank <= 10000 {
		goodReasons = append(goodReasons, fmt.Sprintf("全球头部网站：全球排名第 %d。", resp.Features.Rank))
		trustScore += 90
	} else if resp.Features.Rank <= 50000 {
		// 10,001–50,000: well-known sites
		goodReasons = append(goodReasons, fmt.Sprintf("知名网站（全球排名第 %d）。", resp.Features.Rank))
		trustScore += 50
	} else {
		// 50,001+: ranked but lower traffic
		goodReasons = append(goodReasons, fmt.Sprintf("流量不高但已被收录的网站（全球排名第 %d）。", resp.Features.Rank))
		trustScore += 20
	}

	// --- 2. TLD (Top Level Domain) ---
	if resp.Features.TLD.IsRisky {
		badReasons = append(badReasons, "检测到高风险域名后缀（常与垃圾信息关联）。")
		riskScore += 20
	}

	// Exceptions keep a brand in the path on .gov or .edu only for a login
	// page (/secure/netflix/login): a hacked site, which the registry can't
	// vouch for.
	if resp.Features.TLD.IsTrusted && resp.TyposquatResult.PathBrand != "" {
		badReasons = append(badReasons, fmt.Sprintf("政府或大学域名下出现「%s」的登录页面，说明该网站可能已被入侵。", resp.TyposquatResult.PathBrand))
		riskScore += 60
	} else if resp.Features.TLD.IsTrusted {
		goodReasons = append(goodReasons, "高可信的官方域名后缀（政府/教育）。")
		trustScore += 100
	} else if resp.Features.TLD.IsICANN && !resp.Features.TLD.IsRisky {
		// Keep it simple
		neutralReasons = append(neutralReasons, "标准的、官方认可的域名后缀。")
	}

	if !resp.Features.TLD.IsICANN {
		if resp.Features.TLD.IsHostingPlatform {
			// PSL private entry but operated by a reputable hosting company — small trust signal.
		} else {
			badReasons = append(badReasons, "不受监管或非标准的域名后缀。")
			riskScore += 30
		}
	}

	// --- 3. Security Protocols ---
	// HSTS is one switch at a CDN; on a domain registered this month it says
	// nothing about who runs it.
	if age, known := resp.DomainInfo.Age(); resp.Analysis.SupportsHSTS && known && age <= newDomainDays {
		neutralReasons = append(neutralReasons, "发送了 HSTS 响应头，但对这么新的域名来说说明不了什么。")
	} else if resp.Analysis.SupportsHSTS {
		goodReasons = append(goodReasons, "启用了严格的 HTTPS 安全策略（HSTS）。")
		trustScore += 20
	}

	// --- 4. URL Structure / Obfuscation ---
	if resp.Features.URL.IsURLShortener {
		badReasons = append(badReasons, "检测到短链接服务（会隐藏真实目的地）。")
		riskScore += 25
	}

	// Uses IP
	if resp.Features.URL.UsesIP {
		badReasons = append(badReasons, "直接使用 IP 地址（常见的规避手法）。")
		riskScore += 100
	}
	// The same, dressed up as a name: 1-2-3-4.sslip.io is 1.2.3.4.
	if service, ip, ok := checks.IPHostname(hostOf(resp.URL)); ok {
		if ip != "" {
			badReasons = append(badReasons, fmt.Sprintf("这个地址是通过 %[2]s 写成域名的裸服务器，也就是 %[1]s；这类服务能把任意 IP 地址变成域名。", ip, service))
		} else {
			badReasons = append(badReasons, fmt.Sprintf("这个地址位于 %s；这类服务能把任意 IP 地址变成域名。", service))
		}
		riskScore += 100
	}

	// Punycode
	if resp.Features.URL.ContainsPunycode {
		badReasons = append(badReasons, "检测到 Punycode 字符（可能用于钓鱼仿冒）。")
		riskScore += 100
	}

	// Too deep
	if resp.Features.URL.TooDeep {
		badReasons = append(badReasons, "URL 路径层级过多（可能用于隐藏真实请求）。")
		riskScore += 30
	}

	// Too long
	if resp.Features.URL.TooLong {
		badReasons = append(badReasons, "URL 长度超过常规限制（可能用于缓冲区溢出或隐藏内容）。")
		riskScore += 20
	}

	// Subdomain Count
	if resp.Features.URL.SubdomainCount > 2 {
		badReasons = append(badReasons, "子域名数量异常。")
		riskScore += 15
	}

	// Keywords
	if resp.Features.URL.Keywords.HasKeywords {
		badReasons = append(badReasons, fmt.Sprintf("URL 中含有敏感的安全关键词：%s", strings.Join(resp.Features.URL.Keywords.Found, "、")))
		riskScore += 10
	}

	// --- 5. Infrastructure Forensics ---
	if !resp.Infrastructure.NameserversValid {
		if resp.Features.TLD.IsHostingPlatform {
			// Hosting platforms manage the entire DNS zone; individual subdomains have no own NS records.
		} else {
			badReasons = append(badReasons, "DNS 配置不完整或缺失。")
			riskScore += 10
		}
	}

	// MX records
	if !resp.Infrastructure.MXRecordsValid && !resp.Features.TLD.IsHostingPlatform {
		neutralReasons = append(neutralReasons, "这个域名没有配置邮件服务器。")
		riskScore += 5
	}

	// A page on a brand's own hosting that forwards to the brand's own site
	// (googleresearch.blogspot.com → research.google) ends somewhere the
	// brand runs, which vouches for it.
	if final := resp.Analysis.RedirectionResult.FinalURLHost; resp.Features.TLD.IsHostingPlatform &&
		final != "" && final != hostOf(resp.URL) && sameBrand(resp.Features.TLD.TLD, final) && isVouchedDestination(resp.Analysis.RedirectionResult.FinalURL) {
		goodReasons = append(goodReasons, fmt.Sprintf("最终跳转到 %s，与 %s 属于同一家公司。", final, resp.Features.TLD.TLD))
		trustScore += 50
	}

	// --- 6. Domain History ---
	// On a hosting platform the registration is the provider's, so its age
	// says nothing about the site.
	if resp.DomainInfo != nil && resp.Features.TLD.IsHostingPlatform && hostingRegistration(resp) {
		if _, known := resp.DomainInfo.Age(); known {
			neutralReasons = append(neutralReasons, fmt.Sprintf("这个域名的注册时间属于 %s，不属于本站。", resp.Features.TLD.TLD))
		}
	} else if resp.DomainInfo != nil {
		// Years of clean operation is real evidence, and should count for more
		// than a single header like HSTS. A domain whose registry publishes no
		// creation date gets neither the penalty nor the bonus.
		age, known := resp.DomainInfo.Age()
		switch {
		case !known:
			neutralReasons = append(neutralReasons, "这个域名的注册局没有公布注册日期。")
		case age <= newDomainDays:
			badReasons = append(badReasons, fmt.Sprintf("%s注册。钓鱼网站通常都是新注册的。", ago(age)))
			riskScore += 25
		case age <= establishedAfterDays:
			badReasons = append(badReasons, fmt.Sprintf("域名较新，%s注册。请谨慎。", ago(age)))
			riskScore += 15
		case age <= oneYearDays:
			// Phishing domains are used within weeks of registration and burn
			// out fast; one that has lasted a few months isn't that profile.
			neutralReasons = append(neutralReasons, fmt.Sprintf("域名比较新（已注册 %s）。", agePhrase(age)))
		case (resp.Features.Rank == 0 || resp.Features.Rank > wellKnownRank) &&
			(kitLoginPath(resp.URL) || (victimAddressInLink(resp.URL) && resp.ContentData != nil && resp.ContentData.HasLoginForm)):
			// Hacked sites keep their age: a phishing kit dropped on an old
			// domain borrows years of history it had no part in.
			neutralReasons = append(neutralReasons, fmt.Sprintf("域名已注册 %s，但这个链接看起来是被人植入的登录页面，所以不计域名年龄。", agePhrase(age)))
		case (resp.Features.Rank == 0 || resp.Features.Rank > wellKnownRank) && adobeBaitStrong(resp):
			// Fake Adobe viewers are dropped on hacked sites the same way
			// (luredout.com, nearly 10 years old).
			neutralReasons = append(neutralReasons, fmt.Sprintf("域名已注册 %s，但这个页面看起来是被人植入的，所以不计域名年龄。", agePhrase(age)))
		case age <= fiveYearsDays:
			goodReasons = append(goodReasons, fmt.Sprintf("已稳定运行 %s。", agePhrase(age)))
			if age > threeYearsDays {
				trustScore += 15
			} else {
				trustScore += 10
			}
		default:
			goodReasons = append(goodReasons, fmt.Sprintf("域名历史悠久（已注册 %s）。", agePhrase(age)))
			trustScore += 25
		}

		// DNSSEC Logic Updated
		if resp.DomainInfo.DNSSEC {
			goodReasons = append(goodReasons, "启用了高级 DNS 安全（DNSSEC）。")
			trustScore += 10
		} else {
			// Moved to Neutral. Not having DNSSEC is NOT a sign of phishing for .coms
			neutralReasons = append(neutralReasons, "标准 DNS 安全（未启用 DNSSEC）。")
			// Removed riskScore penalty
		}
	}

	// --- 7. Redirection Analysis ---
	if resp.Analysis.RedirectionResult.IsRedirected {
		// ChainLength counts URLs, including the one scanned. Long chains only
		// matter when they leave the site: hopping www → locale → home on one
		// domain is normal, while phishing bounces through other domains.
		hops := resp.Analysis.RedirectionResult.ChainLength - 1
		if hops > 3 && resp.Analysis.RedirectionResult.HasDomainJump {
			badReasons = append(badReasons, fmt.Sprintf("检测到过长的重定向链（%d 次跳转）。", hops))
			riskScore += 40
		}

		if resp.Analysis.RedirectionResult.HasDomainJump {
			badReasons = append(badReasons, "检测到跨域名重定向（目的地与来源不同）。")
			// Add the destination as a neutral fact so they can see where they are going
			badReasons = append(badReasons, fmt.Sprintf("最终目的地：%s。更多信息请查看报告。", resp.Analysis.RedirectionResult.FinalURLHost))
			riskScore += 50
		}
	}

	// --- 8. Homoglyphs ---
	if resp.Features.URL.HasHomoglyph {
		badReasons = append(badReasons, "检测到同形异义字符攻击（用视觉上相似的字符伪装）。")
		riskScore += 60
	}

	// --- 9. Typosquatting / Combo-squatting ---
	if resp.TyposquatResult.IsSuspicious {
		if resp.TyposquatResult.IsComboSquat {
			badReasons = append(badReasons, fmt.Sprintf("检测到组合抢注：域名包含品牌名「%s」，但不是官方网站。", resp.TyposquatResult.MatchedBrand))
			// On a hosting platform anyone can claim "instagram-login" for free
			// in seconds; a registered domain at least costs something.
			if resp.Features.TLD.IsHostingPlatform {
				riskScore += 40
			} else {
				riskScore += 20
			}
			if bare := brandPageIsBare(resp); bare != "" {
				badReasons = append(badReasons, fmt.Sprintf("在不知名网站上使用「%s」的名称，而且页面%s。", resp.TyposquatResult.MatchedBrand, bare))
				riskScore += 20
			}
		} else {
			badReasons = append(badReasons, fmt.Sprintf("检测到域名仿冒：域名与「%s」非常相似（相差 %d 个字符）。", resp.TyposquatResult.MatchedDomain, resp.TyposquatResult.Distance))
			riskScore += 40
		}
	}

	// Brands in the rest of the URL. Only reached on sites that aren't well
	// known: applyExceptions clears these for those.
	ts := resp.TyposquatResult
	if ts.EmbeddedDomain != "" {
		badReasons = append(badReasons, fmt.Sprintf("地址开头看起来像 %s，实际属于 %s。", ts.EmbeddedDomain, resp.Domain))
		riskScore += 60
	}
	if ts.SubdomainBrand != "" && ts.EmbeddedDomain == "" {
		if bare := brandPageIsBare(resp); bare != "" {
			badReasons = append(badReasons, fmt.Sprintf("在不知名网站上，有一个子域名以「%s」命名，而且页面%s。", ts.SubdomainBrand, bare))
			riskScore += 20
		}
		if ts.SubdomainBrandExact {
			// github.acme.com is a common, honest pattern; it counts for less.
			badReasons = append(badReasons, fmt.Sprintf("有一个子域名以「%s」命名，但该网站并不属于这个品牌。", ts.SubdomainBrand))
			riskScore += 20
		} else {
			badReasons = append(badReasons, fmt.Sprintf("有一个子域名使用了品牌名「%s」，但该网站并不属于这个品牌。", ts.SubdomainBrand))
			riskScore += 40
		}
	}
	if ts.PathBrand != "" {
		if resp.ContentData != nil && resp.ContentData.HasLoginForm {
			badReasons = append(badReasons, fmt.Sprintf("链接路径中包含「%s」，页面要求你登录，但该网站并不属于这个品牌。", ts.PathBrand))
			riskScore += 40
		} else {
			badReasons = append(badReasons, fmt.Sprintf("链接路径中包含「%s」，但该网站并不属于这个品牌。", ts.PathBrand))
			riskScore += 20
		}
	}

	// --- 10. Threat Intelligence (PhishTank) ---
	// PhishTank semantics:
	//   valid=true  → URL IS phishing (the actual threat signal)
	//   valid=false → URL is NOT phishing
	//   verified    → whether the report has been reviewed at all (not a phishing indicator)
	if resp.Phishing != nil && resp.Phishing.InDatabase {
		if resp.Phishing.Valid && resp.Phishing.Verified {
			// Reviewed and confirmed as phishing — highest risk.
			badReasons = append(badReasons, "已确认钓鱼：这是一个经过核实的钓鱼链接。")
			riskScore += 200
			if resp.Phishing.Target != "" {
				badReasons = append(badReasons, fmt.Sprintf("被仿冒的目标：%s", resp.Phishing.Target))
			}
		} else if resp.Phishing.Valid && !resp.Phishing.Verified {
			// Reported as phishing but not yet reviewed.
			badReasons = append(badReasons, "这个链接被举报为钓鱼，正在等待社区核实。")
			riskScore += 70
		}
		// valid=false (regardless of verified) → confirmed NOT phishing → no penalty.
		// A verified=true, valid=false entry means PhishTank reviewed it and cleared it.
	}

	// A list naming this exact link outweighs anything good about its site:
	// the site's age and rank say nothing about one page on it.
	confirmed := false

	// Locally held feeds (PhishTank's dump, OpenPhish, URLhaus).
	if fm := resp.ThreatFeeds; fm != nil && fm.Listed {
		sources := strings.Join(fm.Sources, " 和 ")
		if fm.Match == "url" {
			badReasons = append(badReasons, fmt.Sprintf("已确认威胁：这个链接被 %s 列为钓鱼或恶意软件。", sources))
			riskScore += 200
			confirmed = true
		} else {
			badReasons = append(badReasons, fmt.Sprintf("这个网站的其他页面被 %s 列为钓鱼或恶意软件。", sources))
			riskScore += 80
		}
	}

	// Google Web Risk or Safe Browsing (the same lists; one reason is enough).
	for _, g := range []*threatfeeds.GoogleThreatResult{resp.WebRisk, resp.SafeBrowsing} {
		if g != nil && g.Listed {
			// Google's terms require this attribution wherever its verdict is shown.
			badReasons = append(badReasons, fmt.Sprintf("已确认威胁：Google 将这个链接列为%s（由 Google 提供提示）。", webRiskThreatNames(g.ThreatTypes)))
			riskScore += 200
			confirmed = true
			break
		}
	}

	// A target's address carried in the link, for the page to pre-fill.
	if victimAddressInLink(resp.URL) {
		badReasons = append(badReasons, "链接里带有供页面自动填写的邮箱地址，钓鱼页面常用这种方式为每个目标定制登录页。")
		riskScore += 20
	}

	// --- 11. Page Content & Phishing Signals ---
	// Outlook's login address in a form the real Outlook never serves: on PHP,
	// or with its folders run together into one name.
	if (resp.Features.Rank == 0 || resp.Features.Rank > wellKnownRank) && fakeOWAPath(resp.URL) {
		badReasons = append(badReasons, "这个链接模仿了 Outlook 登录页的地址，但真正的 Outlook 页面从来不是这样。钓鱼工具包用它伪装成企业邮箱登录。")
		riskScore += 60
	}
	if why := adobeBait(resp); why != "" {
		badReasons = append(badReasons, why)
		riskScore += 40
		if adobeBaitStrong(resp) {
			riskScore += 30
		}
	}
	// The template's newer wording ("Crescera — AI platform for automated
	// trading") is also how real AI startups title their sites, so it only
	// counts on a fresh, unknown domain.
	if age, known := resp.DomainInfo.Age(); resp.ContentData != nil && resp.Features.Rank == 0 && known && age <= 90 &&
		aiPlatformTitle.MatchString(resp.ContentData.Title) {
		badReasons = append(badReasons, "一个全新网站却以交易或投资的「AI 平台」为标题，这是常见投资骗局的模板。")
		riskScore += 50
	}
	// A page claiming to be the official Chrome download on a site that isn't
	// Google's ships a fake browser.
	if resp.ContentData != nil && chromeDownloadTitle.MatchString(resp.ContentData.Title) &&
		!(resp.Features.Rank > 0 && resp.Features.Rank <= wellKnownRank) &&
		!hostInDomains(hostOf(resp.URL), constants.HighValueBrands["Google"].Domains()) {
		badReasons = append(badReasons, "页面自称是 Google Chrome 官方下载，但这并不是 Google 的网站。")
		riskScore += 100
	}
	if resp.ContentData != nil && resp.Features.Rank == 0 && tradingScamTitle.MatchString(resp.ContentData.Title) {
		badReasons = append(badReasons, "页面使用了虚假加密货币交易平台的标题模板（「AI Trading Platform | Official Website」），这是常见的投资骗局。")
		riskScore += 50
	}
	if resp.ContentData != nil {
		if resp.ContentData.HasLoginForm {
			// Check if domain is established
			isEstablished := resp.Features.Rank > 0 && resp.Features.Rank <= 100000
			age, known := resp.DomainInfo.Age()
			isOld := known && age > oneYearDays

			if !isEstablished && !isOld {
				badReasons = append(badReasons, "可疑：在新的或没有排名的域名上检测到登录表单。")
				riskScore += 50
			} else {
				neutralReasons = append(neutralReasons, "页面包含登录表单。")
			}
			// Mail providers don't host their sign-in pages on free hosting.
			if resp.Features.TLD.IsHostingPlatform && mailLoginTitle.MatchString(resp.ContentData.Title) {
				badReasons = append(badReasons, fmt.Sprintf("在 %s 上出现邮箱登录页，而真正的邮件服务商不会把登录页放在免费托管平台上。", resp.Features.TLD.TLD))
				riskScore += 30
			}
		}

		if resp.ContentData.HasPaymentForm {
			badReasons = append(badReasons, "警告：检测到支付相关字段（信用卡、CVV 等）。")
			riskScore += 30
		}

		if resp.ContentData.HasPersonalForm {
			neutralReasons = append(neutralReasons, "页面要求填写个人信息（地址、电话等）。")
		}

		if resp.ContentData.HasHiddenIframe {
			badReasons = append(badReasons, "警告：检测到隐藏的 iframe（常用于后台窃取凭据或点击劫持）。")
			riskScore += 40
		}

		if resp.ContentData.HasTracking {
			neutralReasons = append(neutralReasons, "检测到后台跟踪元素（1×1 像素）。")
		}

		if pb := resp.ContentData.ProviderBlock; pb != nil {
			switch pb.Reason {
			case "phishing", "malware", "misleading":
				what := "钓鱼"
				switch pb.Reason {
				case "malware":
					what = "恶意软件"
				case "misleading":
					what = "误导性网站"
				}
				badReasons = append(badReasons, fmt.Sprintf("%s 已将本页标记为疑似%s，并以警告页代替显示。", pb.Provider, what))
				riskScore += 100
			default:
				badReasons = append(badReasons, fmt.Sprintf("%s 已将该网站下线。", pb.Provider))
				riskScore += 40
			}
		}

		if sr := resp.ContentData.ScriptRedirect; sr != nil {
			if sr.ToIP {
				badReasons = append(badReasons, "页面用脚本把访客送往一个裸 IP 地址（垃圾信息和钓鱼跳转页常用）。")
				riskScore += 60
			} else if sr.CrossDomain && isUserContentHost(hostOf(resp.URL)) {
				// A file in cloud storage that only forwards visitors is a
				// throwaway hop: it borrows the provider's name to get past
				// filters and points at the real page somewhere else.
				badReasons = append(badReasons, fmt.Sprintf("一个任何人都能上传的文件，会立刻把访客送往另一个网站：%s。", hostOf(sr.Target)))
				riskScore += 60
			} else if sr.CrossDomain {
				badReasons = append(badReasons, fmt.Sprintf("页面会立刻把访客送往另一个网站：%s。", hostOf(sr.Target)))
				riskScore += 30
			}
		}

		if resp.ContentData.BrandCheck.IsMismatch {
			badReasons = append(badReasons, fmt.Sprintf("品牌不符：页面提到「%s」，却托管在非官方域名上。", resp.ContentData.BrandCheck.BrandFound))
			riskScore += 100
		} else if len(resp.ContentData.BrandCheck.DetectedNames) > 0 {
			goodReasons = append(goodReasons, fmt.Sprintf("已核实的品牌匹配：%s", strings.Join(resp.ContentData.BrandCheck.DetectedNames, "、")))
			trustScore += 20
		}

		if resp.ContentData.HasForms {
			for _, form := range resp.ContentData.Forms {
				// Where a form sends its data only matters for what it collects:
				// site search often posts to a search service on another domain.
				// Newsletter sign-ups post to email-marketing platforms
				// (Mailchimp, HubSpot, Eloqua…): only a password or card
				// going there would be out of place.
				marketing := checks.IsMarketingFormService(form.Action) && !form.ContainsPassword && !form.ContainsPayment
				if form.ExternalAction && !marketing {
					switch {
					case form.ContainsPassword || form.ContainsPayment || form.ContainsPersonal:
						badReasons = append(badReasons, "严重：表单会把数据提交到另一个域名（常见的钓鱼手法）。")
						riskScore += 80
					case form.ContainsUserLike:
						badReasons = append(badReasons, "表单会把邮箱地址或用户名发送到另一个域名。")
						riskScore += 20
					}
				}
				if form.ContainsPassword && !resp.SSLInfo.HasTLS {
					badReasons = append(badReasons, "危险：在不安全的连接上检测到密码表单！")
					riskScore += 200
				}
			}
		}
	}

	// --- 12. Certificate ---
	// Browsers refuse a certificate issued for another name with a full-page
	// warning; a real site doesn't ship one. Phishing kits on hacked hosts and
	// shared servers often do.
	if resp.TLSInfo.Present && resp.TLSInfo.HostnameMismatch && servedOverHTTPS(resp) {
		badReasons = append(badReasons, "安全证书不属于这个网站；浏览器打开前会显示警告。")
		riskScore += 40
	}

	// --- 13. Hacked-site Paths ---
	// Phishing kits are dropped into WordPress system folders on compromised
	// sites. Images, styles and scripts live there legitimately; pages don't.
	if isWordPressSystemPage(resp.URL) {
		badReasons = append(badReasons, "页面位于 WordPress 系统目录中，钓鱼工具包常被藏在被入侵网站的这类目录里。")
		riskScore += 30
		if resp.ContentData != nil && resp.ContentData.HasLoginForm {
			badReasons = append(badReasons, "该页面要求访客登录。")
			riskScore += 30
		}
	}

	// --- 14. Program Downloads ---
	// A well-known host can't vouch for a program any account uploaded to it.
	// Its reputation would carry the score to Safe, so cap it instead.
	capAtSuspicious := false
	if ext := downloadedProgram(resp); ext != "" {
		switch host := uploadHost(resp); {
		case host != "":
			badReasons = append(badReasons, fmt.Sprintf("下载的程序（%s）是用户上传到 %s 的，没有任何信息能确认它来自软件发布者。", ext, host))
			capAtSuspicious = true
		case resp.Features.Rank == 0:
			badReasons = append(badReasons, fmt.Sprintf("从不知名网站下载程序（%s）。", ext))
			riskScore += 40
		default:
			neutralReasons = append(neutralReasons, fmt.Sprintf("链接会下载一个程序（%s）。", ext))
		}
	}

	// A short link we couldn't follow could lead anywhere; the shortener's
	// own reputation says nothing about it.
	if resp.ShortLink != nil && !resp.ShortLink.Resolved {
		badReasons = append(badReasons, "无法确认这个短链接指向哪里，所以不能判定为安全。")
		capAtSuspicious = true
	}

	// --- 15. The link as given, when it led elsewhere ---
	// Trust above came from the page the link lands on; what the link itself
	// and the hops on the way get wrong still counts.
	if o := resp.Origin; o != nil {
		badReasons = append(badReasons, o.BadReasons...)
		riskScore += o.Risk
		if o.Confirmed {
			confirmed = true
		}
	}

	// --- Normalize / cap scores ---
	riskScore = clamp(riskScore)
	trustScore = clamp(trustScore)

	// Balanced formula: neutral point is 50.
	// trustScore pulls score up, riskScore pulls it down.
	// A site with no signals either way scores 50 → "Suspicious".
	finalScore := clamp(50 + int(math.Round(float64(trustScore-riskScore)*0.5)))
	if capAtSuspicious && finalScore >= safeThreshold {
		finalScore = safeThreshold - 1
	}
	if confirmed && finalScore >= riskyThreshold {
		finalScore = riskyThreshold - 1
	}

	var verdict string
	switch {
	// High risk, low trust
	case finalScore < riskyThreshold:
		verdict = "Risky"
	// Moderate risk OR conflicting/insufficient signals
	case finalScore < safeThreshold:
		verdict = "Suspicious"
	// Low risk, sufficient trust
	default:
		verdict = "Safe"
	}

	res := Result{
		RiskScore:  riskScore,
		TrustScore: trustScore,
		FinalScore: finalScore,
		Verdict:    verdict,
		Reasons: Reasons{
			NeutralReasons: neutralReasons,
			GoodReasons:    goodReasons,
			BadReasons:     badReasons,
		},
	}

	return res
}

// safeThreshold is the lowest final score that reads as Safe.
const safeThreshold = 65

// riskyThreshold: scores below it are Risky.
const riskyThreshold = 30

// Domain age bands, in days since registration.
const (
	// newDomainDays covers the weeks a phishing domain is typically used.
	newDomainDays = 30
	// establishedAfterDays is the age past which a domain no longer reads as
	// freshly registered for phishing.
	establishedAfterDays = 90
	oneYearDays          = 365
	threeYearsDays       = 3 * 365
	fiveYearsDays        = 5 * 365
)

// ago phrases an age for "registered …": "今天", or "12 天前".
func ago(days int) string {
	if days <= 0 {
		return "今天"
	}
	return agePhrase(days) + "前"
}

// agePhrase renders a day count as "12 天", "3 个月" or "2 年 3 个月".
func agePhrase(days int) string {
	switch {
	case days <= 0:
		return "今天"
	case days < 30:
		return fmt.Sprintf("%d 天", days)
	}
	if days < 365 {
		return fmt.Sprintf("%d 个月", days/30)
	}
	years := days / 365
	if months := (days % 365) / 30; months > 0 {
		return fmt.Sprintf("%d 年 %d 个月", years, months)
	}
	return fmt.Sprintf("%d 年", years)
}

// downloadedProgram returns the executable extension of the file the link
// downloads, from the scanned URL, where it redirects, or the served file name.
func downloadedProgram(resp Response) string {
	candidates := []string{resp.URL, resp.Analysis.RedirectionResult.FinalURL}
	if resp.ContentData != nil {
		candidates = append(candidates, resp.ContentData.FileName)
	}
	for _, c := range candidates {
		if ext := checks.ExecutableExt(c); ext != "" {
			return ext
		}
	}
	return ""
}

// uploadHost returns the user-upload host the link or its download sits on,
// or "" when neither does.
func uploadHost(resp Response) string {
	if resp.Features.TLD.IsHostingPlatform {
		return resp.Features.TLD.TLD
	}
	for _, h := range []string{hostOf(resp.URL), resp.Analysis.RedirectionResult.FinalURLHost} {
		if d := userContentHost(h); d != "" {
			return d
		}
	}
	return ""
}

var wpSystemPathRe = regexp.MustCompile(`(?i)/wp-(?:content|includes|admin)/`)

// staticAssetExts are file types WordPress legitimately serves from its folders.
var staticAssetExts = map[string]struct{}{
	".css": {}, ".js": {}, ".mjs": {}, ".map": {}, ".png": {}, ".jpg": {}, ".jpeg": {},
	".gif": {}, ".webp": {}, ".avif": {}, ".svg": {}, ".ico": {}, ".bmp": {},
	".woff": {}, ".woff2": {}, ".ttf": {}, ".otf": {}, ".eot": {},
	".mp4": {}, ".webm": {}, ".mp3": {}, ".pdf": {}, ".zip": {}, ".json": {}, ".xml": {}, ".txt": {},
}

// isWordPressSystemPage reports whether a URL is a page (not a static file)
// under wp-content, wp-includes or wp-admin. WordPress's own login and admin
// entry points (/wp-login.php, /wp-admin/) aren't matched.
func isWordPressSystemPage(rawURL string) bool {
	u, err := neturl.Parse(rawURL)
	if err != nil || !wpSystemPathRe.MatchString(u.Path) {
		return false
	}
	p := strings.ToLower(strings.TrimSuffix(u.Path, "/"))
	for _, entry := range []string{"/wp-admin", "/wp-admin/index.php", "/wp-admin/admin-ajax.php", "/wp-admin/admin-post.php"} {
		if strings.HasSuffix(p, entry) {
			return false
		}
	}
	_, static := staticAssetExts[strings.ToLower(path.Ext(u.Path))]
	return !static
}

// servedOverHTTPS reports whether visitors reach the page over HTTPS, so a
// certificate problem is something they'd actually run into.
func servedOverHTTPS(resp Response) bool {
	return strings.HasPrefix(resp.URL, "https://") ||
		(strings.HasPrefix(resp.Analysis.RedirectionResult.FinalURL, "https://") &&
			resp.Analysis.RedirectionResult.FinalURLHost == hostOf(resp.URL))
}

// webRiskThreatNames turns Web Risk's threat types into plain words.
func webRiskThreatNames(types []string) string {
	names := map[string]string{
		"SOCIAL_ENGINEERING": "钓鱼",
		"MALWARE":            "恶意软件",
		"UNWANTED_SOFTWARE":  "不受欢迎的软件",
		// Safe Browsing v5 adds this one.
		"POTENTIALLY_HARMFUL_APPLICATION": "有害应用",
	}
	// Most serious first, whatever order Google lists them in.
	order := []string{"SOCIAL_ENGINEERING", "MALWARE", "UNWANTED_SOFTWARE", "POTENTIALLY_HARMFUL_APPLICATION"}
	var out []string
	for _, t := range order {
		if slices.Contains(types, t) {
			out = append(out, names[t])
		}
	}
	if len(out) == 0 {
		return "危险"
	}
	return strings.Join(out, "和")
}

// hostOf returns the host of a URL, or the input itself when it has none (a bare IP).
func hostOf(raw string) string {
	if u, err := neturl.Parse(raw); err == nil && u.Hostname() != "" {
		return u.Hostname()
	}
	return raw
}

func clamp(score int) int {
	return int(math.Max(0, math.Min(100, float64(score))))
}

// brandPageIsBare says why a page on a little-known site using a brand's name
// looks like a phishing setup: it asks for a login or card, or has nothing to
// show (kits often serve a blank page or a directory listing to anyone but
// their targets). A company's own brand-named site or "github.acme.com" has a
// page of its own. "" when none of these hold or the site is ranked. A page
// that couldn't be loaded doesn't count: brands park spare domains
// (airfrance.es, emiratesnbd.ae) that time out or don't resolve.
func brandPageIsBare(resp Response) string {
	c := resp.ContentData
	if resp.Features.Rank > 0 || c == nil {
		return ""
	}
	switch {
	case c.FileName != "":
		return ""
	case c.HasLoginForm || c.HasPaymentForm:
		return "要求你登录或付款"
	case strings.HasPrefix(c.Title, "Index of /"):
		return "只是一个文件列表"
	case strings.TrimSpace(c.Title) == "" && !c.HasForms:
		return "是空白页面"
	}
	return ""
}

// tradingScamTitle matches the title template of fake "AI trading platform"
// investment sites ("Polriksobotiks | AI Trading Platform | Official
// Website", "Immediate Peak ™ | The Official …"). Kits churn out hundreds of
// names on one template.
var tradingScamTitle = regexp.MustCompile(`(?i)\btrading (platform|app|bot|software)\b.{0,20}\bofficial (website|site|gateway|page)\b|^immediate [a-z0-9 ]{2,20}(™|®)?\s*[|:-]\s*(the )?official\b`)

// aiPlatformTitle matches the newer investment-scam template ("Keenhaven
// Wealth — AI platform for systematic investing").
var aiPlatformTitle = regexp.MustCompile(`(?i)\s[—–|-]\s*ai(-powered)? (trading |investment )?platform for\b`)

// chromeDownloadTitle matches a title claiming the official Chrome download,
// in Chinese (谷歌浏览器 官方/官网/下载) or English.
var chromeDownloadTitle = regexp.MustCompile(`(?i)(谷歌浏览器|google ?chrome|chrome浏览器).{0,20}(官方|官网|下载)|(官方|官网).{0,10}(谷歌浏览器|google ?chrome)|^(download )?google chrome (official|download)`)

// mailLoginTitle matches page titles of email sign-in pages.
var mailLoginTitle = regexp.MustCompile(`(?i)\b(web ?mail|mail ?(authentication|verification|login|portal)|e-?mail (verification|authentication|login|sign[ -]?in|portal)|mailbox|password protect(ion|ed)|verify your (e-?mail|mailbox|account)|sign ?in ?to (your )?(web)?mail|quarantine portal|e-?mail quarantine)\b`)

// emailToken is a whole email address; mailMergeField is a mass mailer's
// placeholder left unfilled ("#[[-Email-]]").
var (
	emailToken     = regexp.MustCompile(`(?i)^[a-z0-9._%+-]{1,64}@[a-z0-9-]+(\.[a-z0-9-]+)*\.[a-z]{2,}$`)
	mailMergeField = regexp.MustCompile(`(?i)\[\[-?e-?mail-?\]\]|\{\{\s*e-?mail\s*\}\}|##e-?mail##|%e-?mail%`)
)

// victimAddressInLink reports whether a link carries an email address in its
// fragment or, base64-encoded, in its query. A plain address in the query
// counts only under an odd name (?eta=, ?uid=): newsletter and account links
// carry one honestly as ?email= or ?e=, and on unsubscribe pages.
func victimAddressInLink(rawURL string) bool {
	u, err := neturl.Parse(rawURL)
	if err != nil {
		return false
	}
	if mailMergeField.MatchString(u.Fragment) {
		return true
	}
	split := func(s string) []string {
		return strings.FieldsFunc(s, func(r rune) bool { return strings.ContainsRune("#&=/?!;,", r) })
	}
	for _, tok := range split(u.Fragment) {
		if emailToken.MatchString(tok) || emailToken.MatchString(fromBase64(tok)) {
			return true
		}
	}
	mailingPage := mailingPath.MatchString(u.Path)
	for name, vals := range u.Query() {
		_, honest := honestAddressParams[strings.ToLower(name)]
		for _, v := range vals {
			if emailToken.MatchString(fromBase64(v)) || (!honest && !mailingPage && emailToken.MatchString(v)) {
				return true
			}
		}
	}
	return false
}

// honestAddressParams are query names under which sites pass an address
// openly; mailingPath matches pages that take one by design.
var (
	honestAddressParams = map[string]struct{}{
		"email": {}, "e-mail": {}, "e_mail": {}, "mail": {}, "e": {}, "em": {},
		"emailaddress": {}, "email_address": {}, "address": {}, "recipient": {},
		"subscriber": {}, "to": {}, "from": {}, "user": {}, "username": {},
		"login": {}, "contact": {}, "invite": {}, "invitee": {}, "q": {},
	}
	mailingPath = regexp.MustCompile(`(?i)unsubscribe|opt-?out|preferences|subscri|newsletter|invite|verify|confirm|reset`)
)

// fromBase64 decodes s as standard or URL-safe base64, padded or not, or
// returns "".
func fromBase64(s string) string {
	// Kits sometimes mark the encoded part ("$bWFu…").
	s = strings.TrimLeft(s, "$~!*")
	if len(s) < 8 {
		return ""
	}
	for _, enc := range []*base64.Encoding{base64.StdEncoding, base64.URLEncoding, base64.RawStdEncoding, base64.RawURLEncoding} {
		if b, err := enc.DecodeString(s); err == nil {
			return string(b)
		}
	}
	return ""
}

// hostingRegistration reports whether the domain info looked up for a site on
// a hosting platform is really the platform's: a customer's subdomain, or a
// user's page on a link-in-bio service. The platform's own site keeps its age.
func hostingRegistration(resp Response) bool {
	host := strings.TrimPrefix(hostOf(resp.URL), "www.")
	if _, ok := constants.UserPageHosts[host]; ok {
		return !isBareURL(resp.URL)
	}
	return host != resp.Features.TLD.TLD && strings.HasSuffix(host, "."+resp.Features.TLD.TLD)
}

// kitHexToken is a long run of hex, as kits put in paths and session
// parameters; owaOnPHP is an Outlook Web Access path served by PHP, which
// real OWA (ASP.NET) never is.
var (
	kitHexToken = regexp.MustCompile(`(?i)[0-9a-f]{16,}`)
	owaOnPHP    = regexp.MustCompile(`(?i)/owa/.*\.php\b`)
	// owaDotted is Outlook's login path written with dots as one folder name
	// (christianbased.org/owa.auth.logon.aspx/). Real OWA is /owa/auth/logon.aspx.
	owaDotted = regexp.MustCompile(`(?i)owa\.auth\.logon`)
)

// fakeOWAPath reports whether a link copies Outlook Web Access's login path in
// a form real OWA never uses.
func fakeOWAPath(rawURL string) bool {
	u, err := neturl.Parse(rawURL)
	return err == nil && (owaOnPHP.MatchString(u.Path) || owaDotted.MatchString(u.Path))
}

// kitLoginPath reports whether a link is a sign-in page that looks planted:
// a sign-in path carrying a long hex token in the path or a session-style
// parameter, or OWA on PHP. Reset and confirmation links carry tokens too,
// under names like token= or code=, so those don't count.
func kitLoginPath(rawURL string) bool {
	u, err := neturl.Parse(rawURL)
	if err != nil {
		return false
	}
	if fakeOWAPath(rawURL) {
		return true
	}
	if !loginPathWords.MatchString(u.Path) {
		return false
	}
	if kitHexToken.MatchString(u.Path) {
		return true
	}
	for name, vals := range u.Query() {
		switch strings.ToLower(name) {
		case "token", "code", "key", "reset", "confirm", "verify", "signature", "sig", "state", "nonce":
			continue
		}
		for _, v := range vals {
			if kitHexToken.MatchString(v) {
				return true
			}
		}
	}
	return false
}

// adobeFile is a page file named after Adobe (Adobe_Installer.html,
// adobekr.html); pdfViewerTitle is the title fake document viewers use.
var (
	adobeFile      = regexp.MustCompile(`(?i)adobe[^/]*\.(s?html?|php)$`)
	pdfViewerTitle = regexp.MustCompile(`(?i)^\s*(adobe )?(pdf|document) ?viewer\s*$`)
)

// adobeSigns reports whether the page file is named after Adobe and whether
// the page is titled as a PDF viewer.
func adobeSigns(path string, resp Response) (named, titled bool) {
	return adobeFile.MatchString(path), resp.ContentData != nil && pdfViewerTitle.MatchString(resp.ContentData.Title)
}

// adobeBaitStrong is adobeBait with both signs at once: a file named after
// Adobe that is titled as a PDF viewer. Real sites don't name their pages
// after Adobe's product and then pose as its viewer.
func adobeBaitStrong(resp Response) bool {
	if adobeBait(resp) == "" {
		return false
	}
	u, err := neturl.Parse(resp.URL)
	if err != nil {
		return false
	}
	named, titled := adobeSigns(u.Path, resp)
	return named && titled
}

// adobeBait describes a fake Adobe "PDF viewer" page, a common lure for a
// login to "open the document": a page file named after Adobe, or titled
// "PDF Viewer", on a little-known site that isn't Adobe's. "" otherwise.
func adobeBait(resp Response) string {
	if resp.Features.Rank > 0 && resp.Features.Rank <= wellKnownRank {
		return ""
	}
	host := hostOf(resp.URL)
	if hostInDomains(host, constants.HighValueBrands["Adobe"].Domains()) {
		return ""
	}
	u, err := neturl.Parse(resp.URL)
	if err != nil {
		return ""
	}
	named, titled := adobeSigns(u.Path, resp)
	switch {
	case named && titled:
		return "页面以 Adobe 命名、标题是 PDF 阅读器，却不在 Adobe 的网站上：这是诱导你登录以「打开文档」的常见手法。"
	case named:
		return "页面以 Adobe 命名，却不在 Adobe 的网站上，这是常见诱饵（假冒的 PDF 阅读器或安装程序）。"
	case titled && resp.ContentData.HasForms:
		return "不知名网站上出现带表单的「PDF 阅读器」页面：这是诱导你登录以「打开文档」的常见手法。"
	}
	return ""
}
