#! /usr/bin/env sh

if [ -z $SENTRY_AUTH_TOKEN ]; then
  echo "SENTRY_AUTH_TOKEN is not set. Skipping release..."
  exit 0
fi

sentry_cmd="$(which sentry-cli)"
if [ -z $sentry_cmd ]; then
  curl -sL https://sentry.io/get-cli/ | bash
fi

SENTRY_ORG="${SENTRY_ORG:-gvempire}"
SENTRY_PROJECT="${SENTRY_PROJECT:-"go-whyye"}"
VERSION="$(
  sentry_cmd releases propose-version |
  xargs
)"

# Get a reference to the command sentry-cli using `which`

sentry_cmd releases new "$VERSION" \
  && sentry_cmd releases set-commits "$VERSION" --auto \
  && sentry_cmd releases finalize "$VERSION"
