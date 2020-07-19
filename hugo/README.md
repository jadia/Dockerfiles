# Hugo Docker Container

Original repo: [https://github.com/jguyomard/docker-hugo](https://github.com/jguyomard/docker-hugo)

## Alias

```bash
alias hugo='docker run --rm -it -v $PWD:/src -u hugo jguyomard/hugo-builder hugo'
alias hugo-server='docker run --rm -it -v $PWD:/src -p 1313:1313 -u hugo jguyomard/hugo-builder hugo server --bind 0.0.0.0'
```


## Version

Hugo: 0.55
Alpine: 3.12
