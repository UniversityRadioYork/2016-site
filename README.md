# 2016-site #

## How to install ##
Configuration for site is located in config.toml, and you'll need a
[MyRadio](https://github.com/UniversityRadioYork/MyRadio) API key with the
requisite permissions copied into a .myradio.key file in the same directory.

Then follow the Requirements below.

### Install with Docker :whale:###
Alternatively, you only need [Docker](https://www.docker.com/)

```bash
$ docker build -t 2016-site .
$ docker run -p 3000:3000 -v $GOPATH/src/github.com/UniversityRadioYork/2016-site:/go/src/github.com/UniversityRadioYork/2016-site 2016-site
```

## Requirements ##
Requires [Go 1.6](https://golang.org/) to compile and run, along with `sassc` to
compile the SCSS files. You may use other SASS compilers if you wish, but be
prepared for unexpected results.

```bash
$ make run # Builds scss files, and runs the server
```

## Development ##
