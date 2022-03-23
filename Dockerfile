# STEP 1 build executable binary
FROM golang:alpine as builder
RUN adduser -D -g '' appuser
RUN apk add --no-cache git
COPY . $GOPATH/src/github.com/ArtemVladimirov/goMTConnect2MQTT
WORKDIR $GOPATH/src/github.com/ArtemVladimirov/goMTConnect2MQTT
#get dependancies
RUN go install
#build the binary
RUN CGO_ENABLED=0 GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/goMTConnect2MQTT

# STEP 2 build a small image
# start from scratch
FROM scratch
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable
COPY --from=builder /go/bin/goMTConnect2MQTT /go/bin/goMTConnect2MQTT
USER appuser
ENTRYPOINT ["/go/bin/goMTConnect2MQTT"]