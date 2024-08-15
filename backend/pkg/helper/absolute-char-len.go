package helper

// AbsoluteCharLen Strip or fill a string to declared `l` length
func AbsoluteCharLen(s string, l int) string {
	for len(s) < l {
		s = s + "*"
	}
	if len(s) > l {
		s = s[:l]
	}
	return s
}
