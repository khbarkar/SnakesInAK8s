#!/usr/bin/env bash
set -euo pipefail

COUNT="${1:-25}"
NAMESPACE="snakefood"

ADJECTIVES=(
  blazing frozen cranky dizzy fluffy grumpy jazzy lucky mystic nerdy
  perky quirky rusty salty sneaky spicy tangy wacky witty zesty
  bouncy crispy dusty foggy giddy hasty jolly lumpy murky nippy
  plucky raspy shiny testy vivid wobbly zippy breezy chunky dapper
  earthy feisty gloomy humble itchy jumpy keen lanky mellow nutty
  ornery peppy queasy rowdy scruffy thorny uppity vexed whimsy yappy
  absurd bonkers clumsy dainty elastic funky groovy hefty icy jittery
)

NOUNS=(
  badger cactus dingo falcon gopher hedgehog iguana jackal koala lemur
  moose narwhal otter panda quokka raccoon squid toucan urchin vulture
  walrus yak zebra alpaca bison cobra donkey eagle ferret gecko
  heron ibis jaguar kiwi lobster marmot newt osprey parrot quail
  raven sloth tapir umbrellabird viper wombat xerus yapok zebrafish
  anchovy beetle cricket dragonfly emu flamingo goose hamster impala jellyfish
)

kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -

echo "spawning $COUNT snakefood pods in namespace $NAMESPACE..."

for i in $(seq 1 "$COUNT"); do
  ADJ=${ADJECTIVES[$((RANDOM % ${#ADJECTIVES[@]}))]}
  NOUN=${NOUNS[$((RANDOM % ${#NOUNS[@]}))]}
  NAME="${ADJ}-${NOUN}-$(printf '%03d' $i)"

  kubectl run "$NAME" \
    --namespace="$NAMESPACE" \
    --image=registry.k8s.io/pause:3.9 \
    --labels="app=snakefood" \
    --restart=Never \
    --overrides='{
      "spec": {
        "terminationGracePeriodSeconds": 0,
        "containers": [{
          "name": "morsel",
          "image": "registry.k8s.io/pause:3.9"
        }]
      }
    }' 2>/dev/null || true

  echo "  spawned $NAME"
done

echo "done. $COUNT pods ready to be devoured."
