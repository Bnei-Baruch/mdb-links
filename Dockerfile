ARG work_dir=/go/src/github.com/Bnei-Baruch/mdb-links
ARG build_number=dev

FROM golang:1.17-alpine3.15 as build

LABEL maintainer="edoshor@gmail.com"

ARG work_dir
ARG build_number

ENV GOOS=linux \
	CGO_ENABLED=0

RUN apk update && \
    apk add --no-cache \
    git

WORKDIR ${work_dir}
COPY . .
RUN go build -ldflags "-w -X github.com/Bnei-Baruch/mdb-links/version.PreRelease=${build_number}"

FROM alpine:3.15

ARG work_dir
WORKDIR /app
COPY misc/*.sh ./
COPY --from=build ${work_dir}/mdb-links .

EXPOSE 8081
CMD ["./mdb-links", "server"]
