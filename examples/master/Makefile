all:
	go env GOCACHE
	go build -o master -gcflags 'all=-N -l' -mod=readonly
clean:
	@rm -f master *.log *.json
