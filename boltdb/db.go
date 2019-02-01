package boltdb

import (
	"fmt"

	"github.com/boltdb/bolt"
)

var (
	db     *bolt.DB
	bucket []byte
)

const dbname = "image.db"

type Boltdb struct {

}

func Newdb() *Boltdb{
	return &Boltdb{}
}

//把数据插入到bolt数据库中，相当于redis中的set命令
func (b *Boltdb) Insert(key, value string) {

	k := []byte(key)
	v := []byte(value)
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)

		err := b.Put(k, v)
		return err
	})
}

//删除一个指定的key中的数据
func (b *Boltdb) Rm(key string) {
	k := []byte(key)
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		err := b.Delete(k)
		return err
	})
}

//读取一条数据
func (b *Boltdb) Read(key string) string {
	k := []byte(key)
	var val []byte
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		val = b.Get(k)
		return nil
	})
	return string(val)
}

//遍历指定的bucket中的数据
func (b *Boltdb) FetchAll(buk []byte) {
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(buk)
		cur := b.Cursor()
		for k, v := cur.First(); k != nil; k, v = cur.Next() {
			fmt.Printf("key is %s,value is %s\n", k, v)
		}
		return nil
	})
}
