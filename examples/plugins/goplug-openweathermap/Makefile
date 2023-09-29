all:
	go env GOCACHE
	go build -o goplug-openweathermap -gcflags 'all=-N -l' -mod=readonly
	go build -o goplug-openweathermap.so -gcflags 'all=-N -l' -mod=readonly -buildmode=plugin
clean:
	@rm -f goplug-openweathermap goplug-openweathermap.so *.json *.log
