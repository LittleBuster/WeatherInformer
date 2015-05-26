/*
	Weather Informer v1.0
	Developed by Denisov Sergey aka LittleBuster
	DenisovS21@gmail.com
*/

package mariadb

import "fmt"
import "time"
import "strings"
import "strconv"
import "database/sql"
import "WeatherInformer/log"
import _ "github.com/go-sql-driver/mysql"

type Record struct {
	Temp uint64
	Hum uint64
	Date string
}

type MariaDB struct {
	db *sql.DB
}

func (mdb *MariaDB) Connect( ip string, user string, password string, database string ) bool {
	var err error

	mdb.db, _ = sql.Open("mysql", user + ":" + password + "@tcp(" + ip + ":3306)/" + database)
	if err != nil {
		log.Local("MariaDB: Error connection " + ip, "Main", log.LOG_ERROR)
		defer mdb.db.Close()
		return false
	}

	fmt.Println("MariaDB: Connected to " + ip)
	return true
}

func (mdb *MariaDB) AddSensorData( Temp uint64, Hum uint64, Table string ) {
	var (
		Date string
		err error
	)

	Date = strings.Split(time.Now().String(), ".")[0]
	
	stmt, err := mdb.db.Prepare("INSERT INTO " + Table + "(temp, hum, date) VALUES('" + 
		strconv.FormatUint(Temp, 10) + "','" + strconv.FormatUint(Hum, 10) + "','" + Date + "')")
	
	if err != nil {
		log.Local("MariaDB: Can not insert into " + Table, "Main", log.LOG_ERROR)
	}
	defer stmt.Close()
	stmt.Exec()
	fmt.Println("MariaDB: Data inserted in " + Table)
}

func (mdb *MariaDB) AddWaterData( Water uint64, Table string ) {
	var (
		Date string
		err error
	)

	Date = strings.Split(time.Now().String(), ".")[0]
	stmt, err := mdb.db.Prepare("INSERT INTO " + Table + "(water, date) VALUES('" + 
		strconv.FormatUint(Water, 10) + "','" + Date + "')")
	
	if err != nil {
		log.Local("MariaDB: Can not insert into " + Table, "Main", log.LOG_ERROR)
	}
	defer stmt.Close()
	stmt.Exec()
	stmt.Close()
}

func (mdb *MariaDB) Close() {
	mdb.db.Close()
	fmt.Println("MariaDB: Disconnected")
}
