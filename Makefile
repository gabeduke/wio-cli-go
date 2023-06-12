WIO = $(BIN_DIR)/wio
BIN_DIR := $(CURDIR)/bin
MKDIR_P = mkdir -p
FLAGS ?=

${BIN_DIR}:
	${MKDIR_P} ${BIN_DIR}

$(WIO): $(BIN_DIR)
	go build -o $(WIO)

build: $(WIO)
	go build -o $(WIO)

login: $(WIO)
	$(WIO) user login $(FLAGS)

create-user: $(WIO)
	$(WIO) user create $(FLAGS)

list-nodes: $(WIO)
	$(WIO) nodes list $(FLAGS)

delete-node: $(WIO)
	$(WIO) nodes delete $(FLAGS)

create-node: $(WIO)
	$(WIO) nodes create $(FLAGS)

register-node: $(WIO)
	$(WIO) nodes register $(FLAGS)

clean:
	rm -rf $(BIN_DIR)
