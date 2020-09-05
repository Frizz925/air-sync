#!/bin/bash
set -e

SCRIPT_DIR=$(realpath $(dirname $0))

bash $SCRIPT_DIR/build.sh

echo "Deploying service to Google App Engine..."
cd $SCRIPT_DIR/server
gcloud app deploy -q
gcloud app deploy -q cron.yaml
echo "Service deployed."
