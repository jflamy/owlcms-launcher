#!/bin/bash -x
export TAG=v1.9.4-alpha03
git tag -d ${TAG}
git push origin --delete ${TAG}
gh release delete ${TAG} --repo owlcms/owlcms-controlpanel --yes

BUILD_MAC=false
BUILD_WINDOWS=true
BUILD_RASPBERRY=false
BUILD_LINUX=false

# Pull the latest changes
git pull

# Update the resource configuration
export DEB_TAG=${TAG#v}
dist/updateRc.sh ${DEB_TAG}

# Substitute the values in release.yaml
sed -i "s/BUILD_MAC: .*/BUILD_MAC: ${BUILD_MAC}/" .github/workflows/release.yaml
sed -i "s/BUILD_WINDOWS: .*/BUILD_WINDOWS: ${BUILD_WINDOWS}/" .github/workflows/release.yaml
sed -i "s/BUILD_RASPBERRY: .*/BUILD_RASPBERRY: ${BUILD_RASPBERRY}/" .github/workflows/release.yaml
sed -i "s/BUILD_LINUX: .*/BUILD_LINUX: ${BUILD_LINUX}/" .github/workflows/release.yaml

# Commit and push the changes
git commit -am "owlcms-launcher $TAG"
git push
git tag -a ${TAG} -m "owlcms-launcher $TAG"
git push origin --tags