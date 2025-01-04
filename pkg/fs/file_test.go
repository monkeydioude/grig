package fs

import "testing"

func TestICanTestFileAppending(t *testing.T) {
	pathToAppend := "./service"
	filesDirToUse := "/A/B/C"

	if AppendToThisFileDirectory(pathToAppend, filesDirToUse) != "/A/B/service" {
		t.Fail()
	}
}
