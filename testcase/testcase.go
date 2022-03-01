package testcase

import (
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/hoangvvo/hcmut-co1027/conf"
	"github.com/hoangvvo/hcmut-co1027/unzip"
)

type TestSuite struct {
	Name      string
	CreatedAt time.Time
}

func GetSuites() ([]TestSuite, error) {
	entries, err := os.ReadDir(conf.CASEDIR)
	if err != nil {
		return nil, err
	}
	var results []TestSuite
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		name := entry.Name()
		name = name[0 : len(name)-4]
		results = append(results, TestSuite{
			Name:      name,
			CreatedAt: info.ModTime(),
		})
	}
	return results, nil
}

func DeleteSuite(name string) error {
	if name == "case-1" {
		return errors.New("không được xóa cái này :)")
	}
	return os.Remove(filepath.Join(conf.CASEDIR, name+".zip"))
}

func RunSuite(suiteName, answer string) ([]Result, error) {
	dirExtract, err := ioutil.TempDir(os.TempDir(), "testextract")
	if err != nil {
		return nil, err
	}
	defer os.Remove(dirExtract)

	_, err = unzip.New().Extract(filepath.Join(conf.CASEDIR, suiteName+".zip"), dirExtract)
	if err != nil {
		return nil, err
	}

	var caseDirs []string
	entries, err := os.ReadDir(dirExtract)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			caseDirs = append(caseDirs, filepath.Join(dirExtract, entry.Name()))
		}
	}

	return CompileAndRun(caseDirs, answer)
}

func AddSuite(fileName string, file multipart.File) error {
	suiteName := fileName[0 : len(fileName)-4]
	if _, err := os.Stat(filepath.Join(conf.CASEDIR, fileName)); errors.Is(err, os.ErrNotExist) {
		dst, err := os.Create(filepath.Join(conf.CASEDIR, fileName))
		if err != nil {
			return err
		}
		defer dst.Close()
		_, err = io.Copy(dst, file)
		if err != nil {
			return err
		}
	} else {
		return errors.New("test suite with with the name " + suiteName + " has already existed")
	}
	return nil
}
