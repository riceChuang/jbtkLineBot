package boltdb

import (
	"github.com/boltdb/bolt"
	"fmt"
)


var	(_db  *Boltdb)



func Initialize(){
	//创建bolt数据库本地文件
	dbc, err := bolt.Open(dbname, 0600, nil)

	//初始化bucket
	bucket = []byte("imageBucket")
	if err != nil {
		fmt.Println("open err:", err)
		return
	} else {
		db = dbc
	}
	//创建bucket
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket(bucket)
		return err
	})

	_db = Newdb()
}


func DB() *Boltdb {
	return _db
}

