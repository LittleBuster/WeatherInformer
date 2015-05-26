/*
	Weather Informer v1.0
	Developed by Denisov Sergey aka LittleBuster
	DenisovS21@gmail.com
*/

package log

import "os"
import "fmt"
import "time"
import "strings"


const (
	LOG_ERROR = "ERROR"
	LOG_WARNING = "WARNING"
	LOG_TEXT = "TEXT"
)

func Local( Text string, Module string, Type string ) {
	var (
		Date string
		Out string
	)

	Date = strings.Split(time.Now().String(), ".")[0]
	Out = "[" + Date + "][" + Type + "][" + Module + "] " + Text + "\n"

	fmt.Println(Out)
	f, _ := os.OpenFile("/root/go/src/WeatherInformer/log.txt", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660);
	f.Write([]byte(Out))
	f.Close()
}
