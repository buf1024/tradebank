package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func gengo(output string, protoDir string, protoFile string) error {
	return nil
}

func genpb(protoc string, protoDir string, protoFile string) error {
	outDir := fmt.Sprintf("--go_out=%s", protoDir)
    pathDir := fmt.Sprintf("--proto_path=%s", protoDir)
	fmt.Printf("\ngenerating protobuf go file...\n")
	cmd := exec.Command(protoc, outDir, pathDir, protoFile)
	var output bytes.Buffer
	cmd.Stdout = &output
    cmd.Stderr = &output
	err := cmd.Run()
	if err != nil {
        return fmt.Errorf("%s", output.String())
	}
    fmt.Printf("protobuf go file generated\n")
	return err
}

func main() {
	pbFile := flag.String("p", "", "the protobuf definition file")
	output := flag.String("o", "command.go", "the output go file")
	help := flag.Bool("h", false, "print the help string")

	flag.Parse()

	if len(*pbFile) <= 0 || *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	_, err := os.Stat(*pbFile)
	if err != nil && os.IsNotExist(err) {
		fmt.Printf("protobuf file %s not exists\n", *pbFile)
		os.Exit(-1)
	}

	file, err := filepath.Abs(*pbFile)
	if err != nil {
		fmt.Printf("get the absolute path of %s failed\n", *pbFile)
		os.Exit(-1)
	}

	fileName := filepath.Base(file)
	i := strings.LastIndex(file, fileName)
	fileDir := file[0:i]

	path, err := exec.LookPath("protoc")
	if err != nil {
		fmt.Printf("can't file protoc executable\n")
		os.Exit(-1)
	}
	protoc, _ := filepath.Abs(path)

	i = strings.LastIndex(fileName, ".")
	if i <= 0 {
		i = len(fileName)
	}
	pbOutFile := fileName[0:i] + ".pb.go"

	cmdOutFile := fileDir + *output

	fmt.Printf("generate info:\n"+
		"protoc        = %s\n"+
		"protobuf dir  = %s\n"+
		"protobuf file = %s\n"+
		"output file   = %s\n",
		protoc, fileDir, file, cmdOutFile)

	err = genpb(protoc, fileDir, file)
	if err != nil {
		fmt.Printf("\ngenerate protobuf go file error:\n%s\n", err)
		os.Exit(-1)
	}
	err = gengo(cmdOutFile, fileDir, file)
	if err != nil {
		fmt.Printf("\ngenerate protobuf command file error:\n%s\n", err)
		os.Exit(-1)
	}
	fmt.Printf("\ngenerate success:\nprotobuf go file = %s\nprotobuf command file = %s\n",
		pbOutFile, cmdOutFile)
}
