# 2016-site #

#### Requirements ####

If you wish to not use docker, you will need the following things installed:
 - [NodeJs](https://nodejs.org/en/ "NodeJs")
 - [Go 1.6](https://golang.org/ "Go")

Otherwise, you only need [Docker](https://www.docker.com/)

## How to install ##
# Install with Docker :whale:#
```bash
$ docker build -t 2016-site .
$ docker run -p 3000:3000 -v $GOPATH/src/github.com/UniversityRadioYork/2016-site:/go/src/github.com/UniversityRadioYork/2016-site 2016-site
```
# Install without Docker #
```bash
$ ./install.sh
$ go run main.go
```
### Development ###
