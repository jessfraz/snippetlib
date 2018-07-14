# snippetlib

[![Travis CI](https://travis-ci.org/jessfraz/snippetlib.svg?branch=master)](https://travis-ci.org/jessfraz/snippetlib)

Server to host code snippets.

## Installation

#### Binaries

- **darwin** [386](https://github.com/jessfraz/snippetlib/releases/download/v0.2.1/snippetlib-darwin-386) / [amd64](https://github.com/jessfraz/snippetlib/releases/download/v0.2.1/snippetlib-darwin-amd64)
- **freebsd** [386](https://github.com/jessfraz/snippetlib/releases/download/v0.2.1/snippetlib-freebsd-386) / [amd64](https://github.com/jessfraz/snippetlib/releases/download/v0.2.1/snippetlib-freebsd-amd64)
- **linux** [386](https://github.com/jessfraz/snippetlib/releases/download/v0.2.1/snippetlib-linux-386) / [amd64](https://github.com/jessfraz/snippetlib/releases/download/v0.2.1/snippetlib-linux-amd64) / [arm](https://github.com/jessfraz/snippetlib/releases/download/v0.2.1/snippetlib-linux-arm) / [arm64](https://github.com/jessfraz/snippetlib/releases/download/v0.2.1/snippetlib-linux-arm64)
- **solaris** [amd64](https://github.com/jessfraz/snippetlib/releases/download/v0.2.1/snippetlib-solaris-amd64)
- **windows** [386](https://github.com/jessfraz/snippetlib/releases/download/v0.2.1/snippetlib-windows-386) / [amd64](https://github.com/jessfraz/snippetlib/releases/download/v0.2.1/snippetlib-windows-amd64)

#### Via Go

```bash
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
 Version: v0.2.1
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
