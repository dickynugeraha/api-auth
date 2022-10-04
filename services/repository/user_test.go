package repository

import (
	"api-auth/domains"
	"database/sql"
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

func (s *Suite) TestUserRepository_SuccessGetUsers() {
	rows := sqlmock.NewRows([]string{"id", "name", "email", "password"}).AddRow("uuid1", "name1", "email1", "pass1").AddRow("uuid2", "name2", "email2", "pass2")
	
	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).WillReturnRows(rows)

	users, err := s.userRepository.Users()

	s.NotEmpty(users)
	s.NotNil(users)
	s.Nil(err)
}

func (s *Suite) TestUserRepository_FailGetUsers() {

	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).WillReturnError(gorm.ErrRecordNotFound)

	users, err := s.userRepository.Users()

	s.Empty(users)
	s.Nil(users)
	s.NotNil(err)
}

func (s *Suite) TestUserRepository_SuccessFindByEmail() {
	email := "kale@gmail.com"
	query := "SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1"
	rows := sqlmock.NewRows([]string{"id","name","email","password"}).AddRow("uuid","kale",email,"password_hashing")

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnRows(rows)

	user := s.userRepository.FindByEmail(email)

	s.NotNil(user)
	s.Equal(email, user.Email)
}

func (s *Suite) TestUserRepository_FailFindByEmail() {
	email := "user_not_found@gmail.com"
	query := "SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1"

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnError(gorm.ErrRecordNotFound)

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

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("invalid_uuid").WillReturnError(gorm.ErrRecordNotFound)

	user := s.userRepository.FindById("invalid_uuid")

	s.Empty(user)
}

func (s *Suite) TestUserRepository_SuccessCreateUser() {
	db, mockTemp, _ := sqlmock.New()
	dbase, _ := gorm.Open("mysql", db)
	dbase.LogMode(true)
	repos := UserRepository{db: dbase}


	input := &domains.Register{
		Name:            "kale",
		Email:           "kale@gmail.com",
		Password:        "password",
		PasswordConfirm: "password",
	}

	query := "INSERT INTO `users` (`id`,`name`,`email`,`password`) VALUES (?,?,?,?)"

	mockTemp.ExpectBegin()
	mockTemp.ExpectExec(regexp.QuoteMeta(query)).WithArgs("uuid", input.Name, input.Email, input.Password).WillReturnResult(sqlmock.NewResult(1,1))
	mockTemp.ExpectCommit()
	
	err := repos.CreateUser("uuid", input)

	s.NoError(err)
}

func (s *Suite) TestUserRepository_FailCreateUser() {
	db, mockTemp, _ := sqlmock.New()
	dbase, _ := gorm.Open("mysql", db)
	dbase.LogMode(true)
	repos := UserRepository{db: dbase}

	input := &domains.Register{
		Name:            "kale",
		Email:           "kale@gmail.com",
		Password:        "password",
		PasswordConfirm: "password",
	}

	query := "INSERT INTO `users` (`id`,`name`,`email`,`password`) VALUES (?,?,?,?)"
 
	mockTemp.ExpectBegin()
	mockTemp.ExpectExec(regexp.QuoteMeta(query)).WithArgs("uuid",input.Name, input.Email, input.Password).WillReturnError(gorm.Errors{})
	mockTemp.ExpectCommit()

	err := repos.CreateUser("uuid", input)

	s.Error(err)
}

func (s *Suite) TestUserRepository_SuccessDeleteUser() {
	db, mockTemp, _ := sqlmock.New()
	dbase, _ := gorm.Open("mysql", db)
	dbase.LogMode(true)
	repos := UserRepository{db: dbase}

	query := "DELETE FROM `users` WHERE (id = ?)"

	mockTemp.ExpectBegin()
	mockTemp.ExpectExec(regexp.QuoteMeta(query)).WithArgs("userId").WillReturnResult(sqlmock.NewResult(1,1))
	mockTemp.ExpectCommit()

	err := repos.DeleteUserById("userId")

	s.NoError(err)
}

func (s *Suite) TestUserRepository_FailDeleteUser() {
	db, mockTemp, _ := sqlmock.New()
	dbase, _ := gorm.Open("mysql", db)
	dbase.LogMode(true)
	repos := UserRepository{db: dbase}

	query := "DELETE FROM `users` WHERE (id = ?)"

	mockTemp.ExpectBegin() 
	mockTemp.ExpectExec(regexp.QuoteMeta(query)).WithArgs("userId").WillReturnError(gorm.ErrRecordNotFound)
	mockTemp.ExpectCommit()

	err := repos.DeleteUserById("userId")

	s.Error(err)
}

func (s *Suite) TestUserRepository_SuccessChangePassword() {
	db, mockTemp, _ := sqlmock.New()
	dbase, _ := gorm.Open("mysql", db)
	dbase.LogMode(true)
	repos := UserRepository{db: dbase}

	input := &domains.ChangePassword{
		Email: "email@gmial.com",
		NewPassword: "hash_password",
		PasswordConfirm: "hash_password",
	}

	query := "UPDATE `users` SET `password` = ? WHERE (id = ?)"

	mockTemp.ExpectBegin()
	mockTemp.ExpectExec(regexp.QuoteMeta(query)).WithArgs(input.NewPassword, "uuid").WillReturnResult(sqlmock.NewResult(0,1))
	mockTemp.ExpectCommit()

	err := repos.UpdatePassword(input, "uuid")

	s.Nil(err)
	s.NoError(err)
}

func (s *Suite) TestUserRepository_FailChangePassword() {
	db, mockTemp, _ := sqlmock.New()
	dbase, _ := gorm.Open("mysql", db)
	dbase.LogMode(true)
	repos := UserRepository{db: dbase}

	input := &domains.ChangePassword{
		Email: "email@gmial.com",
		NewPassword: "hash_password",
		PasswordConfirm: "hash_password",
	}

	query := "UPDATE `users` SET `password` = ? WHERE (id = ?)"

	mockTemp.ExpectBegin()
	mockTemp.ExpectExec(regexp.QuoteMeta(query)).WithArgs("hash_password", "uuid").WillReturnError(gorm.Errors{})
	mockTemp.ExpectCommit()

	err := repos.UpdatePassword(input, "uuid")

	s.NotNil(err)
	s.Error(err)
}

func TestSuiteRepository(t *testing.T) {
	suite.Run(t, new(Suite))
}
