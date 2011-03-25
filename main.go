package main

import (
	"io/ioutil"
	"flag"
	"fmt"
)

var newLicense string

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		return
	}
	license, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		return
	}
	for i := 1; i < flag.NArg(); i++ {
		fmt.Print("Correcting ", flag.Arg(i), "...")
		if !correct(flag.Arg(i), string(license)) {
			fmt.Println("\tFailure!")
			continue
		}
		fmt.Println("\tSuccess!")
	}
}

func correct(path, license string) bool {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}

	var hasLicense bool

	for i, c := range file {
		switch c {
		case ' ', '\t', '\n':
			continue
		case '/':
			if len(file) > i && file[i+1] == '*' {
				hasLicense = true
				break
			}
		default:
			hasLicense = false
			break
		}
	}

	var previous byte
	var lend int
	if hasLicense {
		for i, c := range file {
			if previous == '*' && c == '/' {
				lend = i
				break
			}
			previous = c
		}
		file = file[lend:len(file)]
	}
	file = []byte(license + string(file))
	return ioutil.WriteFile(path, file, 0) == nil
}
