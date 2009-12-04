SRCS := map_server.go tiles.go
OBJS := $(SRCS:.go=.6)

%.6: %.go
	6g -I. $<

map_server: $(OBJS)
	6l -e -o $@ $^

