package user

import (
	"context"
	"errors"
	"server/db"
	"server/util"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	secretKey = "secret"
)

type service struct {
	Repository
	timeout time.Duration
	db      db.Database
}

func NewService(repository Repository, database db.Database) Service {
	return &service{
		repository,
		time.Duration(2) * time.Second,
		database,
	}
}

func (s *service) CreateUser(c context.Context, req *CreateUserReq) (*CreateUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	u := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	r, err := s.Repository.CreateUser(ctx, u)
	if err != nil {
		return nil, err
	}

	res := &CreateUserRes{
		ID:       strconv.Itoa(int(r.Id)),
		Username: r.Username,
		Email:    r.Email,
	}

	return res, nil
}

type MyJWTClaims struct {
	Id       uint64 `json:"id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (s *service) Login(c context.Context, req *LoginUserReq) (*LoginUserRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	user, err := s.Repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &LoginUserRes{}, err
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		return &LoginUserRes{}, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, MyJWTClaims{
		Id:       user.Id,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(user.Id)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return &LoginUserRes{}, err
	}

	return &LoginUserRes{AccessToken: ss, Username: user.Username, ID: strconv.Itoa(int(user.Id))}, nil
}

func (svc *service) getUserById(id uint64) (*User, error) {
	user := &User{}
	row := svc.db.GetDB().QueryRow("SELECT id, username, email, password FROM public.users WHERE id = $1", id)
	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (svc *service) GetUserByJWT(token string) (*User, error) {
	var claims MyJWTClaims
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !jwtToken.Valid {
		return nil, errors.New("invalid token")
	}

	return svc.getUserById(claims.Id)
}
