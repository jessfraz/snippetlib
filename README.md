# snippetlib

[![Circle CI](https://circleci.com/gh/jfrazelle/snippetlib.svg?style=svg)](https://circleci.com/gh/jfrazelle/snippetlib)

Server to host code snippets.

```
$ snippetlib --help
Usage of snippetlib:
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
```
