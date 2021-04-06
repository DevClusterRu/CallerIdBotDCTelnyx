package engine

import (
	"time"
)

func ListenSMS()  {
	for{
		time.Sleep(time.Second)
		raw:=Redis.GetRow("SMSBOTFROMTX")
		if raw!="" {
			SMS.Parse(raw)
			SMSch<-1
		}
	}
}

func ListenDC()  {
	for{
		time.Sleep(time.Second)
		raw:=Redis.GetRow("SMSBOTFROMDC")
		if raw!="" {
			DC.Parse(raw)
			DCch<-1
		}
	}
}