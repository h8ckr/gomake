# Copyright © 2019 Ruben Mosblech
# This file belongs to the creator
# You may use, modify and distribute this file as long as you keep this notice

MAKEFILENAME	= makefile-go
MAKEFILE-GO = $(MAKEFILENAME)

.PHONY: all
all: $(MAKEFILE-GO)
	@$(MAKE) -f $(MAKEFILENAME) all || ret=$$?; exit $$ret

%: | $(MAKEFILE-GO)
	@$(MAKE) -f $(MAKEFILENAME) "$@" || ret=$$?; exit $$ret

$(MAKEFILE-GO):
	$(shell curl -o $(MAKEFILENAME) https://gist.githubusercontent.com/h8ckr/27261b540607de4e366f59c508993ea5/raw)
