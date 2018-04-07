# PLACE CHAAAAAIN

### Setup:

        go get github.com/satori/go.uuid
        go get github.com/tendermint/abci/cmd/abci-cli
        go get -u github.com/tendermint/tendermint/cmd/tendermint

### How to run ?????
- Run tendermint using `tendermint node`
- Run our application: `go run cmd/placechainnode/main.go start`

### Steps:

Inspired by r/place (last years reddit april prank).

We focus the application around a canvas of size `N x N`. Participants set pixels in a distributed fashion.

The right to set a pixel on the canvas is bought with a crypto coin. The crypto coin we design in this project is called `place-coin`.

When a pixel is set, this is done via a transaction of the blockchain. Painting the picture means traversing back the block chain until a value for a pixel is found.

### Participants:

- Arne Beer
- Rafael Epplee
- Hans Ole Hatzel
- Marcel Kamlot
- Felix Ortmann
- Benjamin Warnke


### Kern-Funktionalität
- Pixel-Farbe wählen aus z.B. 8 vordefinierten Farben (kostet coins)
- jeder Benutzer 1 node (keine Account-Verwaltung)
- neue nodes starten mit x coins
- setzen von Pixeln (kostet coins)
- coins verdienen wenn Pixel sichtbar sind

### mögliche-Features
- Fläche zusammenhängender gleichfarbiger Pixel umfärben (kostet coins aber Mengenrabatt)
- Hintergrundfarbe ändern (kostet extrem viele coins)
- Farbe wechseln kostet extra coins
