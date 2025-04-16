FROM golang:1.24.2-alpine3.20 AS builder
ENV CGO_ENABLED=0
COPY . /build

#Build the service binary.
WORKDIR /build
RUN go build -o api -ldflags "-X main.build=main" ./cmd/api/

#Multistage build

FROM alpine:3.21.2

RUN addgroup -g 1000 -S app && \
  adduser -u 1000 -h /cmd -G app -S app

#Copy from stage 0 builder only the binary files
COPY --from=builder --chown=app:app /build/api .
# Create the /doc directory and set ownership
RUN mkdir /doc && \
  chown -R app:app /doc

WORKDIR /
EXPOSE 9999
USER app
CMD "./api"
