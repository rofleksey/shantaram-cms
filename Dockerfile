FROM golang:1.25-alpine AS apiBuilder
WORKDIR /opt
RUN apk update && apk add --no-cache make
COPY . /opt/
RUN go mod download
ARG GIT_TAG=?
RUN make build GIT_TAG=${GIT_TAG}

FROM alpine
ENV ENVIRONMENT=production
WORKDIR /opt
RUN apk update && \
    apk add --no-cache curl ca-certificates && \
    update-ca-certificates && \
    ulimit -n 100000
COPY --from=apiBuilder /opt/shantaram /opt/shantaram
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/v1/healthz || exit 1
CMD [ "./shantaram" ]
