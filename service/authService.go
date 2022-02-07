package service

import (
	"errors"
	"log"

	"github.com/dgrijalva/jwt-go"
	"github.com/manishdangi98/banking-auth/domain"
	"github.com/manishdangi98/banking-auth/dto"
	"github.com/manishdangi98/banking-lib/errs"
)

type AuthService interface {
	Login(dto.LoginRequest) (*string, *errs.AppError)
	Verify(urlParams map[string]string) (bool, error)
}
type DefaultAuthService struct {
	repo           domain.AuthRepository
	rolePermission domain.RolePermission
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*string, *errs.AppError) {
	login, err := s.repo.FindBy(req.Username, req.Password)
	if err != nil {
		return nil, err
	}
	token, err := login.GenrateToken()
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s DefaultAuthService) Verify(urlParams map[string]string) (bool, error) {
	//convert the string token to JWT struct
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return false, err
	} else {
		/*
		   Checking the validity of the token, this verifies the expiry
		   time and the signature of the token
		*/
		if jwtToken.Valid {
			//typecast the token to jwt.MapClaims
			mapClaims := jwtToken.Claims.(jwt.MapClaims)
			//converting the token claims to Claims struct
			if claims, err := domain.BuildClaimsFromJwtMapClaims(mapClaims); err != nil {
				return false, err
			} else {
				if claims.IsUserRole() {
					if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
						return false, nil
					}
				}
				//verify of the role is authorized to use the route
				isAuthorized := s.rolePermission.IsAuthorizedFor(claims.Role, urlParams["routeName"])
				return isAuthorized, nil
			}
		} else {
			return false, errors.New("invalid token")
		}
	}
}
func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		log.Println("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}

func NewLoginService(repo domain.AuthRepository, permission domain.RolePermission) DefaultAuthService {
	return DefaultAuthService{repo, permission}
}
