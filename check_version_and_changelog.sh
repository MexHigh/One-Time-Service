#!/bin/sh

CONFIG_VERSION=$(yq '.version' config.yaml)
if ! [ "$CONFIG_VERSION" = "$CI_COMMIT_TAG" ]
then
  echo "Version in config.yaml does not match version $CI_COMMIT_TAG"
  exit 1
fi

if ! grep -q "## $CI_COMMIT_TAG" CHANGELOG.md
then
  echo "CHANGELOG.md does not have an entry for version $CI_COMMIT_TAG"
  exit 1
fi

exit 0
