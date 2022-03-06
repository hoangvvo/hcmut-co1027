package runner

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hoangvvo/hcmut-co1027/conf"
	"github.com/hoangvvo/hcmut-co1027/unzip"
)

type TestSuite struct {
	Name      string
	Total     int
	TestCases []string
}

func GetSuite(suiteName string) (*TestSuite, error) {
	suiteDir := filepath.Join(conf.SuitesDir, suiteName)
	entries, err := os.ReadDir(suiteDir)
	if err != nil {
		return nil, err
	}
	ts := TestSuite{
		Name:      suiteName,
		TestCases: []string{},
	}
	for _, entry := range entries {
		if entry.IsDir() {
			ts.TestCases = append(ts.TestCases, entry.Name())
		}
	}
	ts.Total = len(ts.TestCases)
	return &ts, nil
}

func GetSuites() ([]TestSuite, error) {
	entries, err := os.ReadDir(conf.SuitesDir)
	if err != nil {
		return nil, err
	}
	var results []TestSuite
	for _, entry := range entries {
		dirName := entry.Name()

		dirsForCnt, _ := os.ReadDir(filepath.Join(conf.SuitesDir, dirName))
		total := len(dirsForCnt) - 1

		results = append(results, TestSuite{
			Name:  dirName,
			Total: total,
		})
	}
	return results, nil
}

func DeleteSuite(suiteName string) error {
	suiteName = filepath.Base(suiteName) // prevent malicious paths
	if suiteName == "case-1" {
		return errors.New("không được xóa cái này :)")
	}

	suiteDir := filepath.Join(conf.SuitesDir, suiteName)
	suiteZip := filepath.Join(conf.ArchiveDir, suiteName+".zip")

	stat, err := os.Stat(suiteDir)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return errors.New("not a directory")
	}

	os.Remove(suiteZip)
	return os.RemoveAll(suiteDir)
}

func fileNameWithoutExtSliceNotation(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func writeFile(filePath string, file multipart.File) error {
	dst, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, file)
	return err
}

func AddSuite(fileName string, file multipart.File) error {

	suiteName := fileNameWithoutExtSliceNotation(filepath.Base(fileName))
	suiteDir := filepath.Join(conf.SuitesDir, suiteName)

	if _, err := os.Stat(suiteDir); !errors.Is(err, os.ErrNotExist) {
		return errors.New("test suite with with the name " + suiteName + " has already existed")
	}

	// add zip file
	zipPath := filepath.Join(conf.ArchiveDir, fileName)
	err := writeFile(zipPath, file)
	if err != nil {
		os.RemoveAll(suiteDir)
		return err
	}

	// unzip
	_, err = unzip.New().Extract(zipPath, suiteDir)
	if err != nil {
		// fail, remove all
		os.RemoveAll(suiteDir)
		return err
	}

	// attempt to convert crlf to lf
	// loop each dirs because arg can be very long for dos2unix
	dirEntries, err := os.ReadDir(suiteDir)
	if err != nil {
	} else {
		for _, entry := range dirEntries {
			if entry.IsDir() {
				exec.Command("bash", "-c", "dos2unix -o "+filepath.Join(suiteDir, entry.Name(), "*.txt")).Run()
			}
		}
	}

	return nil
}
