.PHONY: rsync _ptbend ptbend imgimport

rsync:
	rsync -av out/ ~/Dropbox/genart/goffer

_ptbend:
	go run . -backup -proj ptbend0 -fcount 990

_imgimport:
	go run . -backup -proj imgimport -fcount 160

ptbend: _ptbend rsync
imgimport: _imgimport rsync
