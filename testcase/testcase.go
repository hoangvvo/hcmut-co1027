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
		results = append(results, TestSuite{
			Name:      name,
			CreatedAt: info.ModTime(),
		})
	}
	return results, nil
}

func DeleteSuite(name string) error {
	name = filepath.Base(name)
	if name == "case-1.zip" {
		return errors.New("không được xóa cái này :)")
	}
	return os.Remove(filepath.Join(conf.CASEDIR, name))
}

func RunSuite(suiteName, answer string) ([]Result, error) {
	dirExtract, err := ioutil.TempDir(os.TempDir(), "testextract")
	if err != nil {
		return nil, err
	}
	defer os.Remove(dirExtract)

	suiteName = filepath.Base(suiteName)
	_, err = unzip.New().Extract(filepath.Join(conf.CASEDIR, suiteName), dirExtract)
	if err != nil {
		return nil, err
	}

	var caseDirs []string
	entries, err := os.ReadDir(dirExtract)
	if err != nil {
		return nil, err
	}

	var customMain []byte

	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}
		if info.IsDir() {
			caseDirs = append(caseDirs, filepath.Join(dirExtract, entry.Name()))
		} else if info.Name() == "main.cpp" {
			customMain, err = os.ReadFile(filepath.Join(dirExtract, entry.Name()))
			if err != nil {
				return nil, err
			}
		}
	}

	return CompileAndRun(caseDirs, answer, customMain)
}

func AddSuite(fileName string, file multipart.File) error {
	suiteName := fileName
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
