package log_test

import (
	"io/ioutil"
	"os"
	"testing"

	api "github.com/hanmd82/proglog/api/v1"
	"github.com/hanmd82/proglog/internal/log"
	"github.com/stretchr/testify/require"
)

// TestLog(*testing.T) defines a table of tests
func TestLog(t *testing.T) {
	for scenario, fn := range map[string]func(
		t *testing.T, log *log.Log,
	){
		"append and read a record":    testAppendRead,
		"offset out of range error":   testOutOfRangeErr,
		"init with existing segments": testInitExisting,
	} {
		t.Run(scenario, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "store-test")
			require.NoError(t, err)
			defer os.RemoveAll(dir)

			c := log.Config{}
			c.Segment.MaxStoreBytes = 32

			log, err := log.NewLog(dir, c)
			require.NoError(t, err)

			fn(t, log)
		})
	}
}

// testAppendRead(*testing.T, *log.Log) tests successfully appending to and reading from the log.
func testAppendRead(t *testing.T, log *log.Log) {
	record := &api.Record{
		Value: []byte("hello world"),
	}

	off, err := log.Append(record)
	require.NoError(t, err)
	require.Equal(t, uint64(0), off)

	read, err := log.Read(off)
	require.NoError(t, err)
	require.Equal(t, record, read)
}

// testOutOfRangeErr(*testing.T, *log.Log) tests that the log returns an error when trying to read an offset
// that’s outside of the range of offsets the log has stored.
func testOutOfRangeErr(t *testing.T, log *log.Log) {
	read, err := log.Read(1)
	require.Nil(t, read)
	apiErr := err.(api.ErrOffsetOutOfRange)
	require.Equal(t, uint64(1), apiErr.Offset)
}

// testInitExisting(*testing.T, *log.Log) tests that creating a log bootstraps
// from the data stored by prior log instances.
func testInitExisting(t *testing.T, o *log.Log) {
	record := &api.Record{
		Value: []byte("hello world"),
	}

	for i := 0; i < 3; i++ {
		_, err := o.Append(record)
		require.NoError(t, err)
	}
	require.NoError(t, o.Close())

	n, err := log.NewLog(o.Dir, o.Config)
	require.NoError(t, err)

	off, err := n.Append(record)
	require.NoError(t, err)
	require.Equal(t, uint64(3), off)
}
