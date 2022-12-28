.PHONY: ${BIN}

BIN=statii

${BIN}::
	go build -o $@ main.go

clean:
	rm -f ${BIN}