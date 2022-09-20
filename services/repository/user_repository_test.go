package repository

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type Suite struct {
	suite.Suite
	DB *gorm.DB
	mock sqlmock.Sqlmock
	user User
}

func (s *Suite) SetupSuite(t *testing.T) {
	var (
		db *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	assert.NoError(t, err)

	s.DB, err = gorm.Open("mysql", db)
	assert.NoError(t, err)

	UserRepository(s.DB)
}