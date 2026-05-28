##
## Build
##
ARG GOVERSION
FROM golang:${GOVERSION} AS build


ENV GOPROXY=https://goproxy.cn,direct

COPY . /kk_scheduler

WORKDIR /kk_scheduler

RUN CGO_ENABLED=0 GOOS=linux go build -o /kk_scheduler/output /kk_scheduler/go/internal/main

##
## Deploy
##
FROM alpine:latest



COPY --from=build /kk_scheduler/output /kk_scheduler/output
COPY ui/.output/public /kk_scheduler/public

WORKDIR /kk_scheduler

CMD ["./output"]

