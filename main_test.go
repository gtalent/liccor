package main

import (
	"os"
	"testing"
)

func Test_Liccor_Version(t *testing.T) {
	os.Args[0] = "liccor"
	os.Args[1] = "--version"
	main()
}
