package main

import (
	"fmt"

	"example.com/m/config"
)

func main() {
	fmt.Println(config.LoadConfig())
}
