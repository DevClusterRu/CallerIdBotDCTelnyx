package engine

import (
	"github.com/kardianos/osext"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func (f *Config) Init() {
	Customers = make(map[string][]CustomerQuery)
	Devices = make(map[int]Device)
	DevicesCursor = make(map[string]int)
	f.Variables = make(map[string]string)


	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(folderPath + "/settings.conf")
	//file, err := os.Open("/var/www/caller/calleridrep/src/.env")

	if err != nil {
		panic(err)
	}
	b, _ := ioutil.ReadAll(file)
	rows := strings.Split(string(b), "\n")
	for _, v := range rows {
		if v == "" || strings.Index(v, "#") == 0 {
			continue
		}
		pair := strings.Split(v, "=")
		if len(pair) == 2 {
			pair[0] = strings.TrimSpace(pair[0])
			pair[1] = strings.TrimSpace(pair[1])

			if strings.ToUpper(pair[0]) == "DEVICE" {
				deviceProperties := strings.Split(pair[1], ",")
				if len(deviceProperties) != 7 {
					panic("Wrong device properties in settings.conf")
				}

				devId, _ := strconv.Atoi(strings.TrimSpace(deviceProperties[0]))
				grp, _ := strconv.Atoi(strings.TrimSpace(deviceProperties[6]))

				Devices[devId] = Device{
					devId,
					strings.TrimSpace(deviceProperties[1]),
					strings.TrimSpace(deviceProperties[2]),
					strings.TrimSpace(deviceProperties[3]),
					strings.TrimSpace(deviceProperties[4]),
					strings.TrimSpace(deviceProperties[5]),
					grp,
				}

				_, ok := DevicesCursor[strings.TrimSpace(deviceProperties[5])]
				if !ok {
					DevicesCursor[strings.TrimSpace(deviceProperties[5])] = 0
				}

			} else {
				f.Variables[pair[0]] = pair[1]
			}
		}
	}
}
