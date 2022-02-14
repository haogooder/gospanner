
package mongo

import (
	"errors"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

func TestIsErrNoDocuments(t *testing.T) {
	ast := require.New(t)
	ast.False(IsErrNoDocuments(errors.New("dont match")))
	ast.True(IsErrNoDocuments(ErrNoSuchDocuments))
	ast.True(IsErrNoDocuments(mongo.ErrNoDocuments))
}

func TestIsDup(t *testing.T) {
	ast := require.New(t)
	ast.False(IsDup(nil))
	ast.False(IsDup(errors.New("invaliderror")))
	ast.True(IsDup(errors.New("E11000")))
}
