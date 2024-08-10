#! /usr/bin/env sh

[ "$SKIP_RELEASE" = "1" ] && echo "Skipping release build" && exit 0

if [ -z $SENTRY_AUTH_TOKEN ]; then
  echo "SENTRY_AUTH_TOKEN is not set. Exiting..."
  exit 1
fi

if ! command -v sentry-cli > /dev/null; then
  curl -sL https://sentry.io/get-cli/ | bash
fi

export SENTRY_LOG_LEVEL=info
export SENTRY_ORG=${SENTRY_ORG:-gvempire}
export SENTRY_PROJECT=${SENTRY_PROJECT:-"go-whyye"}
export VERSION=${SENTRY_RELEASE:-$(
  sentry-cli releases propose-version |
  xargs
)}

echo "Version: $VERSION"

sentry-cli releases new "$VERSION" \
  && sentry-cli releases set-commits "$VERSION" --auto \
  && sentry-cli releases finalize "$VERSION"
