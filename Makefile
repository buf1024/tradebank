bindir=./bin
exe=ioms

ioms_go=./ioms/main/ioms.go

all:$(exe)


ioms: $(ioms_go)
	@echo "building $@"
	go build -o $(bindir)/$@ $^

