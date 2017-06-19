#!/usr/bin/env bash
# Usage: misc/release.sh
# Build package, tag a commit, push it to origin, and then deploy the
# package on production server.

set -e

echo "Building..."
make build

version="$(./mdb-links version | awk '{print $NF}')"
[ -n "$version" ] || exit 1
echo $version

echo "Tagging commit and pushing to remote repo"
git commit --allow-empty -a -m "Release $version"
git tag "v$version"
git push origin master
git push origin "v$version"

echo "Uploading executable to server"
scp mdb-links archive@app.mdb.bbdomain.org:/sites/mdb-links/"mdb-links-$version"
ssh archive@app.mdb.bbdomain.org "ln -sf /sites/mdb-links/mdb-links-$version /sites/mdb-links/mdb-links"

echo "Restarting application"
ssh archive@app.mdb.bbdomain.org "supervisorctl restart mdb-links"