package mysql

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_delivered(t *testing.T) {
	db, err := initDBConn()
	defer db.Close()
	assert.Nil(t, err)
	ctx := context.Background()
	err = delivered(ctx, db, "9f27f7f2-96c9-46dc-9814-5cc2943a3808")
	assert.Nil(t, err)
}

func Test_failed(t *testing.T) {
	db, err := initDBConn()
	defer db.Close()
	assert.Nil(t, err)
	ctx := context.Background()
	err = failed(ctx, db, "9f27f7f2-96c9-46dc-9814-5cc2943a3808")
	assert.Nil(t, err)
}
