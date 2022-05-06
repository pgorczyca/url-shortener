package main

import (
	"fmt"

	"github.com/pgorczyca/url-shortener/pkg"
)

func main() {
	fmt.Println(pkg.Base62Encode(9234526231))
}
