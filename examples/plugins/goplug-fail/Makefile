all:
	go env GOCACHE
	go build -o goplug-fail -gcflags 'all=-N -l' -mod=readonly
	go build -o goplug-fail.so -gcflags 'all=-N -l' -mod=readonly -buildmode=plugin
clean:
	@rm -f goplug-fail goplug-fail.so *.json *.log
