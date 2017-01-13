bindir=bin
exe=yoyitd

yoyitd_go=ioms/bank/yoyitd/main/ioms.go

all:$(exe)


yoyitd: $(yoyitd_go)
	@echo "building $@"
	go build -o $(bindir)/$@ $^

