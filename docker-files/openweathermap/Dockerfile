FROM golang:alpine as builder

RUN apk update && apk add git
COPY ../../exporter/src/openweathermap $GOPATH/exporter/src/openweathermap/
WORKDIR $GOPATH/exporter/src/openweathermap/
RUN go env -w GO111MODULE=off
RUN go get -d -v


RUN go build -o /go/bin/openweathermap


FROM alpine
EXPOSE 2112
COPY --from=builder /go/bin/openweathermap /bin/openweathermap
ENTRYPOINT ["/bin/openweathermap"]

# Environment Variables
# OWM_PORT
# OWM_API_KEY
# OWM_LOCATION
# OWM_DURATION