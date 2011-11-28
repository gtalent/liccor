/*
   Copyright 2011 gtalent2@gmail.com

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
	"io/ioutil"
	"fmt"
	"os"
	"strings"
)

func findLicense(dir string) (string, os.Error) {
	d, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	for _, v := range d {
		if v.Name == ".copyright" {
			licenseData, err := ioutil.ReadFile(dir + "/" + v.Name)
			return string(licenseData), err
		}
	}
	return findLicense(dir + "./.")
}

func findSrcFiles(dir string) ([]string, os.Error) {
	l, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	output := make([]string, 0)
	for _, v := range l {
		if v.IsDirectory() {
			files, err := findSrcFiles(dir + "/" + v.Name)
			if err != nil {
				return output, err
			}
			for _, v2 := range files {
				output = append(output, v2)
			}
		} else {
			pt := strings.LastIndex(v.Name, ".")
			//determine how to format the license
			if pt == -1 {
				continue
			}
			switch v.Name[pt:] {
			case ".go", ".c", ".cpp", ".cxx", ".h", ".hpp", ".java":
				output = append(output, dir + "/" + v.Name)
			}
		}
	}
	return output, err
}

func main() {
	licenseData, err := findLicense(".")
	if err != nil {
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
	for i := 0; i < len(files); i++ {
		pt := strings.LastIndex(files[i], ".")
		lic := ""
		//determine how to format the license
		switch files[i][pt:] {
		case ".go":
			fmt.Print("Correcting ", files[i], "...")
			lic = lics["go"]
		case ".c", ".cpp", ".cxx", ".h", ".hpp", ".java":
			fmt.Print("Correcting ", files[i], "...")
			lic = lics["c-like"]
		default:
			fmt.Println("Skipping", files[i])
			continue
		}
		if !correct(files[i], lic) {
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
	file = []byte(license + string(file))
	return ioutil.WriteFile(path, file, 0) == nil
}
