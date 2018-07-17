# snippetlib

[![Travis CI](https://img.shields.io/travis/jessfraz/snippetlib.svg?style=for-the-badge)](https://travis-ci.org/jessfraz/snippetlib)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/jessfraz/snippetlib)

Server to host code snippets.

 * [Installation](README.md#installation)
      * [Binaries](README.md#binaries)
      * [Via Go](README.md#via-go)
 * [Usage](README.md#usage)

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/jessfraz/snippetlib/releases).

#### Via Go

```console
$ go get github.com/jessfraz/snippetlib
```

## Usage

```
$ snippetlib --help
           _                  _   _ _ _
 ___ _ __ (_)_ __  _ __   ___| |_| (_) |__
/ __| '_ \| | '_ \| '_ \ / _ \ __| | | '_ \
\__ \ | | | | |_) | |_) |  __/ |_| | | |_) |
|___/_| |_|_| .__/| .__/ \___|\__|_|_|_.__/
            |_|   |_|

 Server to host code snippets.
 Version: v0.2.2
 Build: f5f7038

  -cert string
        path to ssl certificate
  -d    run in debug mode
  -dbconn string
        database connection string (default "postgres://postgres:@127.0.0.1:5432/db?sslmode=disable")
  -key string
        path to ssl key
  -mailchimp-apikey string
        Mailchimp APIKey for subscribing to newsletters
  -mailchimp-listid string
        Mailchimp List ID for newsletter to subscribe emails to
  -p string
        port for server to run on (default "3000")
  -v    print version and exit (shorthand)
  -version
        print version and exit
```
