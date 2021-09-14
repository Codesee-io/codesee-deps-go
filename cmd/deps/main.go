package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Codesee-io/codesee-deps-go/pkg/errutils"
	"github.com/Codesee-io/codesee-deps-go/pkg/links"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: codesee-deps-go <directory>")
		os.Exit(1)
	}

	root := os.Args[1]
	l, err := links.DetermineLinks(root)
	if err != nil {
		errutils.Fatal(err)
	}

	out, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		errutils.Fatal(err)
	}
	fmt.Println(string(out))
}
