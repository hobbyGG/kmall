package data

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"github.com/hobbyGG/kmall/review-service/internal/conf"
	"github.com/hobbyGG/kmall/review-service/internal/data/query"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewDB, NewData, NewESClient, NewReviewRepo)

// Data .
type Data struct {
	// TODO wrapped database client
	Q     *query.Query
	ESCli *elasticsearch.TypedClient
}

// NewData .
func NewData(db *gorm.DB, es *elasticsearch.TypedClient, logger log.Logger) (*Data, func(), error) {
	cleanup := func() {
		log.NewHelper(logger).Info("closing the data resources")
	}
	// 将数据库连接注册到gen生成的代码中
	query.SetDefault(db)

	// 将全局查询结构体传入data中，这样data就具备了查询的功能
	return &Data{Q: query.Q, ESCli: es}, cleanup, nil
}

func NewDB(c *conf.Data) (*gorm.DB, error) {
	dsn := c.Database.Source
	db, err := gorm.Open(mysql.Open(dsn))
	if err != nil {
		panic(err)
	}
	return db, nil
}

func NewESClient(c *conf.Elasticsearch) *elasticsearch.TypedClient {
	conf := elasticsearch.Config{
		Addresses: c.Address,
	}
	esClient, err := elasticsearch.NewTypedClient(conf)
	if err != nil {
		panic(err)
	}
	return esClient
}
