package service

import (
	"fmt"
	http_v1 "github.com/CyganFx/table-reservation/internal/delivery/http-v1"
	"github.com/CyganFx/table-reservation/internal/domain"
	"github.com/CyganFx/table-reservation/pkg/validator/forms"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	aws_session "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"log"
	"mime/multipart"
	"net/url"
	"strings"
)

const (
	UserRoleId = 2
)

type user struct {
	repo UserRepo
}

func NewUser(repo UserRepo) *user {
	return &user{repo: repo}
}

type UserRepo interface {
	Create(name, email, mobile, hashedPassword string, roleId int) error
	GetById(id int) (*domain.User, error)
	Update(user *domain.User) error
	Authenticate(email, password string) (int, error)
	SetProfileImage(filePath string, userID int) error
	UpdateUserRoleByID(userID, roleID int) error
}

func (u *user) Save(form *forms.FormValidator) (bool, error) {
	form.Required("name", "email", "mobile", "password")
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 5)
	form.MinLength("mobile", 11)
	form.MaxLength("mobile", 12)
	form.MaxLength("name", 50)
	form.MaxLength("email", 100)

	if !form.Valid() {
		return false, nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Get("password")), 12)
	if err != nil {
		return false, errors.Wrap(err, "failed to generate hashed password")
	}

	err = u.repo.Create(form.Get("name"),
		form.Get("email"),
		form.Get("mobile"),
		string(hashedPassword),
		UserRoleId,
	)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (u *user) SignIn(email, password string) (int, error) {
	return u.repo.Authenticate(email, password)
}

func (u *user) FindById(id int) (*domain.User, error) {
	return u.repo.GetById(id)
}

func (u *user) UpdateImage(filePath string, userID int) error {
	return u.repo.SetProfileImage(filePath, userID)
}

func (u *user) Update(user *domain.User) error {
	return u.repo.Update(user)
}

func (u *user) SetConfirmData(ctx *gin.Context, reservationData *http_v1.ReservationData, tableID, eventID int, eventDescription string) (*forms.FormValidator, error) {
	session := sessions.Default(ctx)
	userChoice := session.Get("userChoice").(http_v1.UserChoice)
	userChoice.TableID = tableID
	userChoice.EventID = eventID
	userChoice.EventDescription = eventDescription
	reservationData.UserChoice = userChoice

	session.Set("userChoice", userChoice)

	form := forms.New(url.Values{})

	userID := session.Get("authenticatedUserID")
	if userID != nil {
		u, err := u.FindById(userID.(int))
		if err != nil {
			return nil, err
		}
		form.Add("name", u.Name)
		form.Add("mobile", u.Mobile)
		form.Add("email", u.Email)
	}

	session.Save()
	return form, nil
}

func (u *user) UploadImageToAWSBucket(awsSession *aws_session.Session, MyBucket, filename string, file multipart.File) error {
	uploader := s3manager.NewUploader(awsSession)
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(MyBucket),
		ACL:    aws.String("public-read"),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		return errors.Wrap(err, "failed to upload file")
	}
	return nil
}

func (u *user) DeleteImageFromAWSBucket(awsSession *aws_session.Session, imageURL, myBucket, objectsLocationURL string, infoLog *log.Logger) error {
	objectName := strings.ReplaceAll(imageURL, objectsLocationURL, "")
	svc := s3.New(awsSession)
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(myBucket),
		Key:    aws.String(objectName),
	}

	result, err := svc.DeleteObject(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return fmt.Errorf("failed to delete object from aws bucket, here is aws error code %v",
					aerr.Error())
			}
		} else {
			return fmt.Errorf("failed to delete object from aws bucket, here is aws error code %v",
				aerr.Error())
		}
	}

	infoLog.Printf("Object deleted from aws s3 bucket: %v", result)
	return nil
}

func (u *user) UpdateUserRole(userID, roleID int) error {
	return u.repo.UpdateUserRoleByID(userID, roleID)
}
