/*
   Copyright 2011-2014 gtalent2@gmail.com

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const (
	DEFAULT_LICENSE_FILE = ".liccor"
	VERSION              = "liccor 1.8 (go1)"
	// list of file extensions
	SUFFIX_GO   = ".go"
	SUFFIX_C    = ".c"
	SUFFIX_CPP  = ".cpp"
	SUFFIX_CXX  = ".cxx"
	SUFFIX_H    = ".h"
	SUFFIX_HPP  = ".hpp"
	SUFFIX_JAVA = ".java"
	SUFFIX_JS   = ".js"
)

var (
	flagLicenseFile string
	flagVerbose     bool
	showVersion     bool
)

func verboseLog(msg string) {
	if flagVerbose {
		fmt.Println(msg)
	}
}

func findLicense(dir string) (string, error) {
	verboseLog("Search for a license file at directory '" + dir + "'")

	d, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("Could not find license file")
	}
	for _, v := range d {
		filename := v.Name()
		// search the license file
		if filename == flagLicenseFile || filename == DEFAULT_LICENSE_FILE || filename == "LICENSE" || filename == "LICENSE.txt" {
			licenseData, err := ioutil.ReadFile(dir + "/" + v.Name())
			if err != nil {
				err = fmt.Errorf("Could not access " + filename + " file")
			}
			verboseLog("License file '" + filename + "' found...")
			return string(licenseData), err
		}
	}

	return findLicense(dir + "./.")
}

func findSrcFiles(dir string) ([]string, error) {
	verboseLog("Search source files at '" + dir + "'")

	l, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	output := make([]string, 0)
	for _, v := range l {
		if v.IsDir() {
			// ignore .git dir
			if v.Name() != ".git" {
				files, err := findSrcFiles(dir + "/" + v.Name())
				if err != nil {
					return output, err
				}
				for _, v2 := range files {
					output = append(output, v2)
				}
			}
		} else {
			pt := strings.LastIndex(v.Name(), ".")
			if pt == -1 {
				continue
			}
			switch v.Name()[pt:] {
			case SUFFIX_GO, SUFFIX_C, SUFFIX_CPP, SUFFIX_CXX, SUFFIX_H, SUFFIX_HPP, SUFFIX_JAVA, SUFFIX_JS:
				srcPath := dir + "/" + v.Name()
				output = append(output, srcPath)
				verboseLog("Found source '" + srcPath + "'")
			}
		}
	}
	return output, err
}

func hasLicense(file string) (bool, int) {
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

func correct(path, license string) (bool, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	file := string(input)
	orig := file
	if hasLicense, licenseStart := hasLicense(file); hasLicense {
		//remove old license
		for i := licenseStart; i < len(file); i++ {
			if file[i] == '*' && file[i+1] == '/' {
				i += 2
				if file[i] == '\n' {
					i += 1
				}
				file = file[i:len(file)]
				break
			}
		}
	}
	file = license + file
	output := []byte(file)
	if file != orig {
		err = ioutil.WriteFile(path, output, 0)
		return true, err
	}
	return false, nil
}

func version() {
	if showVersion {
		println(VERSION)
		os.Exit(0)
	}
}

func init() {
	flag.StringVar(&flagLicenseFile, "license", DEFAULT_LICENSE_FILE, "the name of the license file")
	flag.StringVar(&flagLicenseFile, "l", DEFAULT_LICENSE_FILE, "shortcut for license")
	flag.BoolVar(&flagVerbose, "verbose", false, "print verbose output")
	flag.BoolVar(&flagVerbose, "v", false, "shortcut for verbose")
	flag.BoolVar(&showVersion, "version", false, "version of liccor")
	flag.Usage = func() {
		fmt.Print("\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExample usage:")
		fmt.Println("  ./liccor -verbose")
		fmt.Print("\n\n")
	}
}

func main() {
	flag.Parse()
	version()

	licenseData, err := findLicense(".")
	if err != nil {
		fmt.Println(err)
		return
	}
	licenseData = licenseData[0 : len(licenseData)-1]
	lics := make(map[string]string)
	lics["c-like"] = "/*\n * " + strings.Replace(string(licenseData), "\n", "\n * ", -1) + "\n */\n"
	lics["go"] = func() string {
		golic := "/*\n   " + strings.Replace(string(licenseData), "\n", "\n   ", -1) + "\n*/\n"
		golic = strings.Replace(golic, "\n   \n", "\n\n", -1)
		return golic
	}()

	files, err := findSrcFiles(".")
	if err != nil {
		return
	}

	allSuccess := true
	for i := 0; i < len(files); i++ {
		pt := strings.LastIndex(files[i], ".")
		lic := ""
		//determine how to format the license
		switch files[i][pt:] {
		case SUFFIX_GO:
			lic = lics["go"]
		case SUFFIX_C, SUFFIX_CPP, SUFFIX_CXX, SUFFIX_H, SUFFIX_HPP, SUFFIX_JAVA, SUFFIX_JS:
			lic = lics["c-like"]
		}
		changed, err := correct(files[i], lic)
		if changed {
			if err != nil {
				fmt.Println("Correcting '" + files[i][2:] + "'... Failure!")
				allSuccess = false
			} else {
				fmt.Println("Correcting '" + files[i][2:] + "'... Success!")
			}
		}
	}
	if allSuccess {
		fmt.Println("All files up to date!")
	}
}
