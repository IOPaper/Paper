package utils

import (
	"bytes"
	"io"
	"os"
	"strings"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func Remove(path string) error {
	return os.Remove(path)
}

func Move(src, dst string) error {
	return os.Rename(src, dst)
}

func Open(path string) (*os.File, error) {
	{
		dstDir := path[0 : strings.LastIndex(path, "/")+1]
		if !IsExist(dstDir) {
			var err error
			if err = os.MkdirAll(dstDir, os.ModePerm); err != nil {
				return nil, err
			}
		}
	}
	return os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
}

func Write(path string, buf *bytes.Buffer, covered bool) error {
	//{
	//	dstDir := path[0 : strings.LastIndex(path, "/")+1]
	//	if !IsExist(dstDir) {
	//		var err error
	//		if err = os.MkdirAll(dstDir, os.ModePerm); err != nil {
	//			return err
	//		}
	//	}
	//}
	//f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	//if err != nil {
	//	return err
	//}
	if covered {
		Remove(path)
	}
	f, err := Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	n, _ := f.Seek(0, 2)
	_, err = f.WriteAt(buf.Bytes(), n)
	return err
}

func Read(path string, rw io.ReadWriter) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(rw, f)
	return err
}
