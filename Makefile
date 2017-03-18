SASS_COMPILER ?= sassc

SASS_DIR    := sass/
SOURCES     := $(shell find sass/ -name '*.scss')
MAIN_FILE   := sass/main.scss
OUTPUT_FILE := public/css/main.scss.css

all: build-sass

build-sass: $(SOURCES)
	$(SASS_COMPILER) $(MAIN_FILE) $(OUTPUT_FILE)
