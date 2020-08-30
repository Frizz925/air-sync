#!/bin/bash
set -e

SCRIPT_DIR=$(realpath $(dirname $0))

echo "Bulding web frontend..."
cd $SCRIPT_DIR/web
[ -d out ] && rm -rf out
yarn build
echo "Web frontend built."

echo "Copying web frontend into public directory..."
cd $SCRIPT_DIR/server
[ -d public ] && rm -rf public 
cp -rv $SCRIPT_DIR/web/out public
echo "Public directory copied."

echo "Deploying service to Google App Engine..."
gcloud app deploy
echo "Service deployed."

echo "Finished."
