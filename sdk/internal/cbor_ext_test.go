package internal

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ipfs-force-community/go-fvm-sdk/sdk/types"
)

func TestWriteCborArray(t *testing.T) {
	entries := []*types.Entry{
		{
			Flags: 0,
			Key:   "111",
			Codec: 0,
			Value: nil,
		},
		{
			Flags: 0,
			Key:   "222",
			Codec: 0,
			Value: nil,
		},
	}

	writer := bytes.NewBufferString("")
	assert.Nil(t, WriteCborArray(writer, entries))

	reader := bytes.NewBuffer(writer.Bytes())

	newEntries, err := ReadCborArray[types.Entry](reader)
	assert.Nil(t, err)
	assert.Equal(t, entries, newEntries)
}
