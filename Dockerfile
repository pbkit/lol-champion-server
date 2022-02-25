FROM denoland/deno:alpine AS POLLAPO
ARG GITHUB_TOKEN
RUN apk add git
RUN git clone https://github.com/pbkit/pbkit.git
RUN deno install -n pollapo -A --unstable pbkit/cli/pollapo/entrypoint.ts
COPY . /work
WORKDIR /work
RUN pollapo install -t ${GITHUB_TOKEN}

FROM golang:1.17-bullseye AS BUILD
RUN apt update && apt install -y protobuf-compiler
RUN go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
COPY . /work
WORKDIR /work
COPY --from=POLLAPO /work/.pollapo .pollapo
RUN make

FROM gcr.io/distroless/cc-debian11
COPY --from=BUILD /work/lol-champion-server /lol-champion-server
EXPOSE 8080/tcp 8081/tcp
ENTRYPOINT ["/lol-champion-server"]
