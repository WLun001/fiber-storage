// +build mysql

package mysql

import (
	"os"
	"testing"
	"time"

	"github.com/gofiber/utils"
)

var storeConfig = Config{
	DropTable: true,
}

func init() {
	if v := os.Getenv("MYSQL_ADDRESS"); v != "" {
		storeConfig.Address = v
	}
	if v := os.Getenv("MYSQL_USERNAME"); v != "" {
		storeConfig.Username = v
	}
	if v := os.Getenv("MYSQL_PASSWORD"); v != "" {
		storeConfig.Password = v
	}
	if v := os.Getenv("MYSQL_DATABASE"); v != "" {
		storeConfig.DatabaseName = v
	}
}

func Test_MySQL_Set(t *testing.T) {

	store := New(storeConfig)

	id := "hello"
	value := []byte("Hi there!")

	store.Set(id, value, 0)

	var (
		returnedValue []byte
		exp           int64
	)

	store.db.QueryRow(store.sqlSelect, id).Scan(&returnedValue, &exp)

	utils.AssertEqual(t, returnedValue, value)
	utils.AssertEqual(t, exp, int64(0))

}

func Test_MySQL_SetExpiry(t *testing.T) {

	store := New(storeConfig)

	id := "hello"
	value := []byte("Hi there!")
	expiry := time.Second * 10

	store.Set(id, value, expiry)

	now := time.Now().Unix()
	var (
		returnedValue []byte
		exp           int64
	)
	store.db.QueryRow(store.sqlSelect, id).Scan(&returnedValue, &exp)

	delta := exp - now
	upperBound := int64(expiry.Seconds())
	lowerBound := upperBound - 2

	if !(delta <= upperBound && delta > lowerBound) {
		t.Fatalf("Test_SetExpiry: expiry delta out of bounds (is %d, must be %d<x<=%d)", delta, lowerBound, upperBound)
	}

}

func Test_MySQL_Get(t *testing.T) {

	store := New(storeConfig)

	id := "hello"
	value := []byte("Hi there!")

	store.db.Exec(store.sqlInsert, id, utils.UnsafeString(value), 0)

	returnedValue, err := store.Get(id)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, value, returnedValue)

}

func Test_MySQL_Get_NoRows(t *testing.T) {

	store := New(storeConfig)

	id := "hello"

	returnedValue, err := store.Get(id)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 0, len(returnedValue))

}

func Test_MySQL_Get_Expired(t *testing.T) {

	store := New(storeConfig)

	id := "hello"
	value := []byte("Hi there!")

	store.db.Exec(store.sqlInsert, id, utils.UnsafeString(value), time.Now().Add(time.Minute*-1).Unix())

	returnedValue, err := store.Get(id)
	utils.AssertEqual(t, nil, err)
	utils.AssertEqual(t, 0, len(returnedValue))
}

func Test_MySQL_Delete(t *testing.T) {

	store := New(storeConfig)

	id := "hello"
	value := []byte("Hi there!")

	store.db.Exec(store.sqlInsert, id, utils.UnsafeString(value), 0)

	err := store.Delete(id)
	utils.AssertEqual(t, nil, err)

	row := store.db.QueryRow(store.sqlSelect, id)
	err = row.Scan()
	utils.AssertEqual(t, noRows, err.Error())

}

func Test_MySQL_Clear(t *testing.T) {

	store := New(storeConfig)

	id := "hello"
	value := []byte("Hi there!")

	store.db.Exec(store.sqlInsert, id, utils.UnsafeString(value), 0)

	err := store.Clear()
	utils.AssertEqual(t, nil, err)

	row := store.db.QueryRow(store.sqlSelect, id)
	err = row.Scan()
	utils.AssertEqual(t, noRows, err.Error())

}

func Benchmark_MySQL_Set(b *testing.B) {
	store := New(storeConfig)

	key := "aaaa"
	val := []byte("This is a value")

	expiry := time.Second * 60

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		store.Set(key, val, expiry)
	}

}

func Benchmark_MySQL_Get(b *testing.B) {
	store := New(storeConfig)

	key := "aaaa"
	val := []byte("This is a value")

	store.Set(key, val, 0)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		store.Get(key)
	}

}
