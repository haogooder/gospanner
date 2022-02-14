
package mongo

import (
    "context"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
    "testing"

    opts "github.com/haogooder/gospanner/mongo/options"
    "github.com/stretchr/testify/require"
    "go.mongodb.org/mongo-driver/bson"
)

func TestDatabase(t *testing.T) {
    ast := require.New(t)

    var sTimeout int64 = 500000
    var cTimeout int64 = 3000
    var maxPoolSize uint64 = 3000
    var minPoolSize uint64 = 0
    collName := "testopen"
    dbName := "mgotest"

    cfg := Config{
        Uri:              "mongodb://localhost:27017",
        Database:         dbName,
        Coll:             collName,
        ConnectTimeoutMS: &cTimeout,
        SocketTimeoutMS:  &sTimeout,
        MaxPoolSize:      &maxPoolSize,
        MinPoolSize:      &minPoolSize,
    }

    c, err := NewClient(context.Background(), &cfg)
    ast.NoError(err)
    cli := c.Database(cfg.Database)
    ast.Nil(err)
    ast.Equal(dbName, cli.GetDatabaseName())
    coll := cli.Collection(collName)
    ast.Equal(collName, coll.GetCollectionName())
    cli.Collection(collName).DropCollection(context.Background())
    cli.DropDatabase(context.Background())

}

func TestRunCommand(t *testing.T) {
    ast := require.New(t)

    cli := initClient("test")

    opts := opts.RunCommandOptions{RunCmdOptions: options.RunCmd().SetReadPreference(readpref.Primary())}
    res := cli.RunCommand(context.Background(), bson.D{
        {"ping", 1}}, opts)
    ast.NoError(res.Err())
}

// func TestCreateCollection(t *testing.T) {
//	ast := require.New(t)
//
//	cli := initClient("test")
//
//	timeSeriesOpt := options.TimeSeriesOptions{
//		TimeField:"timestamp",
//	}
//	timeSeriesOpt.SetMetaField("metadata")
//	ctx := context.Background()
//	createCollectionOpts := opts.CreateCollectionOptions{CreateCollectionOptions: options.CreateCollection().SetTimeSeriesOptions(&timeSeriesOpt)}
//	if err := cli.CreateCollection(ctx, "syslog", createCollectionOpts); err != nil {
//		ast.NoError(err)
//	}
//	cli.DropCollection(ctx)
//	cli.DropDatabase(ctx)
// }
