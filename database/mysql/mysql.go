package mysql

import (
	"fmt"
	"worker-sample/config"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(ctx *config.ServiceContext) *gorm.DB {
	var datetimePrecision = 2
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       genDSN(ctx.Config), // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
		DefaultStringSize:         256,                // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
		DisableDatetimePrecision:  true,               // disable datetime precision support, which not supported before MySQL 5.6
		DefaultDatetimePrecision:  &datetimePrecision, // default datetime precision
		DontSupportRenameIndex:    true,               // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,               // use change when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,              // smart configure based on used version
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Init MySQL DB failed: %+v", err)
	}
	ctx.DB = db
	return db
}

func CloseDB(ctx *config.ServiceContext) {
	sqldb, err := ctx.DB.DB()
	if err != nil {
		log.Warnf("Get sql DB failed: %+v", err)
		return
	}
	if err := sqldb.Close(); err != nil {
		log.Warnf("Close sql DB failed: %+v", err)
	}
}

// genDSN generate DSN string
func genDSN(c *config.Config) string {
	result := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.MySQL.Username,
		c.MySQL.Password,
		c.MySQL.Host,
		c.MySQL.Port,
		c.MySQL.Schemas)
	return result
}
