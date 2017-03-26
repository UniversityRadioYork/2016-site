SASS_COMPILER ?= sassc

SASS_DIR    := sass/
SASS_SRC    := $(shell find sass/ -name '*.scss')
MAIN_FILE   := sass/main.scss
OUTPUT_FILE := public/css/main.scss.css

SOURCES := $(shell find . -name '*.go')
BINARY  := 2016-site

all: build build-sass

build: $(SOURCES)
	go build

run: all
	./2016-site

build-sass: $(SASS_SRC)
	$(SASS_COMPILER) $(MAIN_FILE) $(OUTPUT_FILE)

build-docker-image:
	docker build -t 2016-site .

docker:
	docker run -it --rm -p 3000:3000 -v $$GOPATH/src/github.com/UniversityRadioYork/2016-site:/go/src/github.com/UniversityRadioYork/2016-site 2016-site

clean:
	rm -f $(OUTPUT_FILE) $(BINARY)

.PHONY: build-sass run clean docker build-docker-image test

test:
	go test ./utils
