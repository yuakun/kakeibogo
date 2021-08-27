package main

import (
	// Gin

	"fmt"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Inproduct DB上のテーブル、カラムと構造体との関連付けが自動的に行われる
type Inproduct struct {
	ID          string `gorm:"type:varchar(11);not null"`
	ProductName string `gorm:"type:varchar(200);not null"`
	Amount      string `gorm:"type:varchar(400)"`
	//Date        time.Time `gorm:"type:datetime(6)"`
}

// Exproduct DB上のテーブル、カラムと構造体との関連付けが自動的に行われる
type Exproduct struct {
	ID          string `gorm:"type:varchar(11);not null"`
	ProductName string `gorm:"type:varchar(200);not null"`
	Amount      string `gorm:"type:varchar(400)"`
	//Date        time.Time `gorm:"type:datetime(6)"`
}

func getGormConnect() *gorm.DB {
	DBMS := "mysql"
	USER := "root"
	PASS := "XXXXX"
	PROTOCOL := "tcp(localhost:3306)"
	DBNAME := "Shopping"
	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME
	db, err := gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}

	// DBエンジンを「InnoDB」に設定
	db.Set("gorm:table_options", "ENGINE=InnoDB")

	// 詳細なログを表示
	db.LogMode(true)

	// 登録するテーブル名を単数形にする（デフォルトは複数形）
	db.SingularTable(true)

	// マイグレーション（テーブルが無い時は自動生成）
	db.AutoMigrate(&Inproduct{}, &Exproduct{})

	fmt.Println("db connected: ", &db)
	return db
}

// 商品テーブルにレコードを追加
func insertProduct(registerProduct *Inproduct) {
	db := getGormConnect()

	// insert文
	db.Table("inproduct").Create(&registerProduct)
	defer db.Close()
}

func insertExProduct(registerProduct *Exproduct) {
	db := getGormConnect()

	// insert文
	db.Table("exproduct").Create(&registerProduct)
	defer db.Close()
}

// FetchAllProductsIncome 商品テーブルのレコードを全件取得
func FetchAllProductsIncome(c *gin.Context) {
	resultProducts := FindAllProductsIncome()

	// URLへのアクセスに対してJSONを返す
	c.JSON(200, resultProducts)
}

// FindAllProductsIncome 商品テーブルのレコードを全件取得
func FindAllProductsIncome() []Inproduct {
	db := getGormConnect()
	var products []Inproduct

	// select文
	db.Table("inproduct").Order("ID asc").Find(&products)
	defer db.Close()
	return products
}

// FetchAllProductsExpense 商品テーブルのレコードを全件取得
func FetchAllProductsExpense(c *gin.Context) {
	resultProducts := FindAllProductsExpense()

	// URLへのアクセスに対してJSONを返す
	c.JSON(200, resultProducts)
}

// FindAllProductsExpense 商品テーブルのレコードを全件取得
func FindAllProductsExpense() []Inproduct {
	db := getGormConnect()
	var products []Inproduct

	// select文
	db.Table("exproduct").Order("ID asc").Find(&products)
	defer db.Close()
	return products
}

// 商品の購入状態を定義
const (
	// NotPurchased は 未購入
	NotPurchased = 0

	// Purchased は 購入済
	Purchased = 1
)

// AddProduct は 商品をDBへ登録する
func AddProduct(c *gin.Context) {
	id := c.PostForm("ID")
	productName := c.PostForm("ProductName")
	amount := c.PostForm("Amount")
	//	date := c.PostForm("Date")

	var product = Inproduct{
		ID:          id,
		ProductName: productName,
		Amount:      amount,
		//		Date:        date,
	}

	insertProduct(&product)
}

// ExpProduct は 商品をDBへ登録する
func ExpProduct(c *gin.Context) {
	id := c.PostForm("ID")
	productName := c.PostForm("ProductName")
	amount := c.PostForm("Amount")

	var product = Exproduct{
		ID:          id,
		ProductName: productName,
		Amount:      amount,
	}

	insertExProduct(&product)
}

// DbDeleteProduct は 商品テーブルの指定したレコードを削除する
func DbDeleteProduct(productID int) {
	/*
	product := []entity.Product{}

	db := getGormConnect()
	// delete
	db.Table("inproduct").Delete(&product, productID)
	defer db.Close()
	*/
}

// inComeDelete は 商品情報をDBから削除する
func inComeDelete(c *gin.Context) {
	productIDStr := c.PostForm("ID")

	productID, _ := strconv.Atoi(productIDStr)

	DbDeleteProduct(productID)
}

func main() {

	r := gin.Default()

	// ここからCorsの設定
	r.Use(cors.New(cors.Config{
		// アクセスを許可したいアクセス元
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		// アクセスを許可したいHTTPメソッド(以下の例だとPUTやDELETEはアクセスできません)
		AllowMethods: []string{
			"POST",
			"GET",
			"OPTIONS",
		},
		// 許可したいHTTPリクエストヘッダ
		AllowHeaders: []string{
			"Access-Control-Allow-Origin",
			"Access-Control-Allow-Credentials",
			"Access-Control-Allow-Headers",
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"Authorization",
		},
		// cookieなどの情報を必要とするかどうか
		AllowCredentials: true,
		// preflightリクエストの結果をキャッシュする時間
		MaxAge: 24 * time.Hour,
	}))

	r.POST("/addProduct", AddProduct)

	r.POST("/expProduct", ExpProduct)

	r.GET("/getIncome", FetchAllProductsIncome)

	r.GET("/getExpense", FetchAllProductsExpense)

	r.POST("/delete", inComeDelete)

	r.Run()
}
