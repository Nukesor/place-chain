# PLACE CHAAAAAIN

## Initial install

##### $GOPATH and checkout

You *must* checkout this repository into `$GOPATH/src/place-chain`. Otherwise you are in eternal hell and nothing will ever work.

##### Install / Dev Env

You have to install required tools and dependencies to work with `place-chain`. Everything is wrapped with the Makefile. You must have `go` installed, `go version >=1.9`.

    $ make get-tools
    $ make install-deps

Or simply

    $ make all

## Build

Use the Makefile

    $ make build

After that, you have a folder called `dist/` with the `place-chain` binary in it. You can also directly install the binary via 

    $ make install

To create a lightweight docker container you can call

    $ make container

## Run

Use the commands, find them in `cmd/placechainnode/commands`.
Either you build first to get the binary or you can directly run the commands via 

    $ go run cmd/placechainnode/main.go <command>


*Setup new node*

    $ place-chain init --chain-id my-chain

This will generate a folder in your home directory called `place-chain`, where all needed files for the blockchain are stored.

*Run Full Tendermint Node*

    $ place-chain start --full-node

This will start three things:
- A fully fledged tendermint node
- A place-chain application
- A ABCI server

The ABCI server takes inbound connections from the TM node core and dispatches it to our custom place-chain app.

*Run ABCI Server only*

    $ place-chain start

This will only start a ABCI server together with our custom app. Any tendermint core process may then connect to the ABCI server with our app.

## Feature description

Inspired by r/place (last years reddit april prank).
We focus the application around a canvas of size `N x N`. Participants set pixels in a distributed fashion.

In order to set pixels to the canvas, users have to register. Registering includes:
- upload of a public key, use elliptic-curve `ed25519`
- upload a user profile (`{name, bio, avatarUrl}`)
User registration is placed onto the blockchain as a transaction.

Registered users can set pixels. To claim ownership of a pixel, a user signs the pixel coordinates and color with the private key, which then is verified with the public key that is stored in the blockchain.

When a pixel is set, this is done via a transaction of the blockchain. Painting the picture at some given point in time means traversing back the block chain until a value for a pixel is found.

##### Example Picture
![example image](img/place-chain_example_1.png)

##### Original idea(s):

The right to set a pixel on the canvas is bought with a crypto coin. The crypto coin we design in this project is called `place-coin`.

- earn coins when you colored them
- Fläche zusammenhängender gleichfarbiger Pixel umfärben (kostet coins aber Mengenrabatt)
- Hintergrundfarbe ändern (kostet extrem viele coins)
- Farbe wechseln kostet extra coins

Update: We will *not* pursue the `place-coin` idea. See the homework section instead.


### Participants:

- Arne Beer
- Rafael Epplee
- Hans Ole Hatzel
- Marcel Kamlot
- Felix Ortmann
- Benjamin Warnke


### Homework
- verteilter Betrieb
- signierte Teilnehmer
- sybil Attacke verhindern