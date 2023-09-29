all:
	go env GOCACHE
	go build -o goplug-minimal -gcflags 'all=-N -l' -mod=readonly
	go build -o goplug-minimal.so -gcflags 'all=-N -l' -mod=readonly -buildmode=plugin
clean:
	@rm -f goplug-minimal goplug-minimal.so *.json *.log
