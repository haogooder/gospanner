
package mongo

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func initClient(col string) *MgClient {
	cfg := Config{
		Uri:      "mongodb://localhost:27017",
		Database: "mgotest",
		Coll:     col,
	}
	var cTimeout int64 = 0
	var sTimeout int64 = 500000
	var maxPoolSize uint64 = 30000
	var minPoolSize uint64 = 0
	cfg.ConnectTimeoutMS = &cTimeout
	cfg.SocketTimeoutMS = &sTimeout
	cfg.MaxPoolSize = &maxPoolSize
	cfg.MinPoolSize = &minPoolSize
	cfg.ReadPreference = &ReadPref{Mode: readpref.PrimaryMode}
	qClient, err := Open(context.Background(), &cfg)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return qClient
}

func TestQmgoClient(t *testing.T) {
	ast := require.New(t)
	var timeout int64 = 50

	// uri 错误
	cfg := Config{
		Uri:              "://127.0.0.1",
		ConnectTimeoutMS: &timeout,
	}

	var err error
	_, err = Open(context.Background(), &cfg)
	ast.NotNil(err)

	// Open 成功
	var maxPoolSize uint64 = 100
	var minPoolSize uint64 = 0

	cfg = Config{
		Uri:              "mongodb://localhost:27017",
		Database:         "mgotest",
		Coll:             "testopen",
		ConnectTimeoutMS: &timeout,
		MaxPoolSize:      &maxPoolSize,
		MinPoolSize:      &minPoolSize,
		ReadPreference:   &ReadPref{Mode: readpref.SecondaryMode, MaxStalenessMS: 500},
	}

	cli, err := Open(context.Background(), &cfg)
	ast.NoError(err)
	ast.Equal(cli.GetDatabaseName(), "mgotest")
	ast.Equal(cli.GetCollectionName(), "testopen")

	err = cli.Ping(5)
	ast.NoError(err)

	res, err := cli.InsertOne(context.Background(), bson.D{{Key: "x", Value: 1}})
	ast.NoError(err)
	ast.NotNil(res)

	cli.DropCollection(context.Background())

	// close Client
	cli.Close(context.TODO())
	_, err = cli.InsertOne(context.Background(), bson.D{{Key: "x", Value: 1}})
	ast.EqualError(err, "client is disconnected")

	err = cli.Ping(5)
	ast.Error(err)

	// primary mode with max stalenessMS, error
	cfg = Config{
		Uri:              "mongodb://localhost:27017",
		Database:         "mgotest",
		Coll:             "testopen",
		ConnectTimeoutMS: &timeout,
		MaxPoolSize:      &maxPoolSize,
		ReadPreference:   &ReadPref{Mode: readpref.PrimaryMode, MaxStalenessMS: 500},
	}

	cli, err = Open(context.Background(), &cfg)
	ast.Error(err)
}

func TestClient(t *testing.T) {
	ast := require.New(t)

	var maxPoolSize uint64 = 100
	var minPoolSize uint64 = 0
	var timeout int64 = 50

	cfg := &Config{
		Uri:              "mongodb://localhost:27017",
		ConnectTimeoutMS: &timeout,
		MaxPoolSize:      &maxPoolSize,
		MinPoolSize:      &minPoolSize,
	}

	c, err := NewClient(context.Background(), cfg)
	ast.Equal(nil, err)
	coll := c.Database("mgotest").Collection("testopen")

	res, err := coll.InsertOne(context.Background(), bson.D{{Key: "x", Value: 1}})
	ast.NoError(err)
	ast.NotNil(res)
	coll.DropCollection(context.Background())
}

func TestClient_ServerVersion(t *testing.T) {
	ast := require.New(t)

	cfg := &Config{
		Uri:      "mongodb://localhost:27017",
		Database: "mgotest",
		Coll:     "transaction",
	}

	ctx := context.Background()
	cli, err := Open(ctx, cfg)
	ast.NoError(err)

	version := cli.ServerVersion()
	ast.NotEmpty(version)
	fmt.Println(version)
}

func TestClient_newAuth(t *testing.T) {
	ast := require.New(t)

	auth := Credential{
		AuthMechanism: "PLAIN",
		AuthSource:    "PLAIN",
		Username:      "mgo",
		Password:      "123",
		PasswordSet:   false,
	}
	cred, err := newAuth(auth)
	ast.NoError(err)
	ast.Equal(auth.PasswordSet, cred.PasswordSet)
	ast.Equal(auth.AuthSource, cred.AuthSource)
	ast.Equal(auth.AuthMechanism, cred.AuthMechanism)
	ast.Equal(auth.Username, cred.Username)
	ast.Equal(auth.Password, cred.Password)

	auth = Credential{
		AuthMechanism: "PLAIN",
		AuthSource:    "PLAIN",
		Username:      "qmg/o",
		Password:      "123",
		PasswordSet:   false,
	}
	_, err = newAuth(auth)
	ast.Equal(ErrNotSupportedUsername, err)

	auth = Credential{
		AuthMechanism: "PLAIN",
		AuthSource:    "PLAIN",
		Username:      "mgo",
		Password:      "12:3",
		PasswordSet:   false,
	}
	_, err = newAuth(auth)
	ast.Equal(ErrNotSupportedPassword, err)

	auth = Credential{
		AuthMechanism: "PLAIN",
		AuthSource:    "PLAIN",
		Username:      "mgo",
		Password:      "1/23",
		PasswordSet:   false,
	}
	_, err = newAuth(auth)
	ast.Equal(ErrNotSupportedPassword, err)

	auth = Credential{
		AuthMechanism: "PLAIN",
		AuthSource:    "PLAIN",
		Username:      "mgo",
		Password:      "1%3",
		PasswordSet:   false,
	}
	_, err = newAuth(auth)
	ast.Equal(ErrNotSupportedPassword, err)

	auth = Credential{
		AuthMechanism: "PLAIN",
		AuthSource:    "PLAIN",
		Username:      "q%3mgo",
		Password:      "13",
		PasswordSet:   false,
	}
	_, err = newAuth(auth)
	ast.Equal(ErrNotSupportedUsername, err)
}
