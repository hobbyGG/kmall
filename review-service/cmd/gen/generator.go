package main

import (
	"flag"

	"github.com/hobbyGG/kmall/review-service/internal/conf"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var flagconf string

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	c := config.New(
		config.WithSource(
			file.NewSource(flagconf),
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}
	g := gen.NewGenerator(gen.Config{
		// CURD代码存放位置
		OutPath: "../../internal/data/query",

		// 生成全局query，生成query接口
		Mode: gen.WithDefaultQuery | gen.WithQueryInterface,

		FieldNullable: true,
	})

	g.UseDB(NewDB(bc.Data))

	// 从连接的数据库为所有表生成Model结构体和CRUD代码
	// 也可以手动指定需要生成代码的数据表
	g.ApplyBasic(g.GenerateAllTable()...)

	// 自定义需要实现的功能，定义一个接口，写好对应方法，在方法上方写上对应sql语句的备注，gen会自动生成对应的方法
	// g.ApplyInterface(func() {}, models...)

	g.Execute()
}

func NewDB(c *conf.Data) *gorm.DB {
	db, err := gorm.Open(mysql.Open(c.Database.Source))
	if err != nil {
		panic(err)
	}
	return db
}
