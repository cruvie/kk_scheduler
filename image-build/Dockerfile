##
## Build
##
ARG GOVERSION
FROM golang:${GOVERSION} AS build


ENV GOPROXY=https://goproxy.cn,direct

COPY . /kk-scheduler

WORKDIR /kk-scheduler

RUN CGO_ENABLED=0 GOOS=linux go build -o /kk-scheduler/output /kk-scheduler/go/internal/main

##
## Deploy
##
FROM alpine:latest



COPY --from=build /kk-scheduler/output /kk-scheduler/output
COPY ui/.output/public /kk-scheduler/public

WORKDIR /kk-scheduler

CMD ["./output"]

