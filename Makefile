
.PHONY: rsync
rsync:
	rsync -av out/ ~/Dropbox/genart/goffer

.PHONY: ptbend
_ptbend:
	go run . -proj ptbend0 -fcount 960

.PHONY: ptbend
ptbend: _ptbend rsync
