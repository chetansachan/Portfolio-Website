package main

import (
	"fmt"
	"portfolio-website/internals/config"
)

func main() {
	fmt.Print(config.Load())
}
