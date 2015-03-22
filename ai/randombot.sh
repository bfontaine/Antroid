#! /bin/bash
#
# randombot.sh
# ------------
#
# An "AI" that do random actions
#

turn() {
  local T A P S N W H
  local i l r

  read T A P S

  echo 1 >&2

  # for each ant
  if [ x"$A" != x0 ]; then
    for i in $(seq $A); do
      read
    done
  fi

  read N
  # for each other ant
  if [ x"$N" != x0 ]; then
    for i in $(seq $N); do
      read
    done
  fi

  # Map
  read W H N
  if [ x"$N" != x0 ]; then
    for i in $(seq $N); do
      read
    done
  fi

  # get a random number in (0,1,2,3)
  r=${RANDOM}
  let "r %= 4"

  # return the command
  case "$r" in
    0) echo "0:right" ;;
    1) echo "0:left" ;;
    2) echo "0:forward" ;;
    3) echo "0:rest" ;;
  esac
}

# play until it's the end
while true; do
  turn
done
