package main

import myos "os"

func main() {
	if true {
		myos.Exit(0) // want "direct call to os.Exit in main function is not allowed; use log.Fatal or return an error instead"
	}
}
