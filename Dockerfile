FROM golang:1.23-alpine AS apiBuilder
WORKDIR /opt
COPY go.mod go.sum /opt/
RUN go mod download
COPY . /opt/
RUN go build -o ./shantaram-cms

FROM node:18-alpine
ENV ENVIRONMENT=production
WORKDIR /opt
RUN apk update && apk add --no-cache curl ca-certificates
COPY --from=apiBuilder /opt/shantaram-cms /opt/shantaram-cms
EXPOSE 8080
CMD [ "./shantaram-cms" ]