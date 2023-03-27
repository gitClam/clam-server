package jwts

import (
	"clam-server/config"
	"clam-server/serverlogger"
	"clam-server/utils/uuid"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"time"

	"sync"
)

type (
	Jwt struct {
		Config Config
	}
)

var (
	jwts *Jwt
	lock sync.Mutex
)

func Serve(ctx *gin.Context) bool {
	ConfigJWT()
	pass, err := jwts.CheckJWT(ctx)
	if pass && err != nil {
		_, err = GenerateToken()
	}
	if err != nil {
		jwts.Config.ErrorHandler(ctx, err.Error())
		return false
	}
	return true
}

// ParseToken 解析token的信息为当前用户
func ParseToken(ctx *gin.Context) (string, bool) {
	mapClaims := (jwts.Get(ctx).Claims).(jwt.MapClaims)

	uid, ok := mapClaims["uuid"].(string)
	if !ok {
		ctx.JSON(http.StatusMethodNotAllowed, gin.H{
			"message": "ParseToken err",
		})
		return "", false
	}
	return uid, true
}

func FromAuthHeader(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		return "", nil // No error, just no token
	}
	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "jwt" {
		return "", fmt.Errorf("Authorization header format must be JWT {token}")
	}

	return authHeaderParts[1], nil
}

func (m *Jwt) logf(format string, args ...interface{}) {
	if m.Config.Debug {
		serverlogger.Warn(fmt.Sprintf(format, args))
	}
}

func (m *Jwt) Get(ctx *gin.Context) *jwt.Token {
	if token, exist := ctx.Get(m.Config.ContextKey); exist {
		return token.(*jwt.Token)
	}
	return nil
}

func (m *Jwt) CheckJWT(ctx *gin.Context) (bool, error) {
	if !m.Config.EnableAuthOnOptions && ctx.Request.Method == http.MethodOptions {
		return true, nil
	}

	token, err := m.Config.Extractor(ctx)
	if err != nil {
		return false, fmt.Errorf("Error extracting token: %v", err)
	}

	if token == "" {
		if m.Config.CredentialsOptional {
			return true, fmt.Errorf("Error: No credentials found (CredentialsOptional=true)")
		}
		return false, fmt.Errorf("Error: No credentials found (CredentialsOptional=false)")
	}

	parsedToken, err := jwt.Parse(token, m.Config.ValidationKeyGetter)
	if err != nil {
		return false, fmt.Errorf("Error parsing token2: %v", err)
	}

	if m.Config.SigningMethod != nil && m.Config.SigningMethod.Alg() != parsedToken.Header["alg"] {
		message := fmt.Sprintf("Expected %s signing method but token specified %s",
			m.Config.SigningMethod.Alg(),
			parsedToken.Header["alg"])
		return false, fmt.Errorf("Error validating token algorithm: %s", message)
	}

	if !parsedToken.Valid {
		return false, fmt.Errorf("parsedToken.Valid = false")
	}

	if m.Config.Expiration {
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			if expired := claims.VerifyExpiresAt(time.Now().Unix(), true); !expired {
				return false, fmt.Errorf("token expired")
			}
		}
	}

	ctx.Set(m.Config.ContextKey, parsedToken)
	return true, nil
}

// ------------------------------------------------------------------------
// ------------------------------------------------------------------------

// ConfigJWT jwt中间件配置
func ConfigJWT() {
	if jwts != nil {
		return
	}

	lock.Lock()
	defer lock.Unlock()

	if jwts != nil {
		return
	}

	c := Config{
		ContextKey: config.GetConfig().Jwt.DefaultContextKey,
		// 这个方法将验证jwt的token
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			// 自己加密的秘钥或者说盐值
			return []byte(config.GetConfig().Jwt.Secret), nil
		},
		// 设置后，中间件会验证令牌是否使用特定的签名算法进行签名
		// 如果签名方法不是常量，则可以使用ValidationKeyGetter回调来实现其他检查
		// 重要的是要避免此处的安全问题：https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
		// 加密的方式
		SigningMethod: jwt.SigningMethodHS256,
		// 验证未通过错误处理方式
		ErrorHandler: func(ctx *gin.Context, errMsg string) {
			ctx.JSON(http.StatusMethodNotAllowed, gin.H{
				"message": "errMsg",
			})
		},
		// 指定func用于提取请求中的token
		Extractor:           FromAuthHeader,
		Expiration:          true,
		Debug:               true,
		EnableAuthOnOptions: false,
	}
	jwts = &Jwt{Config: c}
}

type Claims struct {
	uuid string `json:"name"`
	jwt.StandardClaims
}

// GenerateToken 在登录成功的时候生成token
func GenerateToken() (string, error) {

	expireTime := time.Now().Add(time.Duration(config.GetConfig().Jwt.JwtTimeout) * time.Second)

	claims := Claims{
		uuid.GenUUID(),
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "gin-jwt",
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := tokenClaims.SignedString([]byte(config.GetConfig().Jwt.Secret))
	return token, err
}

type Config struct {
	// The function that will return the Key to validate the JWT.
	// It can be either a shared secret or a public key.
	// Default value: nil
	ValidationKeyGetter jwt.Keyfunc
	// The name of the property in the request where the user (&token) information
	// from the JWT will be stored.
	// Default value: "jwts"
	ContextKey string
	// The function that will be called when there's an error validating the token
	// Default value:
	ErrorHandler func(*gin.Context, string)
	// A boolean indicating if the credentials are required or not
	// Default value: false
	CredentialsOptional bool
	// A function that extracts the token from the request
	// Default: FromAuthHeader (i.e., from Authorization header as bearer token)
	Extractor func(*gin.Context) (string, error)
	// Debug flag turns on debugging output
	// Default: false
	Debug bool
	// When set, all requests with the OPTIONS method will use authentication
	// if you enable this option you should register your route with iris.Options(...) also
	// Default: false
	EnableAuthOnOptions bool
	// When set, the middelware verifies that tokens are signed with the specific signing algorithm
	// If the signing method is not constant the ValidationKeyGetter callback can be used to implement additional checks
	// Important to avoid security issues described here: https://auth0.com/blog/2015/03/31/critical-vulnerabilities-in-json-web-token-libraries/
	// Default: nil
	SigningMethod jwt.SigningMethod
	// When set, the expiration time of token will be check every time
	// if the token was expired, expiration error will be returned
	// Default: false
	Expiration bool
}
