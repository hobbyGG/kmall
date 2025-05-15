package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/hobbyGG/kmall/review-service/internal/conf"
	"github.com/hobbyGG/kmall/review-service/internal/data/query"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDB, NewData, NewReviewRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	Q *query.Query
}

// NewData .
func NewData(db *gorm.DB, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 将数据库连接注册到gen生成的代码中
	query.SetDefault(db)

	// 将全局查询结构体传入data中，这样data就具备了查询的功能
	return &Data{Q: query.Q}, cleanup, nil
}

func NewDB(c *conf.Data) (*gorm.DB, error) {
	dsn := c.Database.Source
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	return db, nil
}
