package main

import (
	"io/ioutil"
	"log"
	"math"
	"os"

	"github.com/spf13/viper"
)

// round up
// https://gist.github.com/pelegm/c48cff315cd223f7cf7b
func Round(val float64, roundOn float64, places int) (newVal float64) {
	var round float64
	pow := math.Pow(10, float64(places))
	digit := pow * val
	_, div := math.Modf(digit)
	_div := math.Copysign(div, val)
	_roundOn := math.Copysign(roundOn, val)
	if _div >= _roundOn {
		round = math.Ceil(digit)
	} else {
		round = math.Floor(digit)
	}
	newVal = round / pow
	return
}

func config() {

	if info, err := os.Stat("eporto.conf"); err != nil {
		if os.IsNotExist(err) {
			if err := ioutil.WriteFile("eporto.conf", []byte(emptyConfig), 0644); err != nil {
				log.Fatal("Error writing out empty configuration")
			} else {
				log.Fatal("Written out empty configuration as eporto.conf")
			}
		}
	} else {
		if info.IsDir() {
			log.Fatal("Configuration file path is directory")
		}
	}

	conf, err := os.OpenFile("eporto.conf", os.O_RDONLY, 0600)
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(conf)

	if err != nil {
		log.Fatal(err)
	}

	viper.SetDefault("pageformat", 90)
	viper.SetDefault("listener", ":8080")
	viper.SetDefault("keyphase", "1")

	Listener = viper.GetString("listener")
	PrinterName = viper.GetString("printer")
	PageFormat = viper.GetInt("pageformat")
	PartnerId = viper.GetString("partnerid")
	KeyPhase = viper.GetString("keyphase")
	SignatureKey = viper.GetString("key")
	WalletUser = viper.GetString("user")
	WalletPassword = viper.GetString("password")

	//log.Println(viper.AllSettings())

	if len(Listener) == 0 {
		log.Fatal("Missing listener configuration value")
	}
	if len(PrinterName) == 0 {
		log.Fatal("Missing printer name configuration value")
	}
	if PageFormat == 0 {
		log.Fatal("Missing page format configuration value")
	}
	if len(PartnerId) == 0 {
		log.Fatal("Missing partnerid configuration value")
	}
	if len(KeyPhase) == 0 {
		log.Fatal("Missing keyphase configuration value")
	}
	if len(SignatureKey) == 0 {
		log.Fatal("Missing key configuration value")
	}
	if len(WalletUser) == 0 {
		log.Fatal("Missing user configuration value")
	}
	if len(WalletPassword) == 0 {
		log.Fatal("Missing password configuration value")
	}
}

const emptyConfig = `listener: :8080
printer:
pageformat:
partnerid:
keyphase: 1
key:
user:
password:`
