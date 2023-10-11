FROM golang:1.21

ENV SASS_VERSION=3.3.6

RUN set -e && \
    curl -L https://github.com/sass/libsass/archive/$SASS_VERSION.tar.gz | tar -xvz -C /usr/local && \
    curl -L https://github.com/sass/sassc/archive/$SASS_VERSION.tar.gz | tar -xvz -C /usr/local && \
    SASS_LIBSASS_PATH=/usr/local/libsass-$SASS_VERSION make BUILD=static -C /usr/local/sassc-$SASS_VERSION && \
    cp /usr/local/sassc-$SASS_VERSION/bin/sassc /usr/local/bin/sassc && \
    rm -rf /usr/local/sassc-$SASS_VERSION /usr/local/libsass-$SASS_VERSION

RUN mkdir -p /go/src/github.com/UniversityRadioYork/201623-site
WORKDIR /go/src/github.com/UniversityRadioYork/201623-site

COPY . /go/src/github.com/UniversityRadioYork/201623-site

EXPOSE 3000

RUN go get -d -v

ENTRYPOINT echo "\033[0;31mWARNING: \033[0mRunning with Docker will change \"localhost\" to \"0.0.0.0\" in your config. Remember to change it back!" && \
sed -i 's/localhost/0.0.0.0/g' *.toml && make run
