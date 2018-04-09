FROM alpine:latest

WORKDIR /opt/place-chain
COPY dist/place-chain .

RUN ./place-chain init --chain-id foo-chain

EXPOSE 46656 46657 80

CMD ["./place-chain", "start", "--full-node"] 