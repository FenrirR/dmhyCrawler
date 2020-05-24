package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileManager(t *testing.T) {
	data := []string{
		"line1",
		"line2",
	}
	path := "./test/testdata"
	CreateFile(path)
	SaveList2Txt(data, path)
	res := ReadTxt2Set(path)
	for _, d := range data {
		assert.True(t, res[d], "should read the data write in")
	}

	n, err := RemoveContents("./test")
	assert.Equal(t, 1, n, "should remove test file")
	assert.Nil(t, err, "should not return error")
}
