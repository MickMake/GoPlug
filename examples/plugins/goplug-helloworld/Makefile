all:
	go env GOCACHE
	go build -o goplug-helloworld -gcflags 'all=-N -l' -mod=readonly
	go build -o goplug-helloworld.so -gcflags 'all=-N -l' -mod=readonly -buildmode=plugin
clean:
	@rm -f goplug-helloworld goplug-helloworld.so *.json *.log
