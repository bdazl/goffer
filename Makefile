.PHONY: rsync _ptbend ptbend imgimport

rsync:
	rsync -av out/ ~/Dropbox/genart/goffer

ptbend:
	go run ./cmd/mkmov -backup -proj ptbend0 -fcount 990

imgimport:
	go run ./cmd/mkmov -backup -proj imgimport -fcount 160

diffeq:
	go run ./cmd/mkmov -backup -proj diffeq -fcount 300

fulkonstett:
	go run ./cmd/mkmov -proj fulkonstett -fcount 1 -preview

ptbend-sync: ptbend rsync
imgimport-sync: imgimport rsync
diffeq-sync: diffeq rsync
