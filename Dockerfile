FROM golang AS build
LABEL autodelete="true"

WORKDIR /app

ADD . /app

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o x-ray-echo-server .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=build /app/x-ray-echo-server .

CMD ["./x-ray-echo-server"]

EXPOSE 2000/tcp
EXPOSE 2000/udp
