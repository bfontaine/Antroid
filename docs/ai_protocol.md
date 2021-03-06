# AI Messages Protocol

The game server communicates with AI programs using standard I/O. Each AI
program expects to receive messages on STDIN and send back messages on STDOUT.

## Format

### Server Message

The first line is as follow:

    T A P S

With `T` the turn number, `A` the ants count per player, `P` the players count,
and `S` the game status (`1` for `playing`, and `0` for `over`).

Note: each AI can control multiple ants, the protocol doesn’t care about that.

It then contains `A` lines, one per ant, with the following content:

    ID X Y DX DY E A B

With `ID` the ant’s ID, `X`, `Y` the position, `DX`, `DY` the direction, `E`
the energy level, `A` the acid level, and `B` the brain state (`1`: controlled,
we don't know the other values for now).

It then contains one line with the number of other players’ ants we’re able to
see at this turn:

    N

And `N` lines:

    X Y DX DY B

These are the same as the lines for our ants except that we don't have their
energy and acid levels nor their ID.

It then contains a line describing the map:

    W H N

With `W` being the known width (the maximum `X` plus one), `H` the known height
(the maximum `Y` plus one), and `N` the number of points.

The message then contains `N` lines, one per point:

    X Y C S

Where `X`, `Y` is the position and `C` is the content, as an int code as
described below. `S` is `1` if the point was seen this turn or `0` if it’s a
point we remember from a previous turn.

Contents:

* `0` (`000`) : grass
* `2` (`010`) : rock
* `4` (`100`) : water
* `1` (`001`) : food (sugar)
* `3` (`011`) : food (mill)
* `5` (`101`) : food (meat)

#### Example

    4 5 2 1
    0 0 0 -1 0 100 100 1
    1 25 78 -1 0 17 100 1
    2 112 42 -1 0 32 100 1
    3 5 3 -1 0 100 56 1
    4 2 10 -1 0 100 100 1
    0
    120 80 124
    0 0 0 1
    0 1 2 1
    4 6 0 1
    12 17 0 0
    3 5 0 1
    6 9 5 0
    ...

This describes a game at the turn `4`, with `5` ants per player, `2` players,
with a `playing` game, and this AI controls `5` ants. Their ids are from `0` to
`4`. The first ant is in (`0`, `0`) and has `100` energy and `100` acid, and
the known map is `120`x`80`.

### Client Message

The client is expected to send one command for all their ants, on one line:

    C

`C` is the command. For example:

    1:forward,2:rest,4:rest,3:right


## Game

On each turn, the game server sends a message to each AI program, which is then
expected to respond with an action to perform with their controlled ants.
