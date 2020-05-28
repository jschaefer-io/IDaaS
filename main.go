package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/jschaefer-io/IDaaS/action"
)

func main() {

	// https://github.com/dgrijalva/jwt-go
	// https://golang.org/pkg/net/smtp/

	//fmt.Println(len(crypto.NewToken()))
	//fmt.Println(crypto.NewToken())

	//db, err := gorm.Open("sqlite3", "./test.db")
	//if err != nil {
	//	fmt.Println(err)
	//	panic("failed to connect database")
	//}
	//defer db.Close()
	//
	//// Migrate the schema
	//db.AutoMigrate(&model.Identity{})
	//
	//for i := 0; i < 1000; i++ {
	//	id := model.Identity{
	//		Email:    "test@test.de",
	//		Password: fmt.Sprintf("%d", i),
	//		Salt:     "Lorem Ipsum",
	//	}
	//	db.Create(&id)
	//}
	//
	//
	//var idTest []model.Identity
	//db.Find(&idTest)
	//for _,i := range idTest{
	//	fmt.Println(i)
	//}
	//
	//


	r := gin.Default()
	r.POST("/", action.AddIdentity)
	r.Run()
}