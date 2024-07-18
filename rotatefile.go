package rotatefile

import (
	"os"
	"path"
	"sync"
	"time"
)

type RotateInterval int

var (
	DefaultFilePerm  os.FileMode = 0664
	DefaultFileFlags             = os.O_CREATE | os.O_WRONLY | os.O_APPEND
)

const (
	PerSecond RotateInterval = 1
	PerMinute                = 60 * PerSecond
	PerHour                  = 60 * PerMinute
	PerDay                   = 24 * PerHour
)

type Writer struct {
	file           *os.File
	filepath       string
	rotateInterval RotateInterval
	rotateNext     time.Time
	mu             sync.Mutex
}

func New(filepath string, rotateInterval RotateInterval) (w *Writer, err error) {
	err = os.MkdirAll(path.Dir(filepath), DefaultFilePerm)
	if err == nil || os.IsExist(err) {
		w = &Writer{filepath: filepath, rotateInterval: rotateInterval}
	}
	return
}

func (f *Writer) doRotate() (err error) {
	if now := time.Now(); now.After(f.rotateNext) { // 需要轮转
		if f.file != nil {
			_ = f.file.Close() // 关闭打开的文件
			oldpath := f.file.Name()
			newpath := oldpath + "." + now.Format("20060102150405")
			if err = os.Rename(oldpath, newpath); err != nil { // 重命名文件
				return
			}
		}

		y, m, d := now.Year(), now.Month(), now.Day()
		switch f.rotateInterval {
		case PerDay:
			f.rotateNext = time.Date(y, m, d+1, 0, 0, 0, 0, time.Local)
		case PerHour:
			f.rotateNext = time.Date(y, m, d, now.Hour()+1, 0, 0, 0, time.Local)
		case PerMinute:
			f.rotateNext = time.Date(y, m, d, now.Hour(), now.Minute()+1, 0, 0, time.Local)
		}

		f.file, err = os.OpenFile(f.filepath, DefaultFileFlags, DefaultFilePerm)
	}

	return
}

func (f *Writer) Write(b []byte) (n int, err error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err = f.doRotate(); err != nil {
		return
	}

	return f.file.Write(b)
}
