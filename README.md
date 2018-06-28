# gowebapp

Simple Lightning Payment Processor

Based on the Simplistic Go Web App [here](https://grisha.org/blog/2017/04/27/simplistic-go-web-app/).

## Installation

### With Docker:

First time or to rebuild everything and delete container data: `docker-first-run.sh`

Hit Ctrl-c to stop the container.

Subsequent times, assuming you want to persist data: `docker-run.sh`

Connect to database from another terminal with: `psql -hlocalhost -Upostgres`

### Non-dockerized instructions

To install:

1. You need PostgreSQL, install it somehow, make an empty database
   called `gowebapp`. (Usually `createdb gowebapp` on the command line is
   all you need).

1. You need Go, preferably the latest version (1.8).

1. Now do this:

```
$ export GOPATH=~/go # optional, adjust as necessary
$ go get github.com/weex/slpp
```

1. That's it. You should now be able to run:
```
$ $GOPATH/bin/gowebapp
```
