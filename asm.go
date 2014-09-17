package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
	//"reflect"
)

type Asm struct {
}

func (asm Asm) Output(s string) {
	fmt.Println(s)
}

func (asm Asm) ParseLine(s string) string {
	if s[0] == '#' || s[0:1] == "//" {
		return ""
	} else {
		return s + "\n"
	}
}

func (asm Asm) Parse(reader *bufio.Reader, writer *bufio.Writer) {
	pline := ""
	pnum := 0
	for {
		line, _, err := (*reader).ReadLine()
		pline += string(line)
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if len(pline) == 0 {
			continue
		}
		pnum = strings.Count(pline, "{") - strings.Count(pline, "}")
		if pline[0] == '#' || (pnum == 0 && strings.Count(pline, ";") > 0) {
			writer.WriteString(asm.ParseLine(pline))
			pline = ""
		}
	}
	writer.WriteString(pline + "\n")
}

func (asm Asm) ReadWrite(ifile string, ofile string) {
	fp, err := os.Open(ifile)
	if len(ifile) == 0 || err != nil {
		fp = os.Stdin
	} else {
		defer fp.Close()
	}
	ofp, err := os.OpenFile(ofile, syscall.O_WRONLY|syscall.O_CREAT|syscall.O_TRUNC, 0777)
	defer ofp.Close()
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReaderSize(fp, 4096)
	writer := bufio.NewWriterSize(ofp, 4096)
	defer writer.Flush()
	//fmt.Println(reflect.TypeOf(writer))
	asm.Parse(reader, writer)
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
	fmt.Println(os.Args)
	bfile, tfile, ifile := GetFiles()
	asm.ReadWrite(ifile, tfile)
	args := []string{"-o", bfile, tfile}
	fmt.Println(args)
	out, err := exec.Command("gcc", args...).Output()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("%s\n", out)
}
