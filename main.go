package main

import "forum/utils"

func main() {
	db := utils.InitialiseDB()
	defer db.Close()
}
