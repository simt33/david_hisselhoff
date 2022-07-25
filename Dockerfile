FROM golang:1.13 AS build

WORKDIR /build

COPY . .
RUN GOOS=linux GOARCH=amd64 go build -o /build/run
RUN chmod +x /build/run

FROM scratch

COPY --from=build /build/run /

EXPOSE 8080
CMD ["/run"]
