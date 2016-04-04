# liccor
A license notice corrector for C/C++, Java, JavaScript, and Go.


## Prerequisite
Install Go: http://golang.org/doc/install


## Installation
Run the following command:

    go get github.com/gtalent/liccor


## Usage
Place a ```.liccor``` file in the root directory of the project, containing your license notice. The license notice in the file should not be commented out, liccor takes care of that.

Run liccor to apply the copyright notice from ```.liccor``` to all the source files of the current directory and all sub-directories.

If you want to use an other license file, you can set it with the ```-license``` flag.

Check out the ```-help``` for more information.


## Contributors
- [Gary Talent](https://github.com/gtalent)
- [Paul Vollmer](https://github.com/paulvollmer)
