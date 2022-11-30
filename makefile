TARGETS := fsmap
LIBRARIES := libfsmap.a

all: $(TARGETS) $(LIBRARIES)

install:
	install -Dm755 -t "$(DESTDIR)/$(PREFIX)/bin" $(TARGETS)
	install -Dm644 -t "$(DESTDIR)/$(PREFIX)/lib" $(LIBRARIES)

clean:
	rm -f $(TARGETS)

%.o: %.go
	gccgo -Wall -Werror -c $^ -o $@

lib%.a: %.o
	ar rcs $@ $^

fsmap: main.go fsmap.o
	gccgo -Wall -Werror $^ -o $@
