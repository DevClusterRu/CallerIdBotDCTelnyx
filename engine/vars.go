package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

var RedisQueryes int = 0
var DeviceErrors int = 0
var NumbersErrors int = 0
var MMSSent int = 0
var SMSSent int = 0


var Conf Config
var Redis RedisObject
var SMSch chan int
var DCch chan int

var DevicesCursor map[string]int

var Devices map[int]Device
type Device struct {
	Number int
	Model   string
	Carrier string
	Brand   string
	Cloud   string
	Gateway string
	Group   int
}

type Config struct {
	Variables map[string]string
}

////////////SMS
type SMSs struct {
	WebHook struct {   //Request from customer prone via Telnyx
		Gateway string `json:"gateway"` //Number of telnyx
		Number  string `json:"number"`	//Number of customer
		Value   string `json:"value"`	//Number for check
	}
}
var SMS SMSs
func (f *SMSs) Parse(raw string)  {
	json.Unmarshal([]byte(raw), &f.WebHook)
}
////////////////


////////////DC
type DCs struct {
	WebHook struct {
		Device_id             int    `json:"device_id"`
		Screenshot            string `json:"screenshot"`
		From_num              string `json:"from_num"`
		Cnam                  string `json:"cnam"`
		DeviceFailed          bool   `json:"deviceFailed"`
		Incoming_number_match bool   `json:"incoming_number_match"`
	}
}
var DC DCs
func (f *DCs) Parse(raw string)  {
	json.Unmarshal([]byte(raw), &f.WebHook)
}
//////////////////

type TelnyxCaller struct {
	Gateway string `json:"from"`
	Number string `json:"to"`
	Message string `json:"text"`
	Image string
	Type string
	Subject string
	Carrier string
}

func (f *TelnyxCaller) Send()  {
	f.Image = strings.Replace(f.Image,"https://call-screens-calleridrep.s3.us-east-2.amazonaws.com/", "https://cs.calleridrep.com/",-1)

	fields:=`{"from": "`+f.Gateway+`", "to": "`+f.Number+`",	"text": "`+f.Message+`"}`
	if f.Type == "MMS" {
		fields = `{"from": "`+f.Gateway+`", "to": "` + f.Number + `",	"subject": "Device `+f.Subject+` -`+f.Carrier+`- Originating Carrier: Century Link",	"text": "` + f.Message + `",	"media_urls": ["` + f.Image + `"]}`
	}

	fmt.Println("Sending to TELNYX this: ",fields)

	req, err := http.NewRequest("POST", Conf.Variables["TELNYX_POINT"], bytes.NewBuffer([]byte(fields)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", Conf.Variables["TELNYX_KEY"])

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	if f.Type == "MMS" {
		MMSSent++
	} else {
		SMSSent++
	}

}

type DCCaller struct {
	Gateway string
	CustomerNumber string
	NumberForCheck string
}

type CustomerQuery struct {
	CustomerNumber string
	Device       Device
}

var Customers map[string][]CustomerQuery //Key - key - number for check, value - array of number of customer

type DCQuery struct {
	CloudCode               string   `json:"cloudCode"`
	ApiKey                  string   `json:"apiKey"`
	Queue                   string   `json:"queue"`
	From                    []string `json:"from"`
	To                      []int    `json:"to"`
	OriginatingCarrierGroup struct {
		Freeswitch struct {
			Centurylink []int `json:"centurylink"`
		} `json:"freeswitch"`
	} `json:"originatingCarrierGroup"`
	Notification_url string `json:"notification_url"`
}

var DCKey string

func (f *DCCaller) Send()  {

	//SMSBOTFROMTX    {"gateway":"+14079041409","number":"+14079041409","value":"+14079041409"}

	dev1:=DeviceByGateway(f.Gateway)
	dev2:=DeviceByGateway(f.Gateway)

	customer1:= CustomerQuery{
		f.CustomerNumber,
		dev1,
	}
	customer2:= CustomerQuery{
		f.CustomerNumber,
		dev2,
	}


	if Customers[f.NumberForCheck] == nil {
		fmt.Println("->I'm try create customer")
		Customers[f.NumberForCheck] = []CustomerQuery{customer1, customer2}
	} else {
		fmt.Println("!!!!!!!!!!I'm try add customer")
		Customers[f.NumberForCheck] = append(Customers[f.NumberForCheck], customer1)
		Customers[f.NumberForCheck] = append(Customers[f.NumberForCheck], customer2)
	}
	fmt.Println("-->",Customers)

	num:=[]string{f.NumberForCheck}
	to:=[]int{customer1.Device.Number, customer2.Device.Number} //Devices

	url := Conf.Variables["DC_POINT"]

	var fields DCQuery
	fields.CloudCode = customer1.Device.Cloud
	fields.ApiKey = DCKey
	if customer1.Device.Group==1 {
		fields.Queue = "high"
	} else {
		fields.Queue = "low"
	}

	fields.From = num
	fields.To = to
	fields.OriginatingCarrierGroup.Freeswitch.Centurylink = to
	fields.Notification_url = Conf.Variables["DC_WEBHOOK"]
	js,_ := json.Marshal(fields)
	fmt.Println(string(js))
	fmt.Println("I'm use ",customer1.Device.Number)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(js))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error on response.\n[ERROR] -", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error while reading the response bytes:", err)
	}
	log.Println(string(body))

}

func DeviceByGateway(gateway string) Device {



	var DevicesPack []Device
	for k, v := range Devices {
		if v.Gateway == gateway {
			DevicesPack = append(DevicesPack, Devices[k])
		}
	}

	if len(DevicesPack)==0{
		return Device{}
	}
	dev := DevicesPack[DevicesCursor[gateway]]
	DevicesCursor[gateway]++
	if DevicesCursor[gateway] > len(DevicesPack)-1 {
		DevicesCursor[gateway] = 0
	}

	return dev
}