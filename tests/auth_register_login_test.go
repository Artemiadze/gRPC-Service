package tests

import (
	"testing"
	"time"

	"github.com/Artemiadze/gRPC-Service/tests/suite"

	ssov1 "github.com/Artemiadze/gRPC-Service/gen/go/sso"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID     = 0
	appID          = 1
	appSecret      = "test-secret"
	passDefaultLen = 10
)

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}

func TestRegisterUsers_Success(t *testing.T) {
	ctx, st := suite.New(t)

	users := []struct {
		email string
		pass  string
	}{
		{"user1@example.com", "password1"},
		{"user2@example.com", "password2"},
	}

	for _, u := range users {
		resp, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
			Email:    u.email,
			Password: u.pass,
		})
		require.NoError(t, err)
		assert.NotZero(t, resp.GetUserId())
	}
}

func TestRegister_DuplicateUser_Failure(t *testing.T) {
	ctx, st := suite.New(t)

	email := "duplicate@example.com"
	pass := "securepass"

	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)

	_, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.ErrorContains(t, err, "already exists")
}

func TestLogin_Success(t *testing.T) {
	ctx, st := suite.New(t)

	email := "login@example.com"
	pass := "loginpass"

	reg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{Email: email, Password: pass})
	require.NoError(t, err)

	loginResp, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, loginResp.GetToken())

	loginTime := time.Now()
	tokenParsed, err := jwt.Parse(loginResp.GetToken(), func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, reg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), 1)
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int64
		expectedErr string
	}{
		{"EmptyPassword", gofakeit.Email(), "", appID, "password is required"},
		{"EmptyEmail", "", randomFakePassword(), appID, "email is required"},
		{"BothEmpty", "", "", appID, "email is required"},
		{"NonMatchingPassword", gofakeit.Email(), randomFakePassword(), appID, "invalid email or password"},
		{"NoAppID", gofakeit.Email(), randomFakePassword(), emptyAppID, "app_id is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: randomFakePassword(),
			})
			_, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestIsAdmin_FalseByDefault(t *testing.T) {
	ctx, st := suite.New(t)

	reg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    "admin-check@example.com",
		Password: "adminpass",
	})
	require.NoError(t, err)

	isAdminResp, err := st.AuthClient.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: reg.GetUserId(),
	})
	require.NoError(t, err)
	assert.False(t, isAdminResp.GetIsAdmin())
}

func TestLogout_Success(t *testing.T) {
	ctx, st := suite.New(t)

	email := "logout@example.com"
	pass := "logoutpass"

	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{Email: email, Password: pass})
	require.NoError(t, err)

	loginResp, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)

	logoutResp, err := st.AuthClient.Logout(ctx, &ssov1.LogoutRequest{
		Token: loginResp.GetToken(),
	})
	require.NoError(t, err)
	assert.True(t, logoutResp.GetSuccess())
}
