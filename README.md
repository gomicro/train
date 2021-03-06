# Train
[![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/gomicro/train/Build/master)](https://github.com/gomicro/train/actions?query=workflow%3ABuild)
[![Go Reportcard](https://goreportcard.com/badge/github.com/gomicro/train)](https://goreportcard.com/report/github.com/gomicro/train)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/gomicro/train)
[![License](https://img.shields.io/github/license/gomicro/train.svg)](https://github.com/gomicro/train/blob/master/LICENSE.md)
[![Release](https://img.shields.io/github/release/gomicro/train.svg)](https://github.com/gomicro/train/releases/latest)

Train is a command line tool for creating pull requests on a project that follows the [Gitlab Flow](https://docs.gitlab.com/ee/workflow/gitlab_flow.html) style of release management.

# Requirements
Golang version 1.12 or higher

# Installation

```
go get github.com/gomicro/train
```

# Usage

See the help text for descriptions of what is available

```
train -h
```

# Versioning
The tool will be versioned in accordance with [Semver 2.0.0](http://semver.org).  See the [releases](https://github.com/gomicro/train/releases) section for the latest version.  Until version 1.0.0 the tool is considered to be unstable.

It is always highly recommended to vendor the version you are using.

# License
See [LICENSE.md](./LICENSE.md) for more information.
