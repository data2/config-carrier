package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sync"
)

type Response struct {
	Action  string
	Code    int
	Message string
	Data    interface{}
}

type Record struct {
	ID        string `gorm:"primaryKey"`
	Namespace string
	Path      string
	Value     string
}

const (
	SUCCESS = 0
	FAIL    = 1
	UNAUTH  = 2
)

func DelData(db *gorm.DB, namespace string, path string) Response {
	err := db.Delete(&Record{}, "namespace =? and path = ?", namespace, path).Error
	if err != nil {
		return Response{Code: FAIL, Message: "del fail."}
	}
	return Response{Code: SUCCESS}
}

var lock = sync.Mutex{}

func SetData(db *gorm.DB, namespace string, path string, value string) Response {
	lock.Lock()
	record := LoadData(db, namespace, path)
	var err error
	if record.Code == SUCCESS && record.Data != (Record{}) {
		err = db.Model(&(record.Data)).Update(path, value).Error
	} else {
		uid, _ := uuid.NewUUID()
		err = db.Create(&Record{Namespace: namespace, ID: uid.String(), Path: path, Value: value}).Error
	}
	lock.Unlock()
	if err != nil {
		return Response{Code: FAIL, Message: "set data fail."}
	}
	return Response{Code: SUCCESS}
}

func LoadData(db *gorm.DB, namespace string, path string) Response {
	var record Record
	err := db.Where("namespace = ? and path = ?", namespace, path).Find(&record).Error
	if err != nil {
		return Response{Code: FAIL, Message: err.Error()}
	}
	return Response{Code: SUCCESS, Data: record}
}

func DelAllData(db *gorm.DB, namespace string) Response {
	err := db.Where(" namespace = ?", namespace).Delete(Record{}).Error
	if err != nil {
		return Response{Code: FAIL, Message: "del all data fail."}
	}
	return Response{Code: SUCCESS}
}

func LoadAllData(db *gorm.DB, namespace string) Response {
	var records []Record
	err := db.Where("namespace=?", namespace).Find(&records).Error
	if err != nil {
		return Response{Code: FAIL, Message: "load all data fail."}
	}
	return Response{Code: SUCCESS, Data: records}
}
