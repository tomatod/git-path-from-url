# git-path-from-url
A Git plugin that finds the local path corresponding to a URL, e.g.:
"https://github.com/mike/example/blob/main/src/hoge/bar.go".

## Usage
### 1. Install
If you are Golang user, you can install it by the next command.
```sh
go get github.com/tomatod/git-path-from-url
```

Otherwise, you can download the binary from [the release page](https://github.com/tomatod/git-path-from-url/releases) and then place it in any directory included by your PATH.

### 2. Usage
You can use it like the next command as a git plugin on a git project directory.

```sh
# for example, you are on a directory /usr/michael/example/src/boo/
git path-from-url http://github.com/mike/example/blob/main/src/hoge/bar.go
```

## License
MIT License. Refer to [LICENSE.txt](LICENSE.txt) for details.
