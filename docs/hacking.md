# Hacking Antroid

## How to read the code

To understand the code in `api/`, start with `io.go`, which describes the
low-level HTTPS client. It uses structs described in `params.go` for the
parameters and in `responses.go` for the responses. Then read `client.go`, the
API client based on the previous one. All the errors are described in
`errors.go`. Games structs are described in `game.go` and their specs (i.e.
their rules) are in `game_spec.go`. Turns are described in `turns.go`. The API
info structs are in `api_info.go`.

You’re done with the API part. Now let’s see the game server. Its code is in
`server.go`. It uses AIs, described in `ai.go` and plugins described in
`plugins.go`. Both of them are wrappers around actors, in `actors.go`. The game
server maintain a partial map between turns, which you can find in `maps.go`.

Some pretty-printing facilities are in `pretty_printing.go`, and that’s it.

## How to read the doc

If you’ve correctly set up your local environment you should be able to start a
Go documentation server with:

    godoc -http=:6060

Then open your browser at <http://localhost:6060> and browse to
`github.com/bfontaine/antroid/api`. You’ll get all the documentation of all
exported types and functions from the API on one page.

## How to add an AI

Refer to `ai_protocol.md` for the protocol. Write a program in the language of
your choice that reads on its standard input turn infos and write commands on
its standard output. The game server is intentionally agnostic regarding which
AI controls which ants. You could accept the ids of the ants you control as an
argument, for example.

In pseudo-code, an AI does the following:

    while not done:
        read turn infos on stdin
        (process them)
        write a command on one line on stdout

It can do anything if it respects this protocol.

## How to add a GUI

GUIs are exactly like AIs except they don’t produce any output on stdout (or at
least we don’t listen to it).

## How to add an API method

Let’s say YRG decides to add a new API method, `/unplay`, that would allow you
to “unplay” a turn, i.e. undo your last turn. It would take a game ID and
return a success status or an error, using the `POST` verb.

First, go in `io.go` and add an API call that takes a `GameIDParams` struct and
perform a `POST /unplay`:

    // CallUnplay performs a call to /unplay.
    func (h *Httclient) CallUnplay(params GameIDParams) *Body {
        return h.post("/unplay", params)
    }

Then add the corresponding method in `client.go`. This one is easy because the
remote server is not expected to send back a structure.


    // UnplayGameIdentifier undoes the last turn on this game, given its
    // identifier.
    func (cl *Client) UnplayGameIdentifier(id GameID) error {
        body := cl.http.CallUnplay(GameIDParams{ID: id})

        return body.ensureEmptyResponse()
    }

Add a method that takes a game instead of its identifier. This one is to be
friendly with developpers.

    // UnplayGame undoes the last turn on this game
    func (cl *Client) UnplayGame(g *Game) error {
        return cl.UnplayGameIdentifier(g.Identifier)
    }

That’s it for the API. Now we want to have a subcommand on `./antroid` to do
that. Edit `antroid.go`, and add the subcommand in the subcommands list:

    unplayCmd = app.Command("unplay", "Unplay the last turn.")

Add a game ID argument:

    unplayID = unplayCmd.Arg("id", "game ID").Required().String()

Now in the `main` function, search for the `switch` that performs all actions.
Add a `case`:

    case unplayCmd.FullCommand():
        // create a GameID from a string
        gID := api.GameID(*unplayID)

        // call the client
        if err := cl.UnplayGameIdentifier(gID); err != nil {
            // if there was an error, log it and exit
            exitErr(err)
        } else {
            // success
            fmt.Println("Last turn on game %s successfully undone\n", gID)
        }

Now re-compile, and you should be able to use it: `./antroid unplay <game id>`.
