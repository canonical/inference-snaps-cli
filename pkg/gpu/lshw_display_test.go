package gpu

import (
	"fmt"
	"log"
	"testing"
)

func TestGetLocalLsHw(t *testing.T) {
	lsHw, err := GetHostLsHw()
	if err != nil {
		t.Fatalf(err.Error())
	}

	log.Println(string(lsHw))
}

func TestParseLsHw(t *testing.T) {
	lsHw, err := GetHostLsHw()
	if err != nil {
		t.Fatalf(err.Error())
	}

	gpus, err := ParseLsHw(lsHw)
	if err != nil {
		t.Fatalf(err.Error())
	}

	fmt.Printf("%+v\n", gpus)
}
