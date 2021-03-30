package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type RotateType int

const (
	RotateTypeDay = iota
	RotateTypeHour
	RotateTypeMinute
	RotateTypeSecond
)

type autoRotateFile struct {
	path       string
	fh         *os.File
	rt         RotateType
	mu         sync.Mutex
	maxBackups int
}

func newAutoRotateFile(path string, rt RotateType) (*autoRotateFile, error) {
	f := new(autoRotateFile)
	f.path = path
	f.rt = rt
	if err := f.open(); err != nil {
		return nil, err
	}

	return f, nil
}

func (a *autoRotateFile) open() error {
	dirPath, _ := path.Split(a.path)
	_, err := os.Stat(dirPath)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, 0777); err != nil {
			return err
		}
	}
	if _, err = os.Stat(a.path); err != nil {
		if os.IsNotExist(err) {
			if tf, err := os.Create(a.path); err != nil {
				return err
			} else {
				tf.Close()
			}
		} else {
			return err
		}
	}

	f, err := os.OpenFile(a.path, os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}

	a.fh = f
	return nil
}

func (a *autoRotateFile) Write(b []byte) (int, error) {
	if err := a.rotate(); err != nil {
		return 0, err
	}
	return a.fh.Write(b)
}

func (a *autoRotateFile) rotate() error {
	a.mu.Lock()
	defer a.mu.Unlock()
	fi, err := a.fh.Stat()
	if err != nil {
		return err
	}
	now := time.Now()
	ct := FileCreateTime(fi)

	rotateFile := ""

	switch a.rt {
	case RotateTypeDay:
		if now.Day() != ct.Day() || now.Month() != ct.Month() || now.Year() != ct.Year() {
			rotateFile = a.path + ct.Format(".20060102")
		}
	case RotateTypeHour:
		if now.Hour() != ct.Hour() || now.Day() != ct.Day() || now.Month() != ct.Month() || now.Year() != ct.Year() {
			rotateFile = a.path + ct.Format(".2006010215")
		}
	case RotateTypeMinute:
		if now.Minute() != ct.Minute() || now.Hour() != ct.Hour() || now.Day() != ct.Day() || now.Month() != ct.Month() || now.Year() != ct.Year() {
			rotateFile = a.path + ct.Format(".200601021504")
		}
	case RotateTypeSecond:
		if now.Second() != ct.Second() || now.Minute() != ct.Minute() || now.Hour() != ct.Hour() || now.Day() != ct.Day() || now.Month() != ct.Month() || now.Year() != ct.Year() {
			rotateFile = a.path + ct.Format(".20060102150405")
		}

	}
	if rotateFile != "" {
		if err := a.fh.Close(); err != nil {
			return err
		}
		if err := os.Rename(a.path, rotateFile); err != nil {
			return err
		}
		if err := a.open(); err != nil {
			return err
		}
		if a.maxBackups > 0 {
			go CleanLogs(a.path, a.maxBackups)
		}
	}

	return nil

}

func CleanLogs(path string, remain int) {
	baseName := filepath.Base(path)
	baseDir := filepath.Dir(path)
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	logFiles := []string{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), baseName) && file.Name() != baseName {
			logFiles = append(logFiles, filepath.Join(baseDir, file.Name()))
		}
	}
	if len(logFiles) > remain {
		for i := 0; i < (len(logFiles) - remain); i++ {
			err := os.Remove(logFiles[i])
			if err != nil {
				fmt.Printf("os.Remove %v error: %v", logFiles[i], err.Error())
			}
		}
	}
}
