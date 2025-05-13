package pkg

func TruncateString(s string, i int) string {
	if len(s) < i {
		return s
	} else {
		return s[:i]
	}
}
