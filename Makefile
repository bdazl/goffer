.PHONY: rsync _ptbend ptbend imgimport

rsync:
	rsync -av out/ ~/Dropbox/genart/goffer

_ptbend:
	go run ./cmd/mkmov -backup -proj ptbend0 -fcount 990

_imgimport:
	go run ./cmd/mkmov -backup -proj imgimport -fcount 160

_diffeq:
	go run ./cmd/mkmov -backup -proj diffeq -fcount 300

ptbend: _ptbend rsync
imgimport: _imgimport rsync
diffeq: _diffeq rsync
