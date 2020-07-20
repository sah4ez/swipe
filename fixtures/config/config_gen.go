//+build !swipe

// Code generated by Swipe v1.22.4. DO NOT EDIT.

//go:generate swipe
package config

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadConfig() (cfg *Config, errs []error) {
	cfg = &Config{
		Bind: "hohoho",
		Name: "Default MethodName",
	}
	flag.StringVar(&cfg.Bind, "bind-addr", cfg.Bind, "")
	if cfg.Bind == "" {
		errs = append(errs, errors.New("flag bind-addr required"))
	}
	cfgNameTmp, ok := os.LookupEnv("NAME")
	if ok {
		cfg.Name = cfgNameTmp
	}
	cfgMaxPriceTmp, ok := os.LookupEnv("MAX_PRICE")
	if ok {
		maxpriceInt, err := strconv.Atoi(cfgMaxPriceTmp)
		if err != nil {
			errs = append(errs, fmt.Errorf("convert MAX_PRICE error: %w", err))
		}
		cfg.MaxPrice = int(maxpriceInt)
	}
	if cfg.MaxPrice == 0 {
		errs = append(errs, errors.New("env MAX_PRICE required"))
	}
	cfg.DB = DB{}
	cfgDBConnTmp, ok := os.LookupEnv("DB2_CONN")
	if ok {
		cfg.DB.Conn = cfgDBConnTmp
	}
	cfg.DB.Foo = Foo{}
	cfgDBFooNameTmp, ok := os.LookupEnv("DB2_FOO_NAME")
	if ok {
		cfg.DB.Foo.Name = cfgDBFooNameTmp
	}
	cfgURLsTmp, ok := os.LookupEnv("URLS")
	if ok {
		partsurlsInt := strings.Split(cfgURLsTmp, ",")
		cfg.URLs = make([]int, len(partsurlsInt))
		for i, s := range partsurlsInt {
			tmpInt, err := strconv.Atoi(s)
			if err != nil {
				errs = append(errs, fmt.Errorf("convert URLS error: %w", err))
			}
			cfg.URLs[i] = int(tmpInt)
		}
	}
	flag.Parse()
	return
}

func (cfg *Config) String() string {
	out := `
--bind-addr ` + fmt.Sprintf("%v", cfg.Bind) + `
NAME=` + fmt.Sprintf("%v", cfg.Name) + `
MAX_PRICE=` + fmt.Sprintf("%v", cfg.MaxPrice) + `
DB2_CONN=` + fmt.Sprintf("%v", cfg.DB.Conn) + `
DB2_FOO_NAME=` + fmt.Sprintf("%v", cfg.DB.Foo.Name) + `
URLS=` + fmt.Sprintf("%v", cfg.URLs) + `
`
	return out
}
