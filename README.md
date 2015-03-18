# Antroid

![logo](docs/antroid_64.png)

## Members

* David Galichet: `galichet.david@yahoo.fr`
* Julien Sagot: `ju.sagot@gmail.com`
* Nicolas Cailloux: `caillouxnicolas@gmail.com`
* Baptiste Fontaine: `b@ptistefontaine.fr`

## Constraints

* Use Scheme somewhere in the project
* Use a lot of libraries

## Build

    make

## Dependencies

### Go

* [Ubuntu](https://github.com/golang/go/wiki/Ubuntu): `sudo apt-get install golang`
* OS X : `brew install go`

If you don’t have a `GOPATH`, create a directory somewhere where you'll put
your Go code, for example `~/Go`.

Add the following to your `.bashrc` (adjust with your settings):

    export GOPATH="$HOME/Go"
    export PATH="$GOPATH/bin:$PATH"

Now reload your shell and run:

    mkdir -p "$GOPATH/src/github.com/bfontaine/antroid"
    cd "$GOPATH/src/github.com/bfontaine/antroid"
    git clone git@github.com:bfontaine/Antroid.git .
