package middleware

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"loante/global"
	"loante/service/model"
	"loante/service/resp"
	"strings"
)


type TokenPayload struct {
	Id int
	Username string
	UserType string
	RoleName string
	RoleId string
}


func parse(token string) (*jwt.Token, error) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(global.C.Safety.Secret), nil
	})
}

func verify(token string) (*TokenPayload, error) {
	parsed, err := parse(token)
	if err != nil {
		return nil, err
	}
	// Parsing token claims
	claims, ok := parsed.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}
	// Getting ID, it's an interface{} so I need to cast it to uint
	id, ok := claims["Id"].(float64)
	if !ok {
		return nil, errors.New("Something went wrong")
	}
	username, ok := claims["Username"].(string)
	if !ok {
		return nil, errors.New("Something went wrong")
	}
	roleId := fmt.Sprintf("%v", claims["RoleId"] )
	//roleId, ok := claims["RoleId"].(string)
	//if !ok {
	//	return nil, errors.New("Something went wrong")
	//}
	//rId,_ := strconv.Atoi(roleId)
	return &TokenPayload{
		Id: int(id),
		Username: username,
		RoleId: roleId,
		UserType:  claims["UserType"].(string),
	}, nil
}

func Auth(c *fiber.Ctx) error {
	h := c.Get("token")
	if h == "" {
		return resp.Err(c, 1001, "token error!")
	}
	// Spliting the header
	chunks := strings.Split(h, " ")

	// If header signature is not like `Bearer <token>`, then throw
	// This is also required, otherwise chunks[1] will throw out of bound error
	if len(chunks) < 2 {
		return resp.Err(c, 1002, "token error!")
	}

	// Verify the token which is in the chunks
	user, err := verify(chunks[1])

	if err != nil {
		return resp.Err(c, 1003, "token error!")
	}

	//是否有权限
	rights := ""
	ok := true

	if rights, ok = global.AdminRights[user.RoleId]; !ok {
		lists,_ := new(model.AdminRight).Gets(fmt.Sprintf("id > 0"))
		for _, item := range lists{
			global.AdminRights[fmt.Sprintf("%d", item.Id)] = item.Rights
			if user.RoleId == fmt.Sprintf("%d", item.Id){
				rights = item.Rights
			}
		}
	}
	if rights != "*" && strings.Index(rights, c.Path()) == -1{
		return resp.Err(c, 1, "insufficient permissions!")
	}
	c.Locals("userType", user.UserType)
	c.Locals("userId", user.Id)
	c.Locals("userName", user.Username)
	c.Locals("roleName", user.RoleName)
	c.Locals("roleId", user.RoleId)
	c.Next()
	return nil
}
