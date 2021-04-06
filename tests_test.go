package main

import (
	"CallerIdBot/engine"
	"regexp"
	"testing"
)

func TestValidate(t *testing.T){
	samples := []string{
		"12305645123",
		"+13022462040",
		"3305655123",
		"430-5655-123",
		"1.560.5565.123",
		"370 5555 123",
		"(714) 788-3435",
	}

	for _,v:=range samples{
		res, err:=engine.Validate(v)
		//reg := regexp.MustCompile(`^\+1\d\d\d\d\d\d\d\d\d$`)
		m, err:= regexp.Match(`^\+1\d{10}$`, []byte(res))

		if (err!=nil){
			t.Error("Error when regex: "+res)
		}

		if (!m) {
			t.Error("String not in format: " + res)
		}

	}

}