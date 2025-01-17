package utility

import (
	"github.com/golang-jwt/jwt"
	"time"
	"walk-server/utility/initial"
)

// JwtData 一些结构体的定义
type JwtData struct {
	OpenID string `json:"open_id"`
	jwt.StandardClaims
}

// GenerateStandardJwt 根据数据生成带有 standard claims 的 jwt token
func GenerateStandardJwt(jwtData *JwtData) (string, error) {
	claims := jwtData
	claims.StandardClaims = jwt.StandardClaims{
		// 过期时间
		ExpiresAt: time.Now().Add(48 * time.Hour).Unix(), // 设置 2 天后过期
		// 指定token发行人
		Issuer: "JHWL",
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString([]byte(initial.Config.GetString("server.JWTSecret")))
	return token, err
}

// ParseToken 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func ParseToken(token string) (*JwtData, error) {
	jwtSecret := []byte(initial.Config.GetString("server.JWTSecret"))
	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &JwtData{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*JwtData); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}
