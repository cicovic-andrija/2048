build:
	mkdir -p ./bin
	go build -o ./bin/term2048 -v ./cmd
	cp ./bin/term2048 ~/bin/2048

clean:
	rm -rf ./bin/
