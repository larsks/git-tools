prefix = $(HOME)
bindir = $(prefix)/bin

SCRIPTS = \
	git-edit \
	git-fixup \
	git-mark \
	git-resume \
	git-update-review \
	git-since

INSTALL = install

all: $(SCRIPTS)

install: all
	$(INSTALL) -m 755 -d $(bindir)
	for script in $(SCRIPTS); do \
		$(INSTALL) -m 755 $$script $(bindir)/; \
	done
	ln -sf git-mark $(bindir)/git-unmark
