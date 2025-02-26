#!/bin/bash

set -euo pipefail

xfconf_query=$(which xfconf-query)

if ! which xfconf-query | grep -q "/xfconf-profile/"; then
  echo "ERROR: Using system's xfconf-query. Check your PATH."
  exit 1
fi

export XFCONF_PROFILE_END_TO_END_TEST=1
export LOG_LEVEL=debug

xfconf-profile apply profile.json

if diff expected-log.txt actual-log.txt; then
  echo "Differences detected"
  exit 1
else
  echo "Passed"
  exit 0
fi
