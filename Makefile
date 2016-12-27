PROJECT_NAME=rtorrent-cleaner
PROJECT_VERSION=0.1.0
BIN_DIR=./bin
RUN_ARGS=-ruTorrent=http://localhost -dir=/home/user/Images

GO=go
GOFMT=gofmt

all: compile run

compile:
	@echo "############################################################"
	@echo "#                         COMPILE                          #"
	@echo "############################################################"
	${GOFMT} -d -w src/github.com/kuroidoruido/rtorrent-cleaner/
	${GOFMT} -d -s src/github.com/kuroidoruido/rtorrent-cleaner/
	${GO} install github.com/kuroidoruido/rtorrent-cleaner

run:
	@echo "############################################################"
	@echo "#                           RUN                            #"
	@echo "############################################################"
	@echo
	@${BIN_DIR}/${PROJECT_NAME} ${RUN_ARGS}
	@echo

install:
	@echo "############################################################"
	@echo "#                   INSTALL DEPENDENCIES                   #"
	@echo "############################################################"

	${GO} get github.com/gorilla/rpc
	${GO} get github.com/divan/gorilla-xmlrpc/xml

release: compile
	@echo "############################################################"
	@echo "#                         RELEASE                          #"
	@echo "############################################################"

	@cp -f ${BIN_DIR}/${PROJECT_NAME} ${PROJECT_NAME}-${PROJECT_VERSION}
	@echo " RELEASE =>" ${PROJECT_NAME}-${PROJECT_VERSION}
