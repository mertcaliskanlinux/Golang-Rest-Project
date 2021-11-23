//Bu Dosyada Temel Olarak
//Kimlik doğrulaması isteyen ve istemeyen endpointler belirlendi ve gelen isteğin doğrulama isteğine göre yönlendirmesi yapıldı.
//Doğrulama istenen bir yere gidilmek isteniyorsa, taleple birlikte iletilmiş olan header bilgisi içerisinden Token alındı.
//Tokenın geçerliği ve doğruluğu kontrol edildi. Eğer her şey yolundaysa talep edilen uç noktaya erişime izin verildi.

package app

import (
	"GoRestProject/models"
	u "GoRestProject/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

var JwtAuthentication = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuth := []string{"api/user/new", "/api/user/login"} // DOĞRULAMA İSTEMEYEN ENDPOİNTLER
		requestPath := r.URL.Path                              // MEVCUT İSTEK YOLU

		// GELEN İSTEĞİN DOĞRULAMA İSTEYİP İSTEMEDİĞİNİ KONTROL EDİLİR
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}
		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization") // Header'dan Token Alınır

		// TOKEN YOKSA 403 Unauthorization HATASI DÖNDÜRÜCEK
		if tokenHeader == "" {
			response = u.Message(false, "Token Gönderilemedi!")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content Type", "application/json")
			u.Respond(w, response)
		}

		splitted := strings.Split(tokenHeader, " ") // TOKEN FORMATINDA GELİP GELMEDİĞİ KONTROL EDİLİR.
		if len(splitted) != 2 {
			response = u.Message(false, "Hatalı ya da Geçersiz Token!")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
		}

		tokenPart := splitted[1] //TOKENİN DOĞRULAMA YAPMAMIZA YARAYAN KISMI ALINIR
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { // Token Hatalı ise 403 hatası dönülür
			response = u.Message(false, "Token Hatalı!")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		if !token.Valid {
			response = u.Message(false, "Token Geçersiz!")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Respond(w, response)
			return
		}

		//Doğrulama Başarılı İse İşleme Devam Edilir
		fmt.Sprintf("Kullanıcı %v", tk.Username) //Kullanıcı Adı konsola Basılır
		ctx := context.WithValue(r.Context(), "user", tk.UserId)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
