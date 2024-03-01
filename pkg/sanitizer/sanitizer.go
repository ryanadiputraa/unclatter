package sanitizer

import "github.com/microcosm-cc/bluemonday"

type Sanitizer interface {
	Sanitize(s string) string
}

type sanitize struct {
	policy *bluemonday.Policy
}

func NewSanitizer() Sanitizer {
	return &sanitize{
		bluemonday.UGCPolicy(),
	}
}

func (sn *sanitize) Sanitize(s string) string {
	return sn.policy.Sanitize(s)
}
