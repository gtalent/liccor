/*
   Copyright 2011-2017 gtalent2@gmail.com

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

package lib

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type Liccor struct {
	Log               Logger
	Source            string
	License           string
	LicenseBeforeText string
	LicenseAfterText  string
}

const (
	DEFAULT_LICENSE_FILE = ".liccor"
)

func (l *Liccor) FindLicense(dir, licenseFile string) (string, error) {
	l.Log.Verbose("Search for a license file at directory '" + dir + "'")

	d, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("Could not find license file")
	}
	for _, v := range d {
		filename := v.Name()
		// search the license file
		if filename == licenseFile || filename == DEFAULT_LICENSE_FILE || filename == "LICENSE" || filename == "LICENSE.txt" {
			licenseData, err := ioutil.ReadFile(dir + "/" + v.Name())
			if err != nil {
				err = fmt.Errorf("Could not access " + filename + " file")
			}
			l.Log.Verbose("License file '" + filename + "' found...")
			return string(licenseData), err
		}
	}

	return l.FindLicense(dir+"./.", licenseFile)
}

func (l *Liccor) FindSrcFiles(dir string) ([]string, error) {
	l.Log.Verbose("Search source files at '" + dir + "'")

	d, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	output := make([]string, 0)
	for _, v := range d {
		if v.IsDir() {
			// ignore .git dir
			if v.Name() != ".git" {
				files, err := l.FindSrcFiles(dir + "/" + v.Name())
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
				l.Log.Verbose("Found source '" + srcPath + "'")
			}
		}
	}
	return output, err
}

func (l *Liccor) HasLicense(file string) (bool, int) {
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

func (l *Liccor) Correct(path, license string) (bool, error) {
	input, err := ioutil.ReadFile(path)
	if err != nil {
		return false, err
	}
	file := string(input)
	orig := file
	if hasLicense, licenseStart := l.HasLicense(file); hasLicense {
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

func (l *Liccor) Process() {
	licenseData, err := l.FindLicense(".", l.License)
	if err != nil {
		fmt.Println(err)
		return
	}
	licenseData = licenseData[0 : len(licenseData)-1]
	if l.LicenseBeforeText != "" {
		l.Log.Verbose("License before text set to '" + l.LicenseBeforeText + "'")
		licenseData = l.LicenseBeforeText + "\n" + licenseData
	}
	if l.LicenseAfterText != "" {
		l.Log.Verbose("License after text set to '" + l.LicenseAfterText + "'")
		licenseData = licenseData + "\n" + l.LicenseAfterText
	}

	lics := make(map[string]string)
	clike := "/*\n * " + strings.Replace(string(licenseData), "\n", "\n * ", -1) + "\n */\n"
	lics["c-like"] = strings.Replace(clike, "\n * \n", "\n *\n", -1)
	lics["go"] = func() string {
		golic := "/*\n   " + strings.Replace(string(licenseData), "\n", "\n   ", -1) + "\n*/\n"
		golic = strings.Replace(golic, "\n   \n", "\n\n", -1)
		return golic
	}()

	files, err := l.FindSrcFiles(l.Source)
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
		changed, err := l.Correct(files[i], lic)
		if changed {
			var file string
			if files[i][:2] == "./" {
				file = files[i][2:]
			} else {
				file = files[i]
			}
			if err != nil {
				fmt.Println("Correcting '" + file + "'... Failure!")
				allSuccess = false
			} else {
				fmt.Println("Correcting '" + file + "'... Success!")
			}
		}
	}
	if allSuccess {
		fmt.Println("All files up to date!")
	}
}
