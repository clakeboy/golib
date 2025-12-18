package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// dsn mongodb://root:WiaQ82n7B3L5Cz*2#10m@172.18.76.150:27017?authSource=admin
type Config struct {
	Host        string `json:"host" yaml:"host"`
	Port        string `json:"port" yaml:"port"`
	User        string `json:"user" yaml:"user"`
	Password    string `json:"password" yaml:"password"`
	Auth        string `json:"auth" yaml:"auth"`
	DBName      string `json:"db_name" yaml:"db_name"`
	PoolSize    uint64 `json:"pool_size" yaml:"pool_size"`
	Timeout     uint64 `json:"timeout" yaml:"timeout"`
	Compression bool   `json:"compression" yaml:"compression"`
}

// mongodb orm use official driver
type Database struct {
	client        *mongo.Client
	currentDBName string
}

func NewDatabase(conf *Config) (*Database, error) {
	opts := options.Client()
	opts.SetHosts([]string{fmt.Sprintf("%s:%s", conf.Host, conf.Port)})
	opts.SetMaxPoolSize(uint64(conf.PoolSize))
	opts.SetMinPoolSize(1)
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetTimeout(30 * time.Second)
	if conf.Auth != "" {
		opts.SetAuth(options.Credential{
			AuthMechanism: "SCRAM-SHA-1",
			AuthSource:    conf.Auth,
			Username:      conf.User,
			Password:      conf.Password,
			PasswordSet:   true,
		})
	}
	if conf.Compression {
		opts.SetCompressors([]string{"snappy", "zlib", "zstd"})
	}
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	return &Database{
		client:        client,
		currentDBName: conf.DBName,
	}, nil
}

func NewDatabaseDsn(dsn string) (*Database, error) {
	opts := options.Client()
	opts.ApplyURI(dsn)
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}

	return &Database{
		client:        client,
		currentDBName: "",
	}, nil
}

// connect to mongodb
func (d *Database) Open() error {
	return d.Ping()
}

// disconnect mongodb
func (d *Database) Close() error {
	return d.client.Disconnect(context.Background())
}

// select database
func (d *Database) SelectDatabase(dbName string) {
	d.currentDBName = dbName
}

func (d *Database) ListDatabase() (mongo.ListDatabasesResult, error) {
	return d.client.ListDatabases(context.Background(), bson.D{})
}

func (d *Database) ListDatabaseNames() ([]string, error) {
	return d.client.ListDatabaseNames(context.Background(), bson.D{})
}

// get default database
func (d *Database) Database(dbName ...string) *mongo.Database {
	if len(dbName) > 0 {
		return d.client.Database(dbName[0])
	}
	return d.client.Database(d.currentDBName)
}

// ping mongodb server
func (d *Database) Ping() error {
	ctx, cencel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cencel()
	return d.client.Ping(ctx, readpref.Primary())
}

// get mongodb host info
func (d *Database) HostInfo() (bson.M, error) {
	ctx, cencel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cencel()
	var hostInfo bson.M
	err := d.client.Database("test").RunCommand(ctx, bson.D{{Key: "hostInfo", Value: 1}}).Decode(&hostInfo)
	return hostInfo, err
}

// get mongodb build info
func (d *Database) BuildInfo() (bson.M, error) {
	ctx, cencel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cencel()
	var hostInfo bson.M
	err := d.client.Database("test").RunCommand(ctx, bson.D{{Key: "buildInfo", Value: 1}}).Decode(&hostInfo)
	return hostInfo, err
}

// get mongodb database status
func (d *Database) DbStats(dbname string) (bson.M, error) {
	ctx, cencel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cencel()
	var hostInfo bson.M
	err := d.client.Database(dbname).RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}}).Decode(&hostInfo)
	return hostInfo, err
}
