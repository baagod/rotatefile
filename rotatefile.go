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
	rotateAt       time.Time
	mu             sync.Mutex
}

func New(filepath string, rotate RotateInterval) (w *Writer, err error) {
	err = os.MkdirAll(path.Dir(filepath), DefaultFilePerm) // 创建文件目录
	if err != nil && !os.IsExist(err) {
		return
	}

	stat, err := os.Stat(filepath) // 获取文件信息
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	now := time.Now()
	w = &Writer{filepath: filepath, rotateInterval: rotate}

	if stat != nil { // 判断文件是否需要现在轮转
		mod := stat.ModTime()
		if w.rotateAt = rotate.next(mod); now.After(w.rotateAt) {
			if err = w.rename(mod); err != nil { // 归档文件
				return
			}
		}
	}

	w.file, err = os.OpenFile(w.filepath, DefaultFileFlags, DefaultFilePerm)
	if err == nil {
		w.rotateAt = rotate.next(now) // 设置下次轮转时间
	}

	return
}

func (w *Writer) doRotate() (err error) {
	if now := time.Now(); now.After(w.rotateAt) { // 需要轮转
		if err = w.rename(now); err != nil {
			return
		}
		w.rotateAt = w.rotateAt.Add(time.Duration(w.rotateInterval) * time.Second)
		w.file, err = os.OpenFile(w.filepath, DefaultFileFlags, DefaultFilePerm)
	}
	return
}

func (w *Writer) Write(b []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if err = w.doRotate(); err != nil {
		return
	}
	return w.file.Write(b)
}

func (w *Writer) rename(now time.Time) error {
	if w.file != nil {
		_ = w.file.Close()
	}
	oldpath := w.filepath
	newpath := oldpath + "." + now.Format("20060102150405")
	return os.Rename(oldpath, newpath)
}

func (interval RotateInterval) next(now time.Time) time.Time {
	y, m, d := now.Year(), now.Month(), now.Day()
	switch interval {
	case PerDay:
		return time.Date(y, m, d+1, 0, 0, 0, 0, time.Local)
	case PerHour:
		return time.Date(y, m, d, now.Hour()+1, 0, 0, 0, time.Local)
	case PerMinute:
		return time.Date(y, m, d, now.Hour(), now.Minute()+1, 0, 0, time.Local)
	}
	return time.Date(y, m, d+1, 0, 0, 0, 0, time.Local)
}
