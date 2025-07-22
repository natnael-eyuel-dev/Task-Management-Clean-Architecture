package infrastructure

// imports
import (
	"github.com/natnael-eyuel-dev/Task-Management-Clean-Architecture/Domain";
	"golang.org/x/crypto/bcrypt";
)

type passwordService struct{}

func NewPasswordService() domain.PasswordService {
	return &passwordService{}
}

// hash password
func (pswserv *passwordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// check password
func (pswserv *passwordService) CheckPassword(hashed, plain string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))
	return err == nil
}