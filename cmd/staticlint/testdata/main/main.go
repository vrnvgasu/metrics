package main

import "os"

func main() {
	if true {
		os.Exit(0) // want "direct call to os.Exit in main function is not allowed; use log.Fatal or return an error instead"
	}

	test()
}

func test() {
	if true {
		os.Exit(0)
	}
}
