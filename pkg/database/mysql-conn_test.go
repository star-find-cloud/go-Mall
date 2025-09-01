package database

import (
	"fmt"
	"testing"
)

func TestMySQL_GetDB(t *testing.T) {
	var mysql, err = NewMySQL()

	fmt.Println(mysql.GetDB(), err)
}
