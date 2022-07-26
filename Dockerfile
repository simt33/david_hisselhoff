FROM golang:1.13 AS build

WORKDIR /build

COPY . .
RUN CGO_ENABLED=0 go build -o /build/run
RUN chmod +x /build/run

FROM scratch

COPY --from=build /build/run /

EXPOSE 8080
CMD ["/run"]
