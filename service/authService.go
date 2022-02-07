package service

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
	"github.com/manishdangi98/banking-auth/domain"
	"github.com/manishdangi98/banking-auth/dto"
	"github.com/manishdangi98/banking-lib/errs"
	"github.com/manishdangi98/banking-lib/logger"
)

type AuthService interface {
	Login(dto.LoginRequest) (*dto.LoginResponse, *errs.AppError)
	Verify(urlParams map[string]string) *errs.AppError
}
type DefaultAuthService struct {
	repo           domain.AuthRepository
	rolePermission domain.RolePermission
}

func (s DefaultAuthService) Login(req dto.LoginRequest) (*dto.LoginRequest, *errs.AppError) {
	var appErr *errs.AppError
	var login *domain.Login
	if login, appErr = s.repo.FindBy(req.Username, req.Password); appErr != nil {
		return nil, appErr
	}
	claims := login.ClaimsForAccessToken()
	authToken := domain.NewAuthToken(claims)
	var accessToken string
	if accessToken, appErr = authToken.NewAccessToken(); appErr != nil {
		return nil, appErr
	}
	return &dto.LoginResponse{AccessToken: accessToken}, nil
}

func (s DefaultAuthService) Verify(urlParams map[string]string) *errs.AppError {
	//convert the string token to JWT struct
	if jwtToken, err := jwtTokenFromString(urlParams["token"]); err != nil {
		return errs.NewAuthenticationError(err.Error())
	} else {
		/*
		   Checking the validity of the token, this verifies the expiry
		   time and the signature of the token
		*/
		if jwtToken.Valid {
			//typecast the token to jwt.MapClaims
			claims := jwtToken.Claims.(*domain.Claims)

			if claims.IsUserRole() {
				if !claims.IsRequestVerifiedWithTokenClaims(urlParams) {
					return errs.NewAuthorizationError("request not verified with the token claims")
				}
			}
			isAuthorized := s.rolePermission.IsAuthorizedFor(claims.Role, urlParams["routeName"])
			if !isAuthorized {
				return errs.NewAuthorizationError(fmt.Sprintf("%s role is not authorized", claims.Role))

			}
			return nil
		} else {
			return errs.NewAuthorizationError("Invalid token")
		}
	}
}

func jwtTokenFromString(tokenString string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		logger.Error("Error while parsing token: " + err.Error())
		return nil, err
	}
	return token, nil
}

func NewLoginService(repo domain.AuthRepository, permission domain.RolePermission) DefaultAuthService {
	return DefaultAuthService{repo, permission}
}
