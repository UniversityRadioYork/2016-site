# 201623-site #

## How to install ##
Full instructions for installation are available in the USING.md file.
If you already know what you're doing, the instructions below should suffice.

Configuration for site is located in config.toml, and you'll need a
[MyRadio](https://github.com/UniversityRadioYork/MyRadio) API key with the
requisite permissions copied into a .myradio.key file in the same directory.

Then follow the Requirements below.

## Requirements ##
Requires [Go 1.6](https://golang.org/) to compile and run, along with `sassc` to
compile the SCSS files. You may use other SASS compilers if you wish, but be
prepared for unexpected results.

Alternatively, you can use Docker alone

## Running the site ##

### Without Docker ###
```bash
$ make run # Builds scss files, and runs the server
```

### With Docker :whale: ###
```bash
$ make build-docker-image #Builds the image, will only have to be re-run if you change the Dockerfile
$ make docker #Runs the image
```


## Editor Config
There is a handy editor config file included in this repo, most editors/IDE's have support for this either natively or through a plugin, see [here](http://editorconfig.org/#download).
