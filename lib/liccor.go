/*
   Copyright 2011-2016 gtalent2@gmail.com

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

package liccor

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/paulvollmer/go-verbose"
)

// Version store the version as string
const Version = "1.8.2"

// Liccor the license corrector
type Liccor struct {
	Log               verbose.Verbose
	License           string
	LicenseBeforeText string
	LicenseAfterText  string
}

// New initialize and return a new Liccor instance
func New() *Liccor {
	l := Liccor{}
	l.Log = *verbose.New(os.Stdout, false)
	return &l
}

// DefaultLicenseFile store the default file to search for
const DefaultLicenseFile = ".liccor"

// FindLicense search for a license file
func (l *Liccor) FindLicense(dir, licenseFile string) (string, error) {
	l.Log.Printf("Search for a license file at directory '%s'\n", dir)

	d, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("Could not find license file. %v", err)
	}
	for _, v := range d {
		filename := v.Name()
		// search the license file
		if filename == licenseFile || filename == DefaultLicenseFile || filename == "LICENSE" || filename == "LICENSE.txt" {
			licenseData, err := ioutil.ReadFile(dir + "/" + v.Name())
			if err != nil {
				err = fmt.Errorf("Could not access " + filename + " file")
			}
			l.Log.Printf("License file '%s' found...\n", filename)
			return string(licenseData), err
		}
	}

	return l.FindLicense(dir+"./.", licenseFile)
}

// FindSrcFiles search for source files
func (l *Liccor) FindSrcFiles(dir string) ([]string, error) {
	l.Log.Printf("Search source files at '%s'\n", dir)

	d, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var output []string
	for _, v := range d {
		if v.IsDir() {
			// ignore .git dir
			if v.Name() != ".git" {
				files, err2 := l.FindSrcFiles(dir + "/" + v.Name())
				if err2 != nil {
					return output, err2
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
			case SuffixGO, SuffixC, SuffixCPP, SuffixCXX, SuffixH, SuffixHPP, SuffixJAVA, SuffixJS:
				srcPath := dir + "/" + v.Name()
				output = append(output, srcPath)
				l.Log.Printf("Found source '%s'\n", srcPath)
			}
		}
	}
	return output, err
}

// HasLicense check if sourcecode has a license at the top
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

// Correct a source file license
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
					i++
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

// Process run the liccor magic
func (l *Liccor) Process() {
	licenseData, err := l.FindLicense(".", l.License)
	if err != nil {
		fmt.Println(err)
		return
	}
	licenseData = licenseData[0 : len(licenseData)-1]
	if l.LicenseBeforeText != "" {
		l.Log.Printf("License before text set to '%s'\n", l.LicenseBeforeText)
		licenseData = l.LicenseBeforeText + "\n" + licenseData
	}
	if l.LicenseAfterText != "" {
		l.Log.Printf("License after text set to '%s'\n", l.LicenseAfterText)
		licenseData = licenseData + "\n" + l.LicenseAfterText
	}
	//fmt.Println("License", licenseData)

	lics := make(map[string]string)
	clike := "/*\n * " + strings.Replace(string(licenseData), "\n", "\n * ", -1) + "\n */\n"
	lics["c-like"] = strings.Replace(clike, "\n * \n", "\n *\n", -1)
	lics["go"] = func() string {
		golic := "/*\n   " + strings.Replace(string(licenseData), "\n", "\n   ", -1) + "\n*/\n"
		golic = strings.Replace(golic, "\n   \n", "\n\n", -1)
		return golic
	}()

	files, err := l.FindSrcFiles(".")
	if err != nil {
		return
	}

	allSuccess := true
	for i := 0; i < len(files); i++ {
		pt := strings.LastIndex(files[i], ".")
		lic := ""
		//determine how to format the license
		switch files[i][pt:] {
		case SuffixGO:
			lic = lics["go"]
		case SuffixC, SuffixCPP, SuffixCXX, SuffixH, SuffixHPP, SuffixJAVA, SuffixJS:
			lic = lics["c-like"]
		}
		changed, err := l.Correct(files[i], lic)
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
