BINARIES := fsmap

install:
	install -Dm755 -t "$(DESTDIR)/$(PREFIX)/bin" $(BINARIES)
