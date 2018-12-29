# snippetlib

[![Travis CI](https://img.shields.io/travis/jessfraz/snippetlib.svg?style=for-the-badge)](https://travis-ci.org/jessfraz/snippetlib)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/jessfraz/snippetlib)

Server to host code snippets.

**Table of Contents**

<!-- toc -->

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
