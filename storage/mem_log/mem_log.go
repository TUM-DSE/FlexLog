package mem_log

//#cgo LDFLAGS: -L. -ltbb
//#include "Log.h"
import "C"
import (
	"sync"
	"unsafe"
)

type Log struct {
	mu  sync.Mutex
	Log C.Log
}

func NewLog() (*Log, error) {
	var ret Log
	ret.Log = C.LogNew()

	return &ret, nil
}

func (log *Log) Append(record string, lsn uint64) error {
	var cRecord = C.CString(record)
	log.mu.Lock()
	C.cAppend(log.Log, cRecord, C.ulong(lsn))
	C.free(unsafe.Pointer(cRecord))
	log.mu.Unlock()
	return nil
}

func (log *Log) Commit(lsn uint64, gsn uint64) error {
	log.mu.Lock()
	C.cCommit(log.Log, C.ulong(lsn), C.ulong(gsn))
	log.mu.Unlock()
	return nil
}

func (log *Log) Read(gsn uint64) (string, error) {
	var nextGsn *uint64
	log.mu.Lock()
	var ret = C.GoString(C.cRead(log.Log, C.ulong(gsn), unsafe.Pointer(nextGsn)))
	log.mu.Unlock()
	return ret, nil
}

func (log *Log) Trim(gsn uint64) error {
	log.mu.Lock()
	C.cTrim(log.Log, C.ulong(gsn))
	log.mu.Unlock()
	return nil
}
