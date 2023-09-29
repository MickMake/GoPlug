all:
	@make clean
	make -C main
	make -C examples/master
	make -C examples/plugins

clean:
	@make -C main clean
	@make -C examples/master clean
	@make -C examples/plugins clean
