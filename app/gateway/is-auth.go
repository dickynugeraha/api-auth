package gateway

// import (
// 	"fmt"
// 	"go/token"
// 	"net/http"
// 	"strings"

// 	"github.com/golang-jwt/jwt/v4"
// )

// func IsAuthorization(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.URL.Path == "/login" {
// 			next.ServeHTTP(w, r)
// 		}

// 		authHeader := r.Header.Get("Authorization")

// 		if !strings.Contains(authHeader, "Bearer") {
// 			http.Error(w, "invalid token not bearer", http.StatusBadRequest)
// 			return
// 		}

// 		justToken := strings.Replace(authHeader, "Bearer ", "", -1)

// 		token, err := jwt.Parse(justToken, func(token *jwt.Token) (interface{}, error){
// 			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fmt.Errorf("Signin method invalid")
// 			} else if method != jwt.SigningMethodHS256 {
// 				return nil, fmt.Errorf("Signin method invalid")
// 			}

// 			return jwt.SigningMethodHS256, nil
// 		})
		
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return 
// 		}

// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if !ok || !token.Valid {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return 
// 		}


// 	})
// }