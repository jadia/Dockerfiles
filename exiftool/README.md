# Exiftool Dockerfile

Incase the URL stops working, the tar.gz of exiftool is included in the directory.
Remove the wget part from the Dockerfile and build the image.

## Usage

```bash
docker run --rm -v $(pwd):/tmp ghcr.io/jadia/dockerfiles/exiftool:latest imageName.jpg
```

Add alias

```bash
alias exiftool='docker run --rm -v $(pwd):/tmp ghcr.io/jadia/dockerfiles/exiftool:latest'
```
