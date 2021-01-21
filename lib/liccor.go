/*
   Copyright 2011 - 2021 gary@drinkingtea.net

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

	"github.com/monochromegane/go-gitignore"
	"github.com/paulvollmer/go-verbose"

	"gopkg.in/yaml.v2"
)

// Version store the version as string
const Version = "1.9.2"

type liccorFile struct {
	Source          []string `yaml:"source"`
	CopyrightNotice string   `yaml:"copyright_notice"`
	Ignore          string   `yaml:"ignore"`
}

// Liccor the license corrector
type Liccor struct {
	Log              verbose.Verbose
	Source           []string
	NoticeBeforeText string
	NoticeAfterText  string
	copyrightNotice  string
	ignore           gitignore.IgnoreMatcher
}

// New initialize and return a new Liccor instance
func New() *Liccor {
	l := Liccor{}
	// ignore error, it's ok if there is no gitignore
	l.ignore, _ = gitignore.NewGitIgnore(".gitignore")
	l.Log = *verbose.New(os.Stdout, false)
	l.Source = []string{"."}
	return &l
}

// Load liccor file search for a license file
func (l *Liccor) LoadConfig(dir, liccorFileName string) error {
	l.Log.Printf("Search for a license file at directory '%s'\n", dir)

	d, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Could not find license file. %v", err)
	}

	fileName := ""
	if liccorFileName == "" { // no specified liccor file, search for default file names
		for _, v := range d {
			// search the license file
			if v.Name() == ".liccor" || v.Name() == ".liccor.yml" || v.Name() == ".liccor.yaml" {
				fileName = v.Name()
				break
			}
		}
	} else if _, err := os.Stat(liccorFileName); err == nil {
		fileName = liccorFileName
	}

	switch {
	case fileName == "":
		break
	case strings.HasSuffix(fileName, ".yaml"), strings.HasSuffix(fileName, ".yml"):
		var lf liccorFile
		data, err := ioutil.ReadFile(dir + "/" + fileName)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(data, &lf)
		if err != nil {
			return err
		}

		l.copyrightNotice = lf.CopyrightNotice
		l.Source = lf.Source
		if lf.Ignore != "" {
			l.ignore = gitignore.NewGitIgnoreFromReader("", strings.NewReader(lf.Ignore))
		}
		return err
	default:
		copyrightNotice, err := ioutil.ReadFile(dir + "/" + fileName)
		if err != nil {
			err = fmt.Errorf("Could not access " + fileName + " file")
		}
		l.Log.Printf("License file '%s' found...\n", fileName)
		l.copyrightNotice = string(copyrightNotice)
		l.copyrightNotice = l.copyrightNotice[:len(l.copyrightNotice)-1]
		return err
	}

	return l.LoadConfig(dir+"./.", liccorFileName)
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
			case SuffixGO, SuffixC, SuffixCPP, SuffixCXX, SuffixH, SuffixHPP, SuffixJAVA, SuffixJS, SuffixTS, SuffixTSX:
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
	if err != nil || (l.ignore != nil && l.ignore.Match(path, false)) {
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
	if l.NoticeBeforeText != "" {
		l.Log.Printf("License before text set to '%s'\n", l.NoticeBeforeText)
		l.copyrightNotice = l.NoticeBeforeText + "\n" + l.copyrightNotice
	}
	if l.NoticeAfterText != "" {
		l.Log.Printf("License after text set to '%s'\n", l.NoticeAfterText)
		l.copyrightNotice = l.copyrightNotice + "\n" + l.NoticeAfterText
	}

	lics := make(map[string]string)
	clike := "/*\n * " + strings.Replace(string(l.copyrightNotice), "\n", "\n * ", -1) + "\n */\n"
	lics["c-like"] = strings.Replace(clike, "\n * \n", "\n *\n", -1)
	lics["go"] = func() string {
		golic := "/*\n   " + strings.Replace(string(l.copyrightNotice), "\n", "\n   ", -1) + "\n*/\n"
		golic = strings.Replace(golic, "\n   \n", "\n\n", -1)
		return golic
	}()

	allSuccess := true
	for _, source := range l.Source {
		files, err := l.FindSrcFiles(source)
		if err != nil {
			fmt.Errorf("Error encountered while searching "+source+", %v", err)
			continue
		}

		for i := 0; i < len(files); i++ {
			pt := strings.LastIndex(files[i], ".")
			lic := ""
			//determine how to format the license
			switch files[i][pt:] {
			case SuffixGO:
				lic = lics["go"]
			case SuffixC, SuffixCPP, SuffixCXX, SuffixH, SuffixHPP, SuffixJAVA, SuffixJS, SuffixTS, SuffixTSX:
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
	}
	if allSuccess {
		fmt.Println("All files up to date!")
	}
}
