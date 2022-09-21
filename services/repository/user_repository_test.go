package repository

import (
	"api-auth/domains"
	"database/sql"
	"fmt"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DB   *gorm.DB
	mock sqlmock.Sqlmock
	// user           User
	userRepository UserRepository
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	assert.NoError(s.T(), err)

	s.DB, err = gorm.Open("mysql", db)
	assert.NoError(s.T(), err)

	s.DB.LogMode(true)

	s.userRepository = UserRepository{db: s.DB}
}
func (s *Suite) TestUserRepository_SuccessFindByEmail() {
	email := "kale@gmail.com"
	query := "SELECT * FROM `users` WHERE (email = ?)"

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password"}).AddRow("uuid", "kale", email, "password_hashing")

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnRows(rows)

	user := s.userRepository.FindByEmail(email)

	s.NotEmpty(user)
}

func (s *Suite) TestUserRepository_FailFindByEmail() {
	email := "user_not_found@gmail.com"
	query := "SELECT * FROM `users` WHERE (email = ?)"

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnError(fmt.Errorf("email not found!"))

	user := s.userRepository.FindByEmail(email)

	s.Nil(user)
	s.Empty(user)
}

func (s *Suite) TestUserRepository_SuccessFindById() {
	user1 := &User{
		ID:       "valid_uuid",
		Name:     "kale",
		Email:    "kale@gmail.com",
		Password: "password_hashing",
	}

	query := "SELECT * FROM `users` WHERE (id = ?)"

	rows := sqlmock.NewRows([]string{"id", "name", "email", "password"}).AddRow("valid_uuid", "kale", "kale@gmail.com", "password_hashing")

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("valid_uuid").WillReturnRows(rows)

	user := s.userRepository.FindById("valid_uuid")

	s.NotEmpty(user)
	s.Equal(user1.ID, user.ID)

}

func (s *Suite) TestUserRepository_FailFindById() {
	query := "SELECT * FROM `users` WHERE (id = ?)"

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("invalid_uuid").WillReturnError(fmt.Errorf("User not found!"))

	user := s.userRepository.FindById("invalid_uuid")

	s.Empty(user)
}

func (s *Suite) TestUserRepository_SuccessCreateUser() {
	input := &domains.Register{
		Name:            "kale",
		Email:           "kale@gmail.com",
		Password:        "password",
		PasswordConfirm: "password",
	}

	query := "INSERT INTO `users` (`id`,`name`,`email`,`password`) VALUES ($1, $2, $3, $4) RETURNING `users`.`id`"
	rows := sqlmock.NewRows([]string{"id"}).AddRow("uuid")
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("uuid", input.Name, input.Email, input.Password).WillReturnRows(rows)

	err := s.userRepository.CreateUser(input)

	s.Nil(err)
}

func TestSuiteRepository(t *testing.T) {
	suite.Run(t, new(Suite))
}
