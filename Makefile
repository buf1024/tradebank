bindir=bin
exe=yaodemall

yaodemall=bank/yaodemall/*.go

all:$(exe)


yaodemall: $(yaodemall)
	@echo "building $@"
	go build --ldflags '-extldflags "-static"' -o $(bindir)/$@ $^

