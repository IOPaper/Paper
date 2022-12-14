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

type WriteIndex interface {
	WriteIndex(pl *PaperIndex) error
}

type OpenIndex interface {
	OpenIndex() (*PaperIndex, error)
}

type FileApi interface {
	WriteIndex
	OpenIndex

	Init() error

	CheckPaperIndexStatus(index string) bool
	OpenPaperWithIndex(index string) (FileIndexApi, error)
	AddPaper(index string, paper *engine.Paper) error
	AddAttachment(key string, r io.Reader) error

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

func (f *fileApi) Init() error {
	if !utils.IsExist(f.withRootPath(IndexFile)) {
		if err := f.WriteIndex(&PaperIndex{
			Len:   0,
			Index: map[string]*PaperIndexMetadata{},
		}); err != nil {
			return errors.Wrap(err, "index init failed")
		}
	}
	return nil
}

// WriteIndex [RAW API]
func (f *fileApi) WriteIndex(pl *PaperIndex) error {
	buf := bytes.Buffer{}
	if err := jsoniter.NewEncoder(&buf).Encode(pl); err != nil {
		return errors.Wrap(err, "encode paper index failed")
	}
	if err := utils.Write(f.withRootPath(IndexFile), &buf); err != nil {
		return errors.Wrap(err, "write paper index failed")
	}
	return nil
}

// OpenIndex [RAW API]
func (f *fileApi) OpenIndex() (*PaperIndex, error) {
	buf := bytes.Buffer{}
	if err := utils.Read(f.withRootPath(IndexFile), &buf); err != nil {
		return nil, errors.Wrap(err, "can't open paper index")
	}
	var pl PaperIndex
	if err := jsoniter.NewDecoder(&buf).Decode(&pl); err != nil {
		return nil, errors.Wrap(err, "can't decode paper index")
	}
	return &pl, nil
}

func (f *fileApi) CheckPaperIndexStatus(index string) bool {
	return utils.IsExist(f.withRootPath(DOCRepo + index))
}

// OpenPaperWithIndex [RAW API]
func (f *fileApi) OpenPaperWithIndex(index string) (FileIndexApi, error) {
	indexPath := f.withRootPath(DOCRepo + index)
	if !utils.IsExist(indexPath) {
		return nil, errors.Errorf("index %s not found", indexPath)
	}
	return &fileApi{root: indexPath}, nil
}

// AddPaper [RAW API]
func (f *fileApi) AddPaper(index string, paper *engine.Paper) error {
	indexPath := f.withRootPath(DOCRepo + index)
	if utils.IsExist(indexPath) {
		return errors.New("paper index already exists")
	}
	api := &fileApi{root: indexPath}
	return api.writeDOC(paper)
}

// GetAttachmentPath [RAW API]
func (f *fileApi) GetAttachmentPath(key string) (string, error) {
	attachmentPath := f.withRootPath(AttachmentRepo + key)
	if !utils.IsExist(attachmentPath) {
		return "", errors.Errorf("attachment %s not found", attachmentPath)
	}
	return attachmentPath, nil
}

// GetAttachment [RAW API]
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

// OpenAttachment [RAW API]
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

func (f *fileApi) AddAttachment(key string, r io.Reader) error {
	attachmentPath := f.withRootPath(AttachmentRepo + key)
	if utils.IsExist(attachmentPath) {
		return nil
	}
	fs, err := utils.MustOpen(attachmentPath)
	if err != nil {
		return errors.Wrap(err, "can't open attachment")
	}
	io.Copy(fs, r)
	return nil
}

// OpenDOC [RAW API]
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

func (f *fileApi) writeDOC(paper *engine.Paper) error {
	docPath := f.withRootPath(DOCName)
	buf := bytes.Buffer{}
	if err := jsoniter.NewEncoder(&buf).Encode(paper); err != nil {
		return errors.Wrap(err, "can't encode this doc")
	}
	if err := utils.Write(docPath, &buf); err != nil {
		return errors.Wrap(err, "can't write this doc")
	}
	return nil
}
