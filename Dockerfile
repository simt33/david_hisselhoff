FROM golang:1.13 AS build

WORKDIR /build

COPY . .
RUN go build

FROM golang:1.13

COPY --from=build /build/starter-snake-go /usr/local/bin/

EXPOSE 8080
CMD ["/usr/local/bin/starter-snake-go"]
