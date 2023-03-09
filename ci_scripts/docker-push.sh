#!/usr/bin/env bash

function print_usage() {
  echo "Use: $0 [-t <docker_image_tag_name>] [-p | -b]"
  echo "use -p for push (it builds and push the image)"
  echo "use -b for build image locally"
}

if [[ $# -ne 3 ]]; then
  print_usage
  exit 0
fi

function docker_build() {
  docker image build \
    --tag=skycoin/git-telegram-bot:"$tag" \
    -f ./Dockerfile .
}

function docker_push() {
  docker login -u "$DOCKER_USERNAME" -p "$DOCKER_PASSWORD"
  docker tag skycoin/git-telegram-bot:"$tag" skycoin/git-telegram-bot:"$tag"
  docker image push skycoin/git-telegram-bot:"$tag"
}

while getopts ":t:pb" o; do
  case "${o}" in
  t)
    tag="$(echo "${OPTARG}" | tr -d '[:space:]')"
    if [[ $tag == "develop" ]]; then
      tag="test"
    elif [[ $tag == "master" ]]; then
      tag="latest"
    fi
    ;;
  p)
    docker_build
    docker_push
    ;;
  b)
    docker_build
    ;;
  *)
    print_usage
    ;;
  esac
done
