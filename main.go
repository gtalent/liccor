/*
   This is a license.
 */
package main

import (
	"io/ioutil"
	"flag"
	"fmt"
	"strings"
)

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		return
	}
	licenseData, err := ioutil.ReadFile(flag.Arg(0))
	if err != nil {
		return
	}
	licenseData = licenseData[0:len(licenseData) - 1]
	lics := make(map[string]string)
	lics["c-like"] = "/*\n * " + strings.Replace(string(licenseData), "\n", "\n * ", -1) + "\n */\n"
	lics["go"] = func () string {
		golic := "/*\n   " + strings.Replace(string(licenseData), "\n", "\n   ", -1) + "\n */\n"
		golic = strings.Replace(golic, "\n   \n", "\n\n", -1)
		return golic
	}()
	for i := 1; i < flag.NArg(); i++ {
		pt := strings.LastIndex(flag.Arg(i), ".")
		lic := ""
		//determine how to format the license
		if pt == -1 {
			fmt.Println("Skipping", flag.Arg(i))
			continue
		}
		switch flag.Arg(i)[pt:] {
		case ".go":
			fmt.Print("Correcting ", flag.Arg(i), "...")
			lic = lics["go"]
		case ".c", ".cpp", ".cxx", ".h", ".hpp", ".java":
			fmt.Print("Correcting ", flag.Arg(i), "...")
			lic = lics["c-like"]
		default:
			fmt.Println("Skipping", flag.Arg(i))
			continue
		}
		if !correct(flag.Arg(i), lic) {
			fmt.Println("\tFailure!")
			continue
		}
		fmt.Println("\tSuccess!")
	}
}

func hasLicense(file []byte) (bool, int) {
	for i, c := range file {
		switch c {
		case ' ', '\t', '\n':
			continue
		case '/':
			i++
			if len(file) > i && file[i] == '*' {
				return true, i
			}
		default:
			return false, -1
		}
	}
	return false, -1
}

func correct(path, license string) bool {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}
	if hasLicense, licenseStart := hasLicense(file); hasLicense {
		//remove old license
		for i := licenseStart; i < len(file); i++ {
			if file [i] == '*' && file[i + 1] == '/' {
				i += 2
				if file[i] == '\n' {
					i += 1
				}
				file = file[i:len(file)]
				break
			}
		}
	}
	file = []byte(license + string(file))
	return ioutil.WriteFile(path, file, 0) == nil
}
