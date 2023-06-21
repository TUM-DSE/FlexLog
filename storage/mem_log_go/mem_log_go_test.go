package mem_log_go

import (
	"strings"
	"testing"
)

func BenchmarkMemLogGo_Append(b *testing.B) {
	log, _ := NewMemLogGo()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = log.Append(strings.Repeat("h", 4000), uint64(i))
	}
}

func TestLog_Append(t *testing.T) {
	log, _ := NewMemLogGo()
	for i := 0; i < 1000; i++ {
		_ = log.Append(strings.Repeat("h", 4000), 14)
		_ = log.Commit(14, 33)
		_, _ = log.Read(33)
	}
}
