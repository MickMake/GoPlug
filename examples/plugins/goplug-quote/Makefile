all:
	go env GOCACHE
	go build -o goplug-quote -gcflags 'all=-N -l' -mod=readonly
	go build -o goplug-quote.so -gcflags 'all=-N -l' -mod=readonly -buildmode=plugin
clean:
	@rm -f goplug-quote goplug-quote.so
