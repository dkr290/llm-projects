FROM golang:1.24.2-alpine3.20 AS builder
ENV CGO_ENABLED=0
COPY . /build

#Build the service binary.
WORKDIR /build
RUN go build -o press-detective -ldflags "-X main.build=main" ./cmd/server/

#Multistage build

FROM alpine:3.21.2
#Add maintainer info
LABEL maintainer="Danail Surudzhiyski && Bernard Brenyah"

RUN addgroup -g 1000 -S pressdetective && \
  adduser -u 1000 -h /cmd -G pressdetective -S pressdetective

#Copy from stage 0 builder only the binary files
COPY --from=builder --chown=pressdetective:pressdetective /build/press-detective .
# Create the /doc directory and set ownership
RUN mkdir /doc && \
  chown -R pressdetective:pressdetective /doc

WORKDIR /

USER pressdetective
CMD "./press-detective"
