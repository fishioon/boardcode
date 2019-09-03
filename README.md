## board code

Exercise your board code.


## build
// only test on macos
```shell
brew install tesseract
make

```

## usage

``` shell
./boardcode &
curl --upload-file testdata/hello.png 127.0.0.1:9981/compile/image
```
