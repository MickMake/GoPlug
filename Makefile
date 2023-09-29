all:
	@make clean
	@make -C main
	@make -C examples/master
	@make -C examples/plugins

clean:
	echo "# $(pwd) #"
	@make -C main clean
	@make -C examples/master clean
	@make -C examples/plugins clean
