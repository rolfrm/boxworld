package main
import "os"
import "strings"
import "fmt"

type Iterator interface {
	Next()(*Iterator)
	GetValue()(interface{})
}

func LoadFileToString(path string)(output string){
	file,err :=os.Open(path)
	if err == nil{
		bufflen := 16
		buff := make([]byte ,bufflen)
		n, rdErr := file.Read(buff)
		for rdErr == nil {
			output = strings.Join([]string{output, string(buff[0:n])},"")
			n, rdErr = file.Read(buff)
			}
		file.Close()
	}
	return
}

func PrintNReturn(in string)(out string){
	out = in
	fmt.Println(in)
	return
}
