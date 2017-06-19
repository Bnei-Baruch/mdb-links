#!/usr/bin/env bash
# Usage: misc/deploy-internal.sh
# Deploy a released version to internal production server

set -e

version="$(./mdb-links version | awk '{print $NF}')"
[ -n "$version" ] || exit 1
echo $version

echo "Uploading executable to server"
scp mdb-links archive@app.mdb.bbdomain.org:/sites/mdb-links/"mdb-links-$version"
ssh archive@app.mdb.bbdomain.org "ln -sf /sites/mdb-links/mdb-links-$version /sites/mdb-links/mdb-links"

echo "Restarting application"
ssh archive@app.mdb.bbdomain.org "supervisorctl restart mdb-links"