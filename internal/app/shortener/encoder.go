package shortener

func Base62Encode(deci uint64) string {
	s := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	hashStr := ""
	for deci > 0 {
		hashStr = string(s[deci%62]) + hashStr
		deci /= 62
	}
	return hashStr
}
