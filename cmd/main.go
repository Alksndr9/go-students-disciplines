package main

import (
	"fmt"

	"gitgub.com/Alksndr9/go-students-disciplines/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TO-DO: logger
	// TO-DO: bd
	// TO-DO: router gin

}
