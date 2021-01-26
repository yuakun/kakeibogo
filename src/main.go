package main

import (
	// Gin

	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Inproduct DB上のテーブル、カラムと構造体との関連付けが自動的に行われる
type Inproduct struct {
	ID          string `gorm:"type:varchar(2);not null"`
	ProductName string `gorm:"type:varchar(200);not null"`
	Amount      string `gorm:"type:varchar(400)"`
	//Date        time.Time `gorm:"type:datetime(6)"`
}

// Exproduct DB上のテーブル、カラムと構造体との関連付けが自動的に行われる
type Exproduct struct {
	ID          string `gorm:"type:varchar(2);not null"`
	ProductName string `gorm:"type:varchar(200);not null"`
	Amount      string `gorm:"type:varchar(400)"`
	//Date        time.Time `gorm:"type:datetime(6)"`
}

func getGormConnect() *gorm.DB {
	DBMS := "mysql"
	USER := "root"
	PASS := "root"
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

// 商品テーブルのレコードを全件取得
func findAllProduct() []Inproduct {
	db := getGormConnect()
	var products []Inproduct

	// select文
	db.Order("ID asc").Find(&products)
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

func main() {
	// product テーブルにデータを運ぶための構造体を初期化
	/*
		var product = Product{
			ProductName: "テスト商品",
			Memo:        "テスト商品です",
			Status:      "01",
		}
	*/

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

	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, Worldaa")
	})

	r.POST("/addProduct", AddProduct)
	/*
		r.GET("/addProduct", func(c *gin.Context) {
			// 構造体のポインタを渡す
			insertProduct(&product)

			// Productテーブルのレコードを全件取得する
			resultProducts := findAllProduct()

			// Productテーブルのレコードを全件表示する
			for i := range resultProducts {
				fmt.Printf("index: %d, 商品ID: %d, 商品名: %s, メモ: %s, ステータス: %s\n",
					i, resultProducts[i].ID, resultProducts[i].ProductName, resultProducts[i].Memo, resultProducts[i].Status)
			}
			//c.String(200, "Hello, World, addProduct")
		})
	*/
	r.POST("/expProduct", ExpProduct)

	r.Run()
}
