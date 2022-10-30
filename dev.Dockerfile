FROM golang:1.18.7-alpine3.16 as builder

RUN apk add --no-cache \
    build-base \
    libmediainfo-dev \
     openssl

WORKDIR /app

ENV GO111MODULE=on
ENV GOPATH /go
#RUN set -eux; \
#	export GOROOT="$(go env GOROOT)"; \
#	./run-tool-which-requires-GOROOT-set.sh \

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

#COPY wait-for-it.sh wait-for-it.sh
#RUN chmod +x wait-for-it.sh

COPY create-env.sh .
RUN chmod +x create-env.sh


RUN go install github.com/cosmtrek/air@latest


CMD ["air", "-c", ".air.toml"]

#FROM scratch
#
#WORKDIR /
#
#COPY --from=builder /app /app
#
#

##RUN set -eux; \
##	export GOROOT="$(go env GOROOT)"; \
##	./air \