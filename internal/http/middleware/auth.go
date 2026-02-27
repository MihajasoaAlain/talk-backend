package middleware

import (
	"log"
	"net/http"
	"regexp"
	"strings"
	"talk-backend/internal/http/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const CtxUserIDKey = "userID"

var uuidV4LikeRe = regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)

func RequireAuth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			log.Println("[AUTH] No Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  response.CodeUnauthorized,
				"error": response.MsgUnauthorized,
			})
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" || parts[1] == "" {
			log.Println("[AUTH] Invalid Authorization format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  response.CodeUnauthorized,
				"error": response.MsgUnauthorized,
			})
			return
		}

		tokenStr := parts[1]
		log.Printf("[AUTH] Token received: %s...", tokenStr[:min(20, len(tokenStr))])

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			log.Printf("[AUTH] Token parse error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  response.CodeUnauthorized,
				"error": response.MsgUnauthorized,
			})
			return
		}

		if token == nil || !token.Valid {
			log.Println("[AUTH] Token is nil or invalid")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  response.CodeUnauthorized,
				"error": response.MsgUnauthorized,
			})
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("[AUTH] Cannot get claims")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  response.CodeUnauthorized,
				"error": response.MsgUnauthorized,
			})
			return
		}

		sub, ok := claims["sub"]
		if !ok {
			log.Println("[AUTH] No sub claim")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  response.CodeUnauthorized,
				"error": response.MsgUnauthorized,
			})
			return
		}

		log.Printf("[AUTH] sub claim: %v (type: %T)", sub, sub)

		var userID string
		switch v := sub.(type) {
		case string:
			if !isUUID(v) {
				log.Printf("[AUTH] Invalid UUID sub: %s", v)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":  response.CodeUnauthorized,
					"error": response.MsgUnauthorized,
				})
				return
			}
			userID = v
		default:
			log.Printf("[AUTH] Unknown sub type: %T", sub)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":  response.CodeUnauthorized,
				"error": response.MsgUnauthorized,
			})
			return
		}

		log.Printf("[AUTH] Success! userID=%s", userID)
		c.Set(CtxUserIDKey, userID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (string, bool) {
	v, ok := c.Get(CtxUserIDKey)
	if !ok {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}

func isUUID(v string) bool {
	return uuidV4LikeRe.MatchString(v)
}
