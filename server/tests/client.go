/*
	Weather Informer Test
	Developed by Denisov Sergey aka LittleBuster
	DenisovS21@gmail.com
*/

package main


import "net"
import "fmt"



func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:5000")
	if err != nil {
		fmt.Println("Error connection")
	}
	data := make([]byte, 512)
	
	data = []byte(`{"data": [0,0,1,1,2,2,0]}\`)
	

	conn.Write( data )
	conn.Close()
}