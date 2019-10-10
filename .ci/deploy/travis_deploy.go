package main

import (
	"fmt"
	"os"

	"github.com/vladimirvivien/echo"
)

func main() {
	e := echo.New()
	fmt.Println("Running binary release...")

	//ensure we're travis and configuree
	if e.Eval("${GITHUB_TOKEN}") == "" {
		fmt.Println("missing GITHUB_TOKEN env")
		os.Exit(1)
	}

	// release on tag push
	if e.Eval("${PUBLISH}") == "" {
		fmt.Println("PUBLISH not set, skipping binary publish")
		e.Runout("goreleaser --rm-dist --skip-validate --skip-publish")
	} else {
		fmt.Println("Publishing binaries with goreleaser")
		e.Runout("goreleaser --rm-dist")
	}
}
