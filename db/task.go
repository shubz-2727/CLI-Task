package db

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/boltdb/bolt"
)

var taskBucket = []byte("tasks")
var completeBucket = []byte("complete")
var db *bolt.DB

type Task struct {
	Key   int
	Value string
}

type Complete struct {
	Key   string
	Value string
}

func Init(dbPath string) error {
	var err error
	db, err = bolt.Open(dbPath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	return db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(taskBucket)
		if err != nil {
			return err
		}
		_, err = tx.CreateBucketIfNotExists(completeBucket)
		if err != nil {
			return err
		}
		return nil
	})
}

func CreateTask(task string) (int, error) {
	var id int
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		id64, _ := b.NextSequence()
		id = int(id64)
		key := itob(id)

		return b.Put(key, []byte(task))
	})
	if err != nil {
		return -1, err
	}
	return id, nil
}

func AllTasks() ([]Task, error) {
	var tasks []Task
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			tasks = append(tasks, Task{
				Key:   btoi(k),
				Value: string(v),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func DeleteTask(key int) error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(taskBucket)
		return b.Delete(itob(key))
	})
}

func CompleteTask(value string) error {

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(completeBucket)
		t := time.Now()

		return b.Put([]byte(t.Format(time.RFC3339)), []byte(value))
	})
}

func TodaysTask() ([]Complete, error) {

	var tasks []Complete
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(completeBucket)
		c := b.Cursor()
		//k, v := c.First()

		min := []byte(time.Now().AddDate(0, 0, -1).Format(time.RFC3339))
		max := []byte(time.Now().AddDate(0, 0, 0).Format(time.RFC3339))

		for k, v := c.Seek(min); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			//fmt.Println(string(k), string(v))
			tasks = append(tasks, Complete{
				Key:   string(k),
				Value: string(v),
			})
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return tasks, nil

}

func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func btoi(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}
