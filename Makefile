

default: clean build

clean:
	bash scripts/clean/clean.sh


build:
	bash scripts/build/build-darwin-arm64.sh
	bash scripts/build/build-darwin-amd64.sh
	bash scripts/build/build-linux-amd64.sh
	bash scripts/build/build-windows-amd64.sh
	bash scripts/build/build-windows-arm64.sh