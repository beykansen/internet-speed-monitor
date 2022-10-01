# DO NOT USE bullseye for buildx see https://github.com/docker/buildx/issues/314
FROM golang:1.19.1-buster as builder
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y update && apt-get -y install build-essential pkg-config && apt-get clean

RUN mkdir -p /internet-speed-monitor
WORKDIR /internet-speed-monitor
ENV GO111MODULE on
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o internet-speed-monitor

###
FROM debian:buster-slim as final
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y update && apt-get -y install curl
RUN curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.deb.sh | bash
RUN DEBIAN_FRONTEND="noninteractive" apt-get -y install speedtest && apt-get clean

RUN mkdir -p /internet-speed-monitor
WORKDIR /internet-speed-monitor
COPY --from=builder /internet-speed-monitor .
RUN chmod +x internet-speed-monitor
CMD ["./internet-speed-monitor"]

