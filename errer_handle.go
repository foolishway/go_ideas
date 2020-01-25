package main
import (
	"io"
	"log"
	"os"
)
type Writter struct {
	w io.Writer
	e error
}

func (W *Writter) Write(bs []byte)  {
	if W.e != nil {
		return
	}
	_, e := W.w.Write(bs)
	W.e = e
}
func main() {
	f, e := os.OpenFile("./test.txt", os.O_WRONLY, os.ModeAppend)
	if e != nil {
		log.Fatalf("open file error %v", e)
	}
	defer f.Close()
	//io.Copy(os.Stdout, f)
	w := &Writter{f, nil}
	w.Write([]byte("Hello world,"))
	w.Write([]byte("I am a newbie in GO,"))
	w.Write([]byte("MY NAME IS WEI."))
	if w.e != nil {
		log.Fatalf("Write error:%v", w.e)
	}
}
