# docker/mockllm.Dockerfile (ticket 08, SPEC-11 §4): the tools/mockllm mock
# LLM server as a container for compose.test.yml. mockllm is a stdlib-only
# separate Go module, so the image builds without the sherpa/sqlite cgo
# budget. Service port: 18080.
FROM golang:1.27-alpine AS build
WORKDIR /src
COPY tools/mockllm/go.mod ./
RUN go mod download || true
COPY tools/mockllm/ ./
RUN CGO_ENABLED=0 go build -trimpath -o /out/mockllm .

FROM alpine:3.20
COPY --from=build /out/mockllm /usr/local/bin/mockllm
EXPOSE 18080
ENTRYPOINT ["mockllm", "-addr", ":18080"]
