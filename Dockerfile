FROM golang:1.17

WORKDIR /workspace/glcs
RUN git clone https://github.com/lnikon/glcs .

RUN go mod tidy
RUN go build ./cmd/glcs


ENV UPCXX_INSTALL="/shared-workspace/libs/upcxx"
ENV PGASGRAPH_INSTALL="/shared-workspace/pgasgraph"
ENV PATH="${UPCXX_INSTALL}/bin:${PGASGRAPH_INSTALL}/build/src/PGASGraphCLI:$PATH"

ARG HOST=
ENV HOST $HOST

ARG PORT=8080
ENV PORT $PORT

EXPOSE $PORT

CMD cp -rf /workspace/* /shared-workspace/ && cd /shared-workspace/glcs && ./glcs --http=$HOST:$PORT
