package engine

import (
	"errors"
	"regexp"
	"strings"
)

func checkLen(s string) string  {
	m, _:= regexp.Match(`^\+1\d{10}$`, []byte(s))
	if (m){
		return s
	} else {
		return ""
	}
}

func Validate(number string)  (string, error) {
	//Find first entry - digit
	reg := regexp.MustCompile(`\d.*\d`)
	valid:=reg.FindString(number)
	valid = strings.Replace(valid,".","",-1)
	valid = strings.Replace(valid," ","",-1)
	valid = strings.Replace(valid,"-","",-1)
	valid = strings.Replace(valid,"(","",-1)
	valid = strings.Replace(valid,")","",-1)
	validB:=[]byte(valid)
	if len(validB)==0{
		return "", errors.New("Something wrong in format:"+ number)
	}
	if string(validB[0])!="1" && string(validB[0])!="+"{
		valid = "+1"+valid
		if checkLen(valid)!=""{
			return valid, nil
		} else {
			return "", errors.New("Something wrong in format")
		}
	}

	if string(validB[0])!="+"{
		valid = "+"+valid
		if checkLen(valid)!=""{
			return valid, nil
		} else {
			return "", errors.New("Something wrong in format")
		}
	}

	return valid, nil

}