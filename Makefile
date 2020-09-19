build:
	mkdir -p ./bin
	go build -o ./bin/2048 -v ./cmd
	cp ./bin/2048 ~/bin/2048

clean:
	rm -rf ./bin/
