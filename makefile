TARGETS := fsmap

all: $(TARGETS)

install:
	install -Dm755 -t "$(DESTDIR)/$(PREFIX)/bin" $(TARGETS)

clean:
	rm -f $(TARGETS)

$(TARGETS): %: cmds/%/main.go
	go build -trimpath -o $@ $<
