FROM golang:1.19 AS build
RUN mkdir testproject && chmod 777 -R ./testproject
COPY . /testproject
WORKDIR /testproject
RUN go test -v ./...
RUN go get -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o .


FROM alpine:latest AS final
RUN mkdir app
RUN chmod 777 -R ./app
COPY --from=build /testproject /app
WORKDIR /app
RUN apk add --no-cache tzdata
ENV TZ=Asia/Almaty
CMD ["./testProject"]
EXPOSE 3000