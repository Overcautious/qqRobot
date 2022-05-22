package setting

import (
	"github.com/go-ini/ini"
	"log"
)

var Cfg *ini.File

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}
}
