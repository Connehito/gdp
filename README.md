# gdp
gdp is a CLI tool for pushing the tag associated with deployment and publishing the release note in GitHub.

[![Build Status](https://travis-ci.org/Connehito/gdp.svg?branch=master)](https://travis-ci.org/Connehito/gdp)
[![Go Report Card](https://goreportcard.com/badge/github.com/Connehito/gdp)](https://goreportcard.com/report/github.com/Connehito/gdp)
[![GoDoc](https://godoc.org/github.com/Connehito/gdp?status.svg)](https://godoc.org/github.com/Connehito/gdp)
[![Latest Version](http://img.shields.io/github/release/Connehito/gdp.svg?style=flat-square)](https://github.com/Connehito/gdp/releases/latest)
[![License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://raw.githubusercontent.com/Connehito/gdp/master/LICENSE)

## Requirements
- [git command](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [hub command](https://github.com/github/hub#installation)

## Installation
```
$ go get -u github.com/Connehito/gdp
```

## Usage

### deploy
- Add the tag to local repository and push the tag to remote(origin) repository.

```
# specify tag
$ gdp deploy -t TAG

# dry-run
$ gdp deploy -t TAG -dry-run

# set tag automatically
$ gdp deploy
```

### publish
- Create the release note in GitHub which based on the merge commits of the tag.

```
# specify tag
$ gdp publish -t TAG

# dry-run
$ gdp publish -t TAG -dry-run

# set tag automatically
$ gdp publish
```

### supported tag's format
- [semantic version](https://semver.org/): e.g. v1.2.3 or 1.2.3
- date version: e.g. 20180525.1

## Contributing
Bug reports and pull requests are welcome on GitHub at https://github.com/Connehito/gdp.

## License
gdp is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
