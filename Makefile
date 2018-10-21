prepare:
	go get github.com/gobuffalo/packr/packr
	go generate github.com/sachaos/go-life/preset

install: prepare
	go install
