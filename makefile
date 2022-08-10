psm = montecarloclub.kanban.kbd

idl_path = $(GOPATH)/src/github.com/MonteCarloClub/KBD/kdb.thrift
psm_path = github.com/MonteCarloClub/KBD

help: Makefile
	@echo
	@echo " Choose a command run:"
	@echo "   build    编译项目"
	@echo "   clean    清理文件"
	@echo "   run      运行项目"
	@echo "   cli      生成客户端代码"
	@echo "   server   生成server端代码"


build:
	sh ./build.sh

clean:
	rm -rf output

run: build
	./output/bootstrap.sh

cli:
	kitex -module $(psm_path)  ${idl_path}

idl:
	kitex -module $(psm_path) -service $(psm)  ${idl_path}

