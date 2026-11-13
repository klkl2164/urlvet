package checks

import "github.com/klkl2164/urlvet/internal/constants"

func IsUrlShortener(domain string) bool {
	_, ok := constants.URLShorteners[domain]
	return ok
}
