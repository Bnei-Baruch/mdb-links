#!/usr/bin/env bash
# Usage: misc/release.sh
# Build package, tag a commit and push it to origin.

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