package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	models2 "learn/src/go_learn/go_database/go_mongdb/models"
	"log"
)

type User struct {
	Name      string
	Password  string
	Email     string
	Active    bool
	Authority int32
}

type timerArchive struct {
	SpotTime       string
	SpotDate       string
	SpotDataID     string
	SpotTrackingID string
	GroupName      string
	ModuleName     string
	WorkpieceNum   string
	ControlMode    string //1=Skt, 2=KSR, 4=IQR, 5=AMC/DCM, 9=AMF

	typeID     string
	SpotNum    string
	ProgramNum string
	spotName   string

	basicCurrent string
	preTime      string
	weldTime     string
	keepTime     string

	weldMessage string

	gunNum  string
	gunName string

	qActive     string
	qSpotValue  string
	spatterRate string

	RKV                string
	IKV                string
	UKV                string
	ControlStrokeCurve string
}

type timerErrMsg struct {
	//ModuleErrorText
	errCode int32
	Message struct {
		//00 = Info, 01 = cause, 02 = solution, 03 = Ausloeser
		Info     string
		Cause    string
		Solution string
	}
}

func Find() {

	client := models2.DB.Mongo
	collection, _ := client.Database("timer").Collection("users").Clone()
	//collection.

	findOptions := options.Find()
	findOptions.SetLimit(5)
	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cur)

	var results []*User

	for cur.Next(context.TODO()) {
		//定义一个文档，将单个文档解码为result
		var result User
		err := cur.Decode(&result)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &result)
		fmt.Println(result)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	//遍历结束后关闭游标
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)
}

func getData() {
	client := models2.DB.Mongo
	db := client.Database("timer")
	col := db.Collection("archive")
	ret := col.FindOne(context.TODO(), map[string]string{"WorkpieceID": "SNG000000000000000001"})
	raw, _ := ret.DecodeBytes()
	fmt.Println(raw)
}
func getData1() {
	client := models2.DB.Mongo
	db := client.Database("timer")
	col := db.Collection("msg")
	ret := col.FindOne(context.TODO(), map[string]int32{"errCode": 400001})
	raw, _ := ret.DecodeBytes()
	fmt.Println(raw)
}

func main() {
	models2.Init() // init mongodb

	//client := models.DB.Mongo
	//collection, _ := client.Database("sng").Collection("user").Clone()
	//abb := User{"abb", 24}
	//insertOne, err := collection.InsertOne(context.TODO(), abb)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Inserted a Single Document: ", insertOne.InsertedID)
	Find()
	getData()
	getData1()
}
