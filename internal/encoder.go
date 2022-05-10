package pkg

func Base62Encode(deci int) string {
	s := "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	hashStr := ""
	for deci > 0 {
		hashStr = string(s[deci%62]) + hashStr
		deci /= 62
	}
	return hashStr
}
