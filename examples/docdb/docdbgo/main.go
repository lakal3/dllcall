//

//go:generate dllcall dbif.go ../docdbcpp/dbif.h
package main

import "log"

func main() {
	err := load_dbif("docdbdll.dll")
	if err != nil {
		log.Fatal(err)
	}
	d := &dbIf{dbName: "test.db"}
	err = d.Open()
	if err != nil {
		log.Fatal(err)
	}
	/// ... Do something with db
	err = d.Close()
	if err != nil {
		log.Fatal(err)
	}
}
