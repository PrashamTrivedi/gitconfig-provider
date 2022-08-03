FROM gitpod/workspace-full:2022-05-17-12-26-08

RUN go install github.com/goreleaser/goreleaser@latest
