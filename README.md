# 2016-site #

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

## Using with your own Myradio

### Creating an API key

Next you need a api_key to allow the website to access myradio's show information,

Connect to the MyRadio database and run the following SQL

```sql
INSERT INTO myury.api_key (key_string, description) VALUES ('ARANDOMSTRINGOFCHARACTERS', '2016-site development api key');
INSERT INTO myury.api_key_auth (key_string, typeid) VALUES ('ARANDOMSTRINGOFCHARACTERS', (SELECT typeid FROM l_action WHERE phpconstant = 'AUTH_APISUDO'));
```

[please choose a better key than 'ARANDOMSTRINGOFCHARACTERS']

You might need add some other database columns to create shows

for example:
-   explict podcasts (to create shows)
-   selector (expected by 2016-site/can remove this from models/index.go 2016-site)

This will fix shows not loading on 2016-site when using the base myradio as,

2016-site uses parts of database that aren't made on myradio creation.

### Finishing steps

After completing all these setups:
You can setup a reverse proxy to "https://localhost:4443" or configure ssl for https connections,

And change 2016-site to use your myradio instance:

In Config.toml:

```
myradio_api = "https://{hostname}/api/v2"
```



## Editor Config
There is a handy editor config file included in this repo, most editors/IDE's have support for this either natively or through a plugin, see [here](http://editorconfig.org/#download).

