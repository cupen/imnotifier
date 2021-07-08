FROM golang:1.16.3-buster AS build
LABEL maintainer="xcupen@gmail.com"

ARG goproxy

ENV GOPROXY=${goproxy}

COPY ./  /notifier/
RUN cd /notifier/ && go build -o server main.go

FROM debian:10.9-slim AS runtime
COPY --from=build /notifier/ /notifier/
WORKDIR /notifier/
ENTRYPOINT ["/notifier/server"]
CMD ["--listen", "127.0.0.1:7788"]
