APP=redirish
PKG=github.com/mss/$(APP)
SRC=$(APP).go
BIN=bin/$(APP)

all: src/$(PKG)/$(SRC) $(BIN)

run: $(BIN)
	$(BIN) $(ARGS)

$(BIN): $(SRC)
	go install $(PKG)

src/$(PKG)/$(SRC):
	mkdir -p $(dir $(@D))
	ln -s ../../.. $(@D)

