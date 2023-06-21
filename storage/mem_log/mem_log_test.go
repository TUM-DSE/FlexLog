package mem_log

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLog_Append(t *testing.T) {
	log, err := NewLog()
	require.NoError(t, err)
	_ = log.Append("hello", 14)
	_ = log.Commit(14, 33)
	fmt.Println(log.Read(33))

}
