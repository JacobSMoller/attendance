
package api
import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/http"
	"github.com/dgrijalva/jwt-go"
)

// Service to define route handlers as methods.
type Service struct {
	DB        *gorm.DB
	GwKey     string
	AuthToken string
	DbHost    string
	DbName    string
	DbUser    string
	DbPw      string
}

// Config to pass env variables to service.
type Config struct {
	DB        *gorm.DB
	GwKey     string `required:"true" split_words:"true"`
	AuthToken string `required:"true" split_words:"true"`
	DbHost    string `required:"true" split_words:"true"`
	DbName    string `required:"true" split_words:"true"`
	DbUser    string `required:"true" split_words:"true"`
	DbPw      string `required:"true" split_words:"true"`
}

// VerifyGWRequest verifies incoming requests to the API by checking the X-Gwapi-Signature header.
func (s *Service) VerifyGWRequest(r *http.Request) error {
	gwJwt := r.Header.Get("X-Gwapi-Signature")
	token, err := jwt.Parse(gwJwt, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(s.AuthToken), nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("Could not verify gwapi signature for request")
	}
	return nil
}


// NewAPIService sets up a new service from env vars.
func NewAPIService(conf Config) *Service {
	return &Service{
		DB: conf.DB,
		GwKey: conf.GwKey,
		AuthToken: conf.AuthToken,
		DbHost: conf.DbHost,
		DbName: conf.DbName,
		DbUser: conf.DbUser,
		DbPw: conf.DbPw,

	}
}
