.PHONY: all clean re
NAME=npuzzle

all:
	go build -o $(NAME)

clean:
	rm $(NAME)

re:	clean all
