FROM golang:1.17

WORKDIR /workspace/glcs
RUN git clone https://github.com/lnikon/glcs .

RUN go mod tidy
RUN go build ./cmd/glcs

EXPOSE 8090

ENV DB_NAME="glcs"
ENV DB_USER="postgres"
ENV DB_PASSWORD="postgres"
ENV UPCXX_INSTALL="/workspace/libs/upcxx"
ENV PGASGRAPH_INSTALL="/workspace/pgasgraph"
ENV PATH="${UPCXX_INSTALL}/bin:$PATH"

CMD cp -rf /workspace/* /shared-workspace && cd /shared-workspace/glcs && ./glcs --http=:8090
