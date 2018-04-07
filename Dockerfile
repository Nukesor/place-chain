FROM golang:1.10.1

WORKDIR /opt/tendermint

COPY setup.sh /opt/tendermint
RUN ./setup.sh

COPY app /opt/tendermint
COPY cmd /opt/tendermint
COPY helper /opt/tendermint
COPY static /opt/tendermint
COPY types /opt/tendermint

COPY .tendermint.genesis.json /opt/tendermint

RUN tendermint init

CMD helper/start.sh