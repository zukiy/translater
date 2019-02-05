package yandex

import (
	"fmt"
	"testing"
)

func TestUnit_Translate(t *testing.T) {
	c := New("trnsl.1.1.20190205T140220Z.1c75b6a3a6d8311c.23fa815511980868b6b5d09eacdfccc1b87ccead", "v1.5")
	res, err := c.Translate("calm", "ru")
	if err != nil {
		println(err.Error())
	}
	println(fmt.Sprintf("%+v", res))
}
