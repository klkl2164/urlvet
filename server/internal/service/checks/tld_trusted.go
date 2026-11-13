package checks

import "github.com/klkl2164/urlvet/internal/constants"

func IsTrustedTld(domain string) (bool, bool, string) {
	tld, icann := GetTld(domain)
	_, ok := constants.TrustedTLDs[tld]

	return ok, icann, tld
}
