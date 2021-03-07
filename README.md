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

### Via Homebrew
```bash
$ brew tap Connehito/gdp
$ brew install gdp
```

### Via go get
```bash
$ go get -u github.com/Connehito/gdp
```

Add $GOPATH/bin to the PATH environment variable.
```bash
export PATH=$PATH:$GOPATH/bin
```

## Usage

### Deploy
Add the tag to local repository and push the tag to remote(origin) repository.

```bash
# specify tag
$ gdp deploy -t TAG

# dry-run
$ gdp deploy -t TAG -d

# force(skipped validation)
$ gdp deploy -t TAG -f

# set tag automatically
$ gdp deploy
```

### Publish
Create the release note in GitHub which based on the merge commits of the tag.

```bash
# specify tag
$ gdp publish -t TAG

# dry-run
$ gdp publish -t TAG -d

# force(skipped validation)
$ gdp publish -t TAG -f

# set tag automatically
$ gdp publish
```

## Specification

### Supported tag's format
- [semantic version](https://semver.org/): e.g. v1.2.3 or 1.2.3
- date version: e.g. 20180525.1 or release_20180525

### How to create generate note
Release note content is generated based on merge commit messages.

So, depending on your branch strategy, it may not be the intended result.

### What is last printed message?
When gdp succeeds, the following message is printed.

`Do not be satisfied with 'released', let's face user's feedback in sincerity!`

This is the watchword of [Connehito's developers](https://connehito.com/).

## Contributing
Bug reports and pull requests are welcome on GitHub at https://github.com/Connehito/gdp.

## License
gdp is available as open source under the terms of the [MIT License](http://opensource.org/licenses/MIT).
