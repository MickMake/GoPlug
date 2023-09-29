all:
	@go env GOCACHE
	go build -o GoPlugCmd -gcflags 'all=-N -l' -mod=readonly
clean:
	@rm -f GoPlugCmd
