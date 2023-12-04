
FROM golang:1.21.4-alpine3.18

LABEL AUTHOR=garfieldkwong<garfieldkwong@gmail.com>


WORKDIR /go/src/app

RUN apk add curl wget gcc libc-dev

RUN apk add git

RUN apk add bash

ENV COOKIES=Please_insert_cookies

ADD run.sh .

#run go get github.com/muyouming/gphotosuploader
RUN go get -u github.com/jinzhu/gorm && go get gphotosuploader

ENV IGNORE_PATTERN='^.*/\.@__thumb/.*$'

CMD ["/bin/bash","run.sh"]
