# gogr

A tool for running commands in multiple directories. There are a lot like it
and this one is mine.

## Installation

```
$ go get github.com/kopoli/gogr
```

## Description

The idea was to have a quick-and-dirty Go implementation of the
![grtool](http://mixu.net/gr/).  The distinguishing feature is that this
requires no runtime dependencies and is able to run the commands concurrently.

It both supports the @tag idea and giving directories directly to the program. 

## Usage

```
# Add directories to tag @this

$ gogr tag add this . .. ../..

# run command in each directory

$ gogr @this ls -l

```

## License

MIT license
