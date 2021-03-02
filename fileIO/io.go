package fileIO

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func ReadFileAsString(path string) (string, error) {
	logFile := ""
	fi, err := os.Open(path)
	if err != nil {
		return "", err
	}

	defer fi.Close()
	r := bufio.NewReader(fi)

	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
			//return
		}
		if 0 == n {
			break
		} else {
			// 将读取到的数据交给 callback 处理
			logFile += string(buf[:n])
		}
	}
	return logFile, nil
}

func WriteStringToFile(path string, stringData string, append bool) error {
	var fl *os.File
	var err error
	if append {
		fl, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
	} else {
		fl, err = os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	}
	if err != nil {
		return err
	}
	//err = fl.Chdir()
	//if err != nil {
	//	//fmt.Println(err.Error())
	//	return err
	//}

	defer fl.Close()

	n, err := fl.WriteString(stringData)
	if err != nil {
		//fmt.Println(err.Error())
		return err
	}
	if n < len(stringData) {
		//fmt.Println("write byte num error")
		return fmt.Errorf("write byte num error")
	}
	return nil
}
