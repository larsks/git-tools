prefix = /usr/local
bindir = $(prefix)/bin

SCRIPTS = \
	git-edit \
	git-fixup \
	git-mark \
	git-resume \
	git-update-review

INSTALL = install

all:

install: all
	$(INSTALL) -m 755 -d $(bindir)
	for script in $(SCRIPTS); do \
		$(INSTALL) -m 755 $$script $(bindir)/; \
	done
	ln -sf git-mark $(bindir)/git-unmark
