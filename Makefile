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
	go run ./cmd/mkmov -proj fulkonstett -fcount 1

djanl-del:
	rm -f out/djanl.d/imgs/*

djanl: djanl-del
	go run ./cmd/mkmov -parallel -fps 5 -fcount 20 -proj djanl -w 2048 -h 2048

djanl-1: djanl-del
	go run ./cmd/mkmov -fps 24 -fcount 1 -proj djanl -w 2048 -h 2048

djanl-small: djanl-del
	go run ./cmd/mkmov -fps 4 -fcount 4 -proj djanl -w 2048 -h 2048

djanl-final: djanl-del
	go run ./cmd/mkmov -parallel -verbose -fps 24 -fcount 1500 -proj djanl -w 2048 -h 2048

ptbend-sync: ptbend rsync
imgimport-sync: imgimport rsync
diffeq-sync: diffeq rsync
