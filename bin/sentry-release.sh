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

# If version is empty, set a default
export VERSION=$(
  sentry-cli releases propose-version \
    || echo "${SENTRY_RELEASE_PREFIX:-local}-$(date +%Y%m%d)"
)

echo "Releasing Version: $VERSION..."

sentry-cli releases new "$VERSION" \
  && (sentry-cli releases set-commits --auto "$VERSION" || echo "Failed to set related commit") \
  && sentry-cli releases finalize "$VERSION"

sentry-cli releases deploys $VERSION new -e ${APP_ENV:-production}

