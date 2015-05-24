/*
	Weather Informer v1.0
	Developed by Denisov Sergey aka LittleBuster
	DenisovS21@gmail.com
*/

package main

import "os"
import "fmt"
import "net"
import "bytes"
import "strconv"
import "encoding/json"
import "WeatherInformer/mariadb"
import "WeatherInformer/log"


type WeatherInfo struct {
	Data []uint64
}

type Config struct {
	Ip string
	Port uint64
	SensorTables []string
	WaterTable string
	DbIp string
	DbUser string
	DbPwd string
	DbName string
}

type Weather struct {
	cfg Config
	Listener net.Listener	
}


func (w *Weather) Load(fname string) {
	finfo, _ := os.Stat(fname)
	sz := finfo.Size()

	f, err := os.Open(fname)
    if err != nil {
        log.Local("Error reading config file", "Main", log.LOG_ERROR)
    }
    defer f.Close()

    b := make([]byte, sz)
    f.Read( b )

    json.Unmarshal(b, &w.cfg)
}

func (w *Weather) Start () {
	var err error

	w.Listener, err = net.Listen("tcp", w.cfg.Ip + ":" + strconv.FormatUint(w.cfg.Port, 10))
	if err != nil {
		log.Local("Can not binding ip address", "Main", log.LOG_ERROR)
	}
	defer w.Listener.Close()

	fmt.Println("Server listen: ", strconv.FormatUint(w.cfg.Port, 10))

	for {
		conn, err := w.Listener.Accept()
		if err != nil {
			log.Local("Can not accept client ", "Main", log.LOG_ERROR)
		}
		go w.ReadData(conn)
	}
}

func (w *Weather) ReadData(c net.Conn) {
	var wInfo WeatherInfo

	data := make([]byte, 512)
	c.Read( data )	

	json.Unmarshal(bytes.Split(data, []byte(`\`))[0], &wInfo)

	//add to Database
	mdb := new(mariadb.MariaDB)
	mdb.Connect(w.cfg.DbIp, w.cfg.DbUser, w.cfg.DbPwd, w.cfg.DbName)
	if wInfo.Data[1] != 0 {
		mdb.AddSensorData(wInfo.Data[0], wInfo.Data[1], w.cfg.SensorTables[0])
	}
	if wInfo.Data[3] != 0 {
		mdb.AddSensorData(wInfo.Data[2], wInfo.Data[3], w.cfg.SensorTables[1])
	}
	if wInfo.Data[5] != 0 {
		mdb.AddSensorData(wInfo.Data[4], wInfo.Data[5], w.cfg.SensorTables[2])
	}
	mdb.AddWaterData(wInfo.Data[6], w.cfg.WaterTable)
	mdb.Close()
	
	c.Close()
}

func main() {	
	var weather Weather
	weather.Load("config.json")
	weather.Start()	
}