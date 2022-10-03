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

// func (s *Suite) TestUserRepository_FailGetUsers() {

// 	s.mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).WillReturnError(fmt.Errorf("cannot fetch all user"))

// 	users, err := s.userRepository.Users()

// 	s.Empty(users)
// 	s.Nil(users)
// 	s.NotNil(err)
// }

func (s *Suite) TestUserRepository_SuccessFindByEmail() {
	email := "kale@gmail.com"
	query := "SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1"
	rows := sqlmock.NewRows([]string{"id","name","email","password"}).AddRow("uuid","kale",email,"password_hashing")

	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnRows(rows)

	user := s.userRepository.FindByEmail(email)

	s.NotNil(user)
	s.Equal(email, user.Email)
}

// func (s *Suite) TestUserRepository_FailFindByEmail() {
// 	email := "user_not_found@gmail.com"
// 	query := "SELECT * FROM `users` WHERE (email = ?) ORDER BY `users`.`id` ASC LIMIT 1"

// 	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(email).WillReturnError(errors.New(""))

// 	user := s.userRepository.FindByEmail(email)

// 	s.Nil(user)
// 	s.Empty(user)
// }

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

// func (s *Suite) TestUserRepository_FailFindById() {
// 	query := "SELECT * FROM `users` WHERE (id = ?)"

// 	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("invalid_uuid").WillReturnError(errors.New(""))

// 	user := s.userRepository.FindById("invalid_uuid")

// 	s.Empty(user)
// }

func (s *Suite) TestUserRepository_SuccessCreateUser() {
	input := &domains.Register{
		Name:            "kale",
		Email:           "kale@gmail.com",
		Password:        "password",
		PasswordConfirm: "password",
	}

	query := "INSERT INTO `users` (`id`,`name`,`email`,`password`) VALUES ($1,$2,$3,$4) RETURNING `users`.`id`"
	rows := sqlmock.NewRows([]string{"id"}).AddRow("uuid")

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("uuid", input.Name, input.Email, input.Password).WillReturnRows(rows)
	s.mock.ExpectCommit()
	
	err := s.userRepository.CreateUser(input)

	s.NoError(err)
}

// func (s *Suite) TestUserRepository_FailCreateUser() {
// 	input := &domains.Register{
// 		Name:            "kale",
// 		Email:           "kale@gmail.com",
// 		Password:        "password",
// 		PasswordConfirm: "password",
// 	}

// 	query := "INSERT INTO `users` (`id`,`name`,`email`,`password`) VALUES (?,?,?,?)"
// 	// query := "INSERT INTO `users` (`id`,`name`,`email`,`password`) VALUES ($1,$2,$3,$4)"
// 	// query := `INSERT INTO "users" ("id","name","email","password"\) VALUES \(\?,\?,\?,\?\)`
 
// 	s.mock.ExpectBegin()
// 	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("uuid", input.Name, input.Email, input.Password).WillReturnError(errors.New(""))
// 	s.mock.ExpectCommit()

// 	err := s.userRepository.CreateUser(input)

// 	s.Error(err)
// }

func (s *Suite) TestUserRepository_SuccessDeleteUser() {
	// query := "DELETE FROM `users` WHERE (id = ?)"
	query := "DELETE FROM `users`  WHERE (id = ?) RETURNING `users`.`id`"
	rows := sqlmock.NewRows([]string{"id"}).AddRow("userId")

	defer s.DB.Close()
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs("userId").WillReturnRows(rows)
	// s.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs("userId").WillReturnResult(sqlmock.NewResult(0,0))
	s.mock.ExpectCommit()


	err := s.userRepository.DeleteUserById("userId")

	s.NoError(err)
}

// func (s *Suite) TestUserRepository_FailDeleteUser() {
// 	query := "DELETE FROM `users` WHERE (id = ?)"

// 	s.mock.ExpectBegin()
// 	s.mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs("userId").WillReturnError(errors.New(""))
// 	s.mock.ExpectCommit()

// 	err := s.userRepository.DeleteUserById("userId")

// 	s.Error(err)
// }

func (s *Suite) TestUserRepository_SuccessChangePassword() {
	input := &domains.ChangePassword{
		Email: "email@gmial.com",
		NewPassword: "hash_password",
		PasswordConfirm: "hash_password",
	}

	query := "UPDATE `users` SET `password` = ? WHERE (id = ?) RETURNING `users`.`id`"
	rows := sqlmock.NewRows([]string{"id"}).AddRow("uuid")

	defer s.DB.Close()
	s.mock.ExpectBegin()
	s.mock.ExpectQuery(query).WithArgs(input.NewPassword, "uuid").WillReturnRows(rows)
	s.mock.ExpectCommit()


	err := s.userRepository.UpdatePassword(input, "uuid")

	s.Nil(err)
	s.NoError(err)
}

// func (s *Suite) TestUserRepository_FailChangePassword() {
// 	input := &domains.ChangePassword{
// 		Email: "email@gmial.com",
// 		NewPassword: "hash_password",
// 		PasswordConfirm: "hash_password",
// 	}

// 	query := "UPDATE `users` SET `password` = ? WHERE (id = ?)"

// 	s.mock.ExpectBegin()
// 	s.mock.ExpectQuery(query).WithArgs("hash_password", "uuid").WillReturnError(errors.New(""))
// 	s.mock.ExpectCommit()

// 	err := s.userRepository.UpdatePassword(input, "uuid")

// 	s.NotNil(err)
// 	s.Error(err)
// }

func TestSuiteRepository(t *testing.T) {
	suite.Run(t, new(Suite))
}
