FROM golang:1.18.2-alpine AS build-env
WORKDIR /work
COPY . .
RUN go build -o app .

FROM alpine:latest
WORKDIR /work
EXPOSE 5000
COPY --from=build-env /work/app .
CMD ["./app"]