package user

import (
	"encoding/base64"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/yegamble/go-tube-api/database"
	"github.com/yegamble/go-tube-api/modules/api/errorhandler"
	"github.com/yegamble/go-tube-api/modules/api/video"
	"golang.org/x/crypto/argon2"
	"gorm.io/gorm"
	"math/rand"
	"os/user"
	"time"
)

type User struct {
	gorm.Model
	ID 					uuid.UUID  `json:"id" gorm:"primary_key"`
	FirstName           string    `json:"first_name" gorm:"type:text" validate:"required,min=1,max=30"`
	LastName            string    `json:"last_name" gorm:"type:text" validate:"required,min=1,max=30"`
	Email               string    `json:"email" gorm:unique",type:text" validate:"required,email,min=6,max=32"`
	Username            string     `json:"username" gorm:"unique" validate:"required,alphanum,min=1,max=32"`
	DisplayName         string    `json:"username" gorm:"unique" validate:"max=50"`
	DateOfBirth         time.Time `json:"date_of_birth" gorm:"type:datetime" validate:"required"`
	Gender              string    `json:"gender" gorm:"type:text"`
	CurrentCity         string    `json:"current_city" gorm:"type:text"`
	HomeTown            string    `json:"hometown" gorm:"type:text"`
	Bio                 string    `json:"bio" gorm:"type:text"`
	ProfilePhoto        string     `json:"profile_photo" gorm:"type:text"`
	HeaderPhoto         string     `json:"header_photo" gorm:"type:text"`
	Password            string     `json:"password" gorm:"type:text" validate:"required,min=8,max=120"`
	Videos              []video.Video    `json:"videos" gorm:"type:text"`
	Subscribers         []*user.User  `json:"subscribers" gorm:"type:array"`
	PGPKey              string     `json:"pgp_key" gorm:"type:text"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DeletedAt           gorm.DeletedAt
}

type HashConfig struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
}

func UserFormParser(c *fiber.Ctx) ([]*errorhandler.ErrorResponse, error) {

	var body User

	formErr := ValidateStruct(&body)
	if formErr != nil {
		return formErr,nil
	}

	err := c.BodyParser(&body)
	if err != nil {
		return nil,err
	}

	return nil,nil
}

 func createUser(u *User) error {
	 db := database.DBConn

	 //var body User
	 u.ID = uuid.New()
	 db.Create(u)

	 return  nil
 }

func ValidateStruct(user *User) []*errorhandler.ErrorResponse {
	var errors []*errorhandler.ErrorResponse
	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element errorhandler.ErrorResponse
			element.FailedField = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

//encodes a string input to argon hash
func encodeToArgon(input string) string {

	c := &HashConfig{
		time:    1,
		memory:  64 * 1024,
		threads: 4,
		keyLen:  32,
	}

	// Generate a Salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return ""
	}

	hash := argon2.IDKey([]byte(input), salt, c.time, c.memory, c.threads, c.keyLen)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	format := "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
	full := fmt.Sprintf(format, argon2.Version, c.memory, c.time, c.threads, b64Salt, b64Hash)
	return full

}