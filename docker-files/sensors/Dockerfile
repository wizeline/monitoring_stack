FROM golang:alpine as builder

RUN apk update && apk add git
COPY ../../exporter/src/sensors/ $GOPATH/exporter/src/sensors/
WORKDIR $GOPATH/exporter/src/sensors/
RUN go env -w GO111MODULE=off
RUN go get -d -v


RUN go build -o /go/bin/sensors


FROM alpine
EXPOSE 2111
COPY --from=builder /go/bin/sensors /bin/sensors
ENTRYPOINT ["/bin/sensors"]

# Environment Variables
# OWM_PORT
# OWM_API_KEY
# OWM_LOCATION
# OWM_DURATION