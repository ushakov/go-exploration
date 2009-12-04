SRCS := map_server.go tiles.go
OBJS := $(SRCS:.go=.6)

map_server: $(OBJS)
	6l -e -o $@ $^

%.6: %.go
	6g -I. $<

map_server.6: tiles.6

clean:
	rm -f *~ *.6

