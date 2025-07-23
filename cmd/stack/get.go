package main

import (
	"fmt"
	"os"

	"github.com/canonical/go-snapctl"
)

func get(key string) {
	value, err := snapctl.Get(key).Run()
	if err != nil {
		fmt.Printf("Error getting value: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(value)
}
