# go-life

Terminal based Conway's Game of Life. Implemented in Go.

![demo](https://user-images.githubusercontent.com/6121271/47264728-44ec2d80-d557-11e8-8994-d4af53126fe5.gif)

## Features

* Run on your terminal!
* Insert pattern from presets
* Color themes
* Mouse editing

## Install

### Binary

Go to release page and download.

```shell
$ wget https://github.com/sachaos/go-life/releases/download/v0.1.0/go-life_darwin_amd64 -O /usr/local/bin/go-life
$ chmod +x /usr/local/bin/go-life
```

### Manually Build

You need go version 1.11.

```shell
$ git clone https://github.com/sachaos/go-life.git
$ make install
```

## How to use

### Start

```shell
$ go-life
```

### Set pattern & theme

```shell
$ go-life --theme Ocean --pattern glider-gun
```

### Help

```shell
$ go-life --help
```

### Keymap

```
SPC: stop
Enter: step
c: clear
r: random
h: hide this message & status
p: switch preset
t: switch theme
LeftClick: switch state
RightClick: insert preset
```
