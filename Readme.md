# liccor [![Build Status](https://travis-ci.org/gtalent/liccor.svg?branch=master)](https://travis-ci.org/gtalent/liccor)
A license notice corrector for C/C++, Java, JavaScript, and Go.


## Prerequisite
Install Go: http://golang.org/doc/install


## Installation
Run the following command:

    go get github.com/gtalent/liccor


## Usage
Place a ```.liccor``` file in the root directory of the project, containing your license notice. The license notice in the file should not be commented out, liccor takes care of that.

Run liccor to apply the copyright notice from ```.liccor``` to all the source files of the current directory and all sub-directories.

Check out the ```--help``` for more information.

## .liccor
The ```.liccor``` file only contains the desired copyright notice.

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

## .liccor.yml
The ```.liccor.yml``` file contains the desired copyright notice, which files
to ignore, and which directories to target for updating. The ignore section
uses the ```.gitignore``` syntax.

```yaml
---
source:
- src
ignore: |-
  deps/*
copyright_notice: |-
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
```

## Contributors
- [Gary Talent](https://github.com/gtalent)
- [Paul Vollmer](https://github.com/paulvollmer)
