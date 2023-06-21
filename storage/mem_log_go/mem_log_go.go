package mem_log_go

import (
	"fmt"
	"github.com/nathanieltornow/PMLog/storage"
	"sync"
)

type MemLogGo struct {
	lsnMu              sync.Mutex
	lsnToWaitingRecord map[uint64]chan string
	gsnMu              sync.RWMutex
	gsnToRecord        map[uint64]string
}

func Create() (storage.Log, error) {
	return NewMemLogGo()
}

func NewMemLogGo() (*MemLogGo, error) {
	mlg := new(MemLogGo)
	mlg.lsnToWaitingRecord = make(map[uint64]chan string)
	mlg.gsnToRecord = make(map[uint64]string)
	return mlg, nil
}

func (mlg *MemLogGo) Append(record string, lsn uint64) error {
	mlg.lsnMu.Lock()
	recCh, ok := mlg.lsnToWaitingRecord[lsn]
	if !ok {
		recCh = make(chan string, 1)
		mlg.lsnToWaitingRecord[lsn] = recCh
	}
	mlg.lsnMu.Unlock()
	recCh <- record
	return nil
}

// Commit commits the records which is stored on local-sequence-number lsn with the global-sequence-number gsn on
// on the log of color. The record can be committed for multiple colors.
func (mlg *MemLogGo) Commit(lsn uint64, gsn uint64) error {
	mlg.lsnMu.Lock()
	recCh, ok := mlg.lsnToWaitingRecord[lsn]
	if !ok {
		recCh = make(chan string, 1)
		mlg.lsnToWaitingRecord[lsn] = recCh
	}
	mlg.lsnMu.Unlock()
	rec := <-recCh
	var counter int
	for i := 0; i < 15000; i++ {
		counter++
	}
	mlg.gsnMu.Lock()
	mlg.gsnToRecord[gsn] = rec
	mlg.gsnMu.Unlock()
	defer func() {
		mlg.lsnMu.Lock()
		delete(mlg.lsnToWaitingRecord, lsn)
		mlg.lsnMu.Unlock()
	}()
	return nil
}

// Read returns the record stored on global-sequence-number gsn on the log of color
func (mlg *MemLogGo) Read(gsn uint64) (string, error) {
	mlg.gsnMu.RLock()
	defer mlg.gsnMu.RUnlock()
	rec, ok := mlg.gsnToRecord[gsn]
	if !ok {
		return "", fmt.Errorf("failed to find record")
	}
	return rec, nil
}

// Trim deletes all records of the log of color before global-sequence-number gsn
func (mlg *MemLogGo) Trim(gsn uint64) error {
	if gsn == 0 {
		mlg.gsnMu.Lock()
		mlg.gsnToRecord = make(map[uint64]string)
		mlg.gsnMu.Unlock()
		return nil
	}
	mlg.gsnMu.Lock()
	for recGsn := range mlg.gsnToRecord {
		if recGsn < gsn {
			delete(mlg.gsnToRecord, recGsn)
		}
	}
	mlg.gsnMu.Unlock()
	return nil
}
