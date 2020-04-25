FROM golang:alpine as builder
COPY ./ /home
RUN go env -w GO111MODULE="on" && cd /home && \
    go build -ldflags "-s" main.go

FROM alpine:latest
COPY  --from=builder /home/main /home/sonar-gerrit-review
RUN chmod +x /home/sonar-gerrit-review
ENTRYPOINT ["/home/sonar-gerrit-review"]
