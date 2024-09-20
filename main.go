package main

import (
	"context"
	"fmt"
	"lang_tree/search"
	"os"
)

func main() {
	var tree = search.NewLanguageTree(context.Background())
	err := tree.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// In a for loop parse input,
	// either parse user input, or eccept command line parameters
}
