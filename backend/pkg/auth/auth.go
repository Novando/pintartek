package auth

import "strings"

func GetTokenFromBearer(bearer string) string {
	token := ""
	if strings.Contains(bearer, "Bearer") {
		token = strings.TrimPrefix(bearer, "Bearer ")
	}
	return token
}
