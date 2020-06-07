# Exiftool Dockerfile

```bash
docker run --rm -v $(pwd):/tmp exiftool imageName.jpg
```

Add alias

```bash
alias exiftool='docker run --rm -v $(pwd):/tmp exiftool'
```
