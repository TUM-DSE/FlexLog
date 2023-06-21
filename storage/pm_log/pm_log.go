package pm_log

//#cgo LDFLAGS: -L. -lstorage -lpmemobj
//#include "PMLog.h"
import "C"
import "unsafe"

type PMLog struct {
	pmLog C.PMLog
}

func NewPMLog() (*PMLog, error) {
	var ret PMLog
	ret.pmLog = C.startUp()

	return &ret, nil
}

func (log *PMLog) Append(record string, lsn uint64) error {
	var cRecord = C.CString(record)
	C.cAppend(log.pmLog, cRecord, C.ulong(lsn))
	C.free(unsafe.Pointer(cRecord))

	return nil
}

func (log *PMLog) Commit(lsn uint64, gsn uint64) error {
	C.cCommit(log.pmLog, C.ulong(lsn), C.ulong(gsn))
	return nil
}

func (log *PMLog) Read(gsn uint64) (string, error) {
	var nextGsn *uint64	
	var ret = C.GoString(C.cRead(log.pmLog, C.ulong(gsn), unsafe.Pointer(nextGsn)))	
	return ret, nil	
}

func (log *PMLog) Trim(gsn uint64) error {
	C.cTrim(log.pmLog, C.ulong(gsn))
	return nil
}
