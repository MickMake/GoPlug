all:
	go env GOCACHE
	go build -o goplug-simple -gcflags 'all=-N -l' -mod=readonly
	go build -o goplug-simple.so -gcflags 'all=-N -l' -mod=readonly -buildmode=plugin
clean:
	@rm -f goplug-simple goplug-simple.so *.json *.log
