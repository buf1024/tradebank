bindir=bin
exe=yaodemall

yaodemall=bank/yaodemall/*.go

all:$(exe)


yaodemall: $(yaodemall)
	@echo "building $@"
	go build -gcflags "-N -l" -o $(bindir)/$@ $^
#go build -gcflags "-N -l" --ldflags '-extldflags "-static"' -o $(bindir)/$@ $^

