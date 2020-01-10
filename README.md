# Download manga from https://mangahato.com

## Usage

```
go run ./src/go/main.go
```

## TODO

- [ ] Argument parser
- [ ] URL support

## 1. Retreive manga chapters

- Selector for chapter container

```
#tab-chapter table tbody tr
```

- Selector for images

```
center .chapter-content img
```

## URL mapping

- I decided to skip `href` url mapping because it is a pain.

URL format: `https://mangahato.com` + `/` + `[manga-document|chapter-document]`