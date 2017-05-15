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

package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/gtalent/liccor/lib"
)

func main() {
	app := cli.NewApp()
	app.Name = "liccor"
	app.Version = liccor.Version
	app.Author = "the liccor contributors"
	app.Email = "https://github.com/gtalent/liccor/graphs/contributors"
	app.Usage = "A license notice corrector for C/C++, Java, JavaScript, and Go."
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "liccorFile, l",
			Usage: "the name of the liccor file",
		},
		cli.StringFlag{
			Name:  "src",
			Usage: "directory containing the source files to update",
		},
		cli.StringFlag{
			Name:  "before, b",
			Value: "",
			Usage: "the string before the copyright notice",
		},
		cli.StringFlag{
			Name:  "after, a",
			Value: "",
			Usage: "the string after the copyright notice",
		},
		cli.BoolFlag{
			Name:  "verbose, V",
			Usage: "print verbose output",
		},
	}
	app.Action = func(c *cli.Context) error {
		liccor := liccor.New()
		liccor.Log.Active = c.Bool("verbose")
		liccor.LoadConfig(".", c.String("liccorFile"))
		if c.String("src") != "" {
			liccor.Source = []string{c.String("src")}
		}
		liccor.NoticeBeforeText = c.String("before")
		liccor.NoticeAfterText = c.String("after")
		liccor.Process()
		return nil
	}
	app.Run(os.Args)
}
