# STEP 1 build executable binary
FROM golang:alpine as builder
RUN adduser -D -g '' appuser
RUN apk add --no-cache git
COPY . $GOPATH/src/MTConnect2MQTT/MTConnect2MQTT/
WORKDIR $GOPATH/src/MTConnect2MQTT/MTConnect2MQTT/
#get dependancies
RUN go get -d -v
#build the binary
RUN CGO_ENABLED=0 GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/MTConnect2MQTT

# STEP 2 build a small image
# start from scratch
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /go/bin/MTConnect2MQTT /go/bin/MTConnect2MQTT
USER appuser
ENTRYPOINT ["/go/bin/MTConnect2MQTT"]