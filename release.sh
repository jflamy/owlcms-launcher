#!/bin/bash -x
export TAG=v1.9.1-alpha01
export DEB_TAG=${TAG#v}
git pull
dist/updateRc.sh ${DEB_TAG}
git commit -am "owlcms-launcher $TAG"
git push
git tag -a ${TAG} -m "owlcms-launcher $TAG"
git push origin --tags