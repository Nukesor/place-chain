FROM alpine:latest

WORKDIR /opt/place-chain
COPY dist/place-chain .

RUN ./place-chain init --chain-id foo-chain

EXPOSE 46656 46657 8080

CMD ["./place-chain", "start", "--full-node"] 