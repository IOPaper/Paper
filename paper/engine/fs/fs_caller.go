package fs

import (
	"bytes"
	"github.com/IOPaper/Paper/paper/engine"
	"github.com/IOPaper/Paper/utils"
	"github.com/json-iterator/go"
	"github.com/pkg/errors"
	"io"
	"os"
)

type FileApi interface {
	OpenWithIndex(index string) (FileIndexApi, error)

	GetAttachmentPath(key string) (string, error)
	GetAttachment(key string) (io.Reader, error)
	OpenAttachment(key string) (*os.File, error)
}

type FileIndexApi interface {
	OpenDOC() (*engine.Paper, error)
}

type fileApi struct {
	root string
}

func NewFileApi(root string) FileApi {
	if !utils.IsExist(root) {
		utils.Mkdir(root)
	}
	return &fileApi{root: root}
}

func (f *fileApi) withRootPath(path string) string {
	return f.root + path
}

func (f *fileApi) OpenFile(path string) (*os.File, error) {
	return utils.Open(path)
}

func (f *fileApi) Write(path string, b []byte) error {
	return utils.Write(path, bytes.NewBuffer(b))
}

func (f *fileApi) OpenWithIndex(index string) (FileIndexApi, error) {
	indexPath := f.withRootPath(DOCRepo + index)
	if !utils.IsExist(indexPath) {
		return nil, errors.Errorf("index %s not found", indexPath)
	}
	return &fileApi{root: indexPath}, nil
}

func (f *fileApi) GetAttachmentPath(key string) (string, error) {
	attachmentPath := f.withRootPath(AttachmentRepo + key)
	if !utils.IsExist(attachmentPath) {
		return "", errors.Errorf("attachment %s not found", attachmentPath)
	}
	return attachmentPath, nil
}

func (f *fileApi) GetAttachment(key string) (io.Reader, error) {
	attachmentPath, err := f.GetAttachmentPath(key)
	if err != nil {
		return nil, err
	}
	buf := bytes.Buffer{}
	if err = utils.Read(attachmentPath, &buf); err != nil {
		return nil, errors.Wrap(err, "can't open this attachment")
	}
	return &buf, nil
}

func (f *fileApi) OpenAttachment(key string) (*os.File, error) {
	attachmentPath, err := f.GetAttachmentPath(key)
	if err != nil {
		return nil, err
	}
	file, err := utils.Open(attachmentPath)
	if err != nil {
		return nil, errors.Wrap(err, "can't open this attachment")
	}
	return file, nil
}

func (f *fileApi) OpenDOC() (*engine.Paper, error) {
	docPath := f.withRootPath(DOCName)
	if !utils.IsExist(docPath) {
		return nil, errors.Errorf("doc %s not found", docPath)
	}
	buf := bytes.Buffer{}
	if err := utils.Read(docPath, &buf); err != nil {
		return nil, errors.Wrap(err, "can't open this doc")
	}
	var paper engine.Paper
	if err := jsoniter.NewDecoder(&buf).Decode(&paper); err != nil {
		return nil, errors.Wrap(err, "can't parse this doc")
	}
	return &paper, nil
}
