# 2016-site #

## How to install ##

### Development ###

#### Requirements ####

Requires [Go 1.6](https://golang.org/) to compile and run, along with something
to compile the SCSS files.

```bash
$ sassc sass/main.scss public/css/main.scss.css # Or whatever sass program you prefer
$ go run *.go
```
