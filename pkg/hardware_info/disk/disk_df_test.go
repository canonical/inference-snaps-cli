package disk

import (
	"testing"
)

func TestHostSnapDir(t *testing.T) {
	dfData, err := hostDf("/", "/var/snap/snapd")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dfData)

	dirInfos, err := parseDf(dfData)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(dirInfos)
}
