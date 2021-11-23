package models

import (
	u "GoRestProject/utils"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

//JWT Struct

type Token struct {
	UserId   uint
	Username string
	jwt.StandardClaims
}

//Kullanıcı Tablosu Struct

type Account struct {
	gorm.Model //Migration İşlemi yaparken veritabanı üzerinde account tablosu yaratılması için belirtilir

	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

//Gelen Bilgileri Doğrulama Fonksiyonu

func (account *Account) Validate() (map[string]interface{}, bool) {

	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email Adresi Hatalıdır!"), false
	}

	if len(account.Password) < 8 {
		return u.Message(false, "Şifreniz En az 8 karakter olmalıdır!"), false
	}

	temp := &Account{}

	//Email Adresinin kayıtlı olup olmadığını kontrol ettir

	err := GetDB().Table("accounts").Where("email=?", account.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Bağlantı hatası oluştu. Lütfen tekrar deneyiniz!"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email adresi başka bir kullanıcı tarafından kullanılıyor."), false
	}
	return u.Message(false, "Her şey yolunda!"), true

}

// Kullanıcı Hesabı Yaratma Fonksiyonu

func (account *Account) Create() map[string]interface{} {
	if resp, ok := account.Validate(); !ok {
		return resp
	}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Bağlantı hatası oluştu.Kullanıcı Yaratılamadı!!")
	}

	// Yarılan Hesap için JWT oluşturulur

	tk := &Token{UserId: account.ID}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" // Yanıt İçerisinden parola silinir

	response := u.Message(true, "Hesap Başarıyla Yaratıldı")
	response["account"] = account
	return response
}

//Giriş Yapma Fonksiyonu
func Login(email, password string) map[string]interface{} {

	account := &Account{}
	err := GetDB().Table("accounts").Where("email = ?", email).First(account).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email adresi bulunamadı!")
		}
		return u.Message(false, "Bağlantı hatası oluştu. Lütfen tekrar deneyiniz!")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { // Parola eşleşmedi
		return u.Message(false, "Parola hatalı! Lütfen tekrar deneyiniz!")
	}

	// Giriş başarılı
	account.Password = ""

	// JWT yaratılır
	tk := &Token{UserId: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString // JWT yanıta eklenir

	resp := u.Message(true, "Giriş başarılı!")
	resp["account"] = account
	return resp
}

// Kullanıcı bilgilerini getirme fonksiyonu
func GetUser(u uint) *Account {
	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" { // Kullanıcı bulunamadı
		return nil
	}

	acc.Password = ""
	return acc
}
