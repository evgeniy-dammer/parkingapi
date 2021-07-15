package models

//GetCountObjects return count of rows in table
func GetCountObjects(table string) int {
	var counter int
	err := GetDB().QueryRow("SELECT count(*) FROM " + table + ";").Scan(&counter)

	if err != nil {
		panic(err)
	}

	return counter
}
