package main

import (
	"fmt"
	"os"
)

func run() error {
	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
