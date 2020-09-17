<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [snippetlib](#snippetlib)
  - [Installation](#installation)
      - [Binaries](#binaries)
      - [Via Go](#via-go)
  - [Usage](#usage)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

# snippetlib

[![make-all](https://github.com/jessfraz/snippetlib/workflows/make%20all/badge.svg)](https://github.com/jessfraz/snippetlib/actions?query=workflow%3A%22make+all%22)
[![make-image](https://github.com/jessfraz/snippetlib/workflows/make%20image/badge.svg)](https://github.com/jessfraz/snippetlib/actions?query=workflow%3A%22make+image%22)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/jessfraz/snippetlib)

Server to host code snippets.

**Table of Contents**

<!-- toc -->

- [Installation](#installation)
    + [Binaries](#binaries)
    + [Via Go](#via-go)
- [Usage](#usage)

<!-- tocstop -->

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/jessfraz/snippetlib/releases).

#### Via Go

```console
$ go get github.com/jessfraz/snippetlib
```

## Usage

```
$ snippetlib -h
snippetlib -  Server to host code snippets.

Usage: snippetlib <command>

Flags:

  --cert              path to ssl certificate (default: <none>)
  -d                  enable debug logging (default: false)
  --dbconn            database connection string (default: postgres://postgres:@127.0.0.1:5432/db?sslmode=disable)
  --key               path to ssl key (default: <none>)
  --mailchimp-apikey  Mailchimp APIKey for subscribing to newsletters (default: <none>)
  --mailchimp-listid  Mailchimp List ID for newsletter to subscribe emails to (default: <none>)
  -p                  port for server to run on (default: 3000)

Commands:

  version  Show the version information.
```
