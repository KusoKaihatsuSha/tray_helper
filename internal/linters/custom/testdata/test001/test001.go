package pkg1

import "os"

func pkg1() {
	os.Exit(0) // want "function 'Exit' not permit in 'main' file and 'main' function"
}
