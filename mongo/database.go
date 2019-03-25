package mongo

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"time"
)

//dsn mongodb://root:WiaQ82n7B3L5Cz*2#10m@172.18.76.150:27017?authSource=admin
type Config struct {
	Host     string `json:"host" yaml:"host"`
	Port     string `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Auth     string `json:"auth" yaml:"auth"`
	DBName   string `json:"db_name" yaml:"db_name"`
	PoolSize int    `json:"pool_size" yaml:"pool_size"`
}

//mongodb orm use official driver
type Database struct {
	client        *mongo.Client
	currentDBName string
}

func NewDatabase(conf *Config) (*Database, error) {
	opts := options.Client()
	opts.SetHosts([]string{fmt.Sprintf("%s:%s", conf.Host, conf.Port)})
	opts.SetMaxPoolSize(uint16(conf.PoolSize))
	opts.SetConnectTimeout(20 * time.Second)
	if conf.Auth != "" {
		opts.SetAuth(options.Credential{
			Username:   conf.User,
			Password:   conf.Password,
			AuthSource: conf.Auth,
		})
	}

	client, err := mongo.NewClientWithOptions("", opts)
	if err != nil {
		return nil, err
	}

	return &Database{
		client:        client,
		currentDBName: conf.DBName,
	}, nil
}

//connect to mongodb
func (d *Database) Open() error {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	return d.client.Connect(ctx)
}

//disconnect mongodb
func (d *Database) Close() error {
	return d.client.Disconnect(context.Background())
}

//select database
func (d *Database) SelectDatabase(dbName string) {
	d.currentDBName = dbName
}

//get default database
func (d *Database) DefaultDB() {
	d.client.Database(d.currentDBName)
}

//get collection
func (d *Database) Collection(collectionName string) {
	d.client.Database(d.currentDBName).Collection(collectionName)
}
