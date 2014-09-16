package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

type Asm struct {
	s string
}

func (asm Asm) Output(s string) {
	fmt.Println(s)
}

func (asm Asm) Parse(ifile string, ofile string) {
	fp, err := os.Open(ifile)
	if len(ifile) == 0 || err != nil {
		fp = os.Stdin
	} else {
		defer fp.Close()
	}
	ofp, err := os.OpenFile(ofile, syscall.O_WRONLY|syscall.O_CREAT, 0777)
	defer ofp.Close()
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReaderSize(fp, 4096)
	writer := bufio.NewWriterSize(ofp, 4096)
	defer writer.Flush()
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		writer.WriteString(string(line) + "\n")
	}
}

func GetFiles() (string, string, string) {
	bfile := "a.out"
	tfile := "temp.s"
	ifile := ""
	switch len(os.Args) {
	default:
		fallthrough
	case 4:
		bfile = os.Args[3]
		fallthrough
	case 3:
		tfile = os.Args[2]
		fallthrough
	case 2:
		ifile = os.Args[1]
	case 1:
	case 0:
	}
	return bfile, tfile, ifile
}
func main() {
	var asm Asm
	bfile, tfile, ifile := GetFiles()
	asm.Parse(ifile, tfile)
	args := []string{"-o", bfile, tfile}
	fmt.Println(args)
	fmt.Println(os.Args)
	out, err := exec.Command("gcc", args...).Output()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", out)
}
