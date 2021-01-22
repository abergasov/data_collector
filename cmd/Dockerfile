FROM golang:1.15 AS build
ARG gitcommithash
RUN echo "Based on commit: $gitcommithash"
COPY . /go/src/data_collector
WORKDIR /go/src/data_collector
RUN make build_docker hash="$gitcommithash"

# 2 шаг - создание контейнера для запуска
FROM debian:buster-slim
WORKDIR /app
COPY --from=build /go/src/data_collector/bin/* /app/
COPY --from=build /go/src/data_collector/configs/* /app/configs/
COPY --from=build /go/src/data_collector/scripts/* /app/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
RUN chmod +x /app/*
CMD /app/collector