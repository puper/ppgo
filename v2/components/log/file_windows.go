package log

import (
	"os"
	"syscall"
	"time"
)

func FileCreateTime(fi os.FileInfo) time.Time {
	wFileSys := fi.Sys().(*syscall.Win32FileAttributeData)
	tNanSeconds := wFileSys.CreationTime.Nanoseconds() /// 返回的是纳秒
	tSec := tNanSeconds / 1e9                          ///秒
	return time.Unix(tSec, 0)
}
