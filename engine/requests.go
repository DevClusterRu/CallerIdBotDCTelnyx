package engine

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type DCKeyResponse struct {
	Message int    `json:"message"`
	ApiKey  string `json:"ApiKey"`
	Result  bool   `json:"result"`
}

func PostDCKeyGen() {

	url := "https://u5rbap0zmc.execute-api.us-east-2.amazonaws.com/PROD/auth/createKey"
	fields := `{"lifetime": 100000}`
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(fields)))
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

	var DCKR DCKeyResponse
	json.Unmarshal(body, &DCKR)
	if DCKR.Result == true {
		DCKey = DCKR.ApiKey
		fmt.Println(DCKey)
	} else {
		time.Sleep(time.Second * 5)
		PostDCKeyGen()
	}
}

func RefrechDCKey() {
	for {
		PostDCKeyGen()
		time.Sleep(time.Minute * 60)
	}
}

func TelnyxIncoming() {

	fmt.Println("Telnyx incoming!")
	fmt.Println(SMS.WebHook)
	fmt.Println("----------------------------------")
	telnyx := TelnyxCaller{
		Gateway: SMS.WebHook.Gateway,
		Number:  SMS.WebHook.Number,
		Message: "CID: Testing in progress...",
		Type:    "SMS",
	}
	telnyx.Send()




	fmt.Println("Send to Telnyx incoming!")
	fmt.Println(SMS.WebHook)
	fmt.Println("----------------------------------")
	numberValid, err := Validate(SMS.WebHook.Value)

	if err != nil {
		fmt.Println("Wrong number! Use number of customer", err)
		numberValid = SMS.WebHook.Number
	}

	//Wait key from DC
	for {
		if DCKey != "" {
			break
		} else {
			time.Sleep(time.Second * 2)
			fmt.Print(".")
		}
	}

	dcCaller := DCCaller{
		Gateway:        SMS.WebHook.Gateway,
		CustomerNumber: SMS.WebHook.Number,
		NumberForCheck: numberValid,
	}
	dcCaller.Send()

}

func DCIncoming() {
	if DC.WebHook.From_num != "" && DC.WebHook.Screenshot != "" {

		fmt.Println("<--------------", DC.WebHook)

		//Now find recipients
		customerArray := Customers[DC.WebHook.From_num]

		for i, v := range customerArray {
			cnam := ""
			if DC.WebHook.Cnam != "" {
				cnam = "; CNAM: " + DC.WebHook.Cnam
			}
			if DC.WebHook.Device_id == v.Device.Number {

				//IF FAIL
				if DC.WebHook.DeviceFailed {

					log.Println("Device fail!")

					telnyx := TelnyxCaller{
						Gateway: v.Device.Gateway,
						Number:  v.CustomerNumber,
						Subject: "CID: Device error",
						Image:   "https://cs.calleridrep.com/devicefailed.jpg",
						Type:    "MMS",
						Message: "CID: Device error",
						Carrier: v.Device.Carrier,
					}
					go telnyx.Send()
					DeviceErrors++
					continue
				}

				if !DC.WebHook.Incoming_number_match {

					log.Println("Number fail!")

					telnyx := TelnyxCaller{
						Gateway: v.Device.Gateway,
						Number:  v.CustomerNumber,
						Subject: "CID: Device error",
						Image:   "https://cs.calleridrep.com/callfailed.png",
						Type:    "MMS",
						Message: "CID: Number error",
						Carrier: v.Device.Carrier,
					}
					go telnyx.Send()
					NumbersErrors++
					continue
				}

				//If SUCCESS

				log.Println("Success!")

				telnyx := TelnyxCaller{
					Gateway: v.Device.Gateway,
					Number:  v.CustomerNumber,
					Subject: v.Device.Brand,
					Image:   DC.WebHook.Screenshot,
					Type:    "MMS",
					Message: "CID: Number scanned: " + DC.WebHook.From_num + cnam + " - Thank you for using CallerIDReputation.com",
					Carrier: v.Device.Carrier,
				}
				go telnyx.Send()

				//Remove this recipient
				//Copy last value into this key
				if len(customerArray) > 0 {
					customerArray[i] = customerArray[len(customerArray)-1]
					//Erase last value
					customerArray = customerArray[:len(customerArray)-1]
					//Rewrite Customers
					Customers[DC.WebHook.From_num] = customerArray
					//If this be last customer, erase key
					if len(Customers[DC.WebHook.From_num]) == 0 {
						delete(Customers, DC.WebHook.From_num)
					}
				}
			}

		}
		//Now we clean recipients for this number
		fmt.Println("Now map is: ", Customers)

	}
}
