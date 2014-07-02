package main

import (
	"fmt"
	"github.com/sysr-q/assets"
)

func main() {
	// The world's greatest example.
	script := string(assets.MustRead("hackers.txt"))
	fmt.Println(script)
}
