package main

import (
	"fmt"
	"os"

	"github.com/vladimirvivien/echo"
)

func main() {
	e := echo.New()
	fmt.Println("Running binary release...")

	// ensure we're in travis
	if e.Eval("${TRAVIS}") == "" {
		fmt.Println("This script can only run in Travis CI environment")
		os.Exit(1)
	}

	// release on tag push
	if !e.Empty("${TRAVIS_TAG}") {
		fmt.Println("Releasing binaries with Goreleaser")
		result := e.Run("curl -sL https://git.io/goreleaser | bash")
		if result != "" {
			fmt.Println("Goreleaser may have failed:", result)
			os.Exit(1)
		}
	} else {
		fmt.Println("Not a tag pushed, skipping release")
	}
}
