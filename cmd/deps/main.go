package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Codesee-io/codesee-deps-go/pkg/errutils"
	"github.com/Codesee-io/codesee-deps-go/pkg/links"
)

var (
	version = "dev"
	commit  = "dev"
	date    = time.Now().Format(time.RFC3339)
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: codesee-deps-go <directory>")
		os.Exit(1)
	}

	if os.Args[1] == "-v" || os.Args[1] == "--version" {
		fmt.Printf("codesee-deps-go version %s\ncommit: %s\nbuilt at: %s\n", version, commit, date)
		os.Exit(0)
	}

	root := os.Args[1]
	l, err := links.DetermineLinks(root)
	if err != nil {
		errutils.Fatal(err)
	}

	out, err := json.Marshal(l)
	if err != nil {
		errutils.Fatal(err)
	}
	fmt.Println(string(out))
}
