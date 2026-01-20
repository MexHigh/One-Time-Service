#!/bin/sh

CURRENT_TAG="$1"
CONFIG_VERSION=$(yq '.version' config.yaml)
if ! [ "$CONFIG_VERSION" = "$CURRENT_TAG" ]
then
  echo "Version in config.yaml does not match version $CURRENT_TAG"
  exit 1
fi

if ! grep -q "## $CURRENT_TAG" CHANGELOG.md
then
  echo "CHANGELOG.md does not have an entry for version $CURRENT_TAG"
  exit 1
fi

exit 0
