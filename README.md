# Dockerfiles

This repository contains the Dockerfiles for variety of applications.
The GitHub actions workflow builds and pushes the image to Github Container Registry.
But Github Actions are only triggered when there are changes in the VERSION file on any application.

Most of the workflow is inspired from this blog: [https://www.learncloudnative.com/blog/2020-02-20-github-action-build-push-docker-images/](https://www.learncloudnative.com/blog/2020-02-20-github-action-build-push-docker-images/)


The images can be found here: [https://github.com/jadia?tab=packages&q=dockerfiles](https://github.com/jadia?tab=packages&q=dockerfiles)

## Add new application

1. Create a new folder with the image/application name.
2. Create `Dockerfile` and `VERSION` file in the folder and push the changes.
3. When `Dockerfile` or the `VERSION` file is updated or added, Github Actions will build and push the image to Github container registry.
