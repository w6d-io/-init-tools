ARG GOVERSION=1.13.9
FROM golang:$GOVERSION as builder
ARG GOVERSION=1.13.9
ARG VCS_REF
ARG BUILD_DATE
ARG VERSION
ARG PROJECT_NAMESPACE
ARG PROJECT_PATH
ARG PROJECT_NAME
ENV GO111MODULE="on" \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build/${PROJECT_PATH}
RUN mkdir -p /build/bin
COPY . .
RUN ./build.sh

FROM alpine
ARG VCS_REF
ARG BUILD_DATE
ARG VERSION
ARG PROJECT_NAMESPACE
ARG PROJECT_PATH
ARG PROJECT_URL
ARG PROJECT_NAME
ARG USER_EMAIL="david.alexandre@w6d.io"
ARG USER_NAME="David ALEXANDRE"
LABEL maintainer="${USER_NAME} <${USER_EMAIL}>" \
        org.label-schema.vcs-ref=$VCS_REF \
        org.label-schema.vcs-url=$PROJECT_URL \
        org.label-schema.build-date=$BUILD_DATE \
        org.label-schema.version=$VERSION


WORKDIR /opt/bin
COPY --from=builder /build/bin /opt/bin
