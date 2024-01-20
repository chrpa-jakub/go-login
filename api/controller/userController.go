package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/chrpa-jakub/register-api/database"
	"github.com/chrpa-jakub/register-api/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)
var body struct {
  Login string `json:"login"`
  Password string `json:"password"`
}

var ctx = context.Background()

func Register(c *gin.Context){

  if c.BindJSON(&body) != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve body."})
    return
  }

  if database.DB.Exists(ctx, body.Login).Val() != 0 {

    c.JSON(http.StatusBadRequest, gin.H{"error": "User with this login already exists!"})
    return
  }

  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)

  if err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error":"Failed to hash password."})
    return
  }

  user := model.User{Login: body.Login, PasswordHash: string(hashedPassword)}
  userJson, _ := json.Marshal(user)

  database.DB.Set(ctx, user.Login, userJson, 0)

  jwtToken := createJwt(user)
  c.SetCookie("jwt", jwtToken, 60*60*24*30, "", "", true, true)
  c.JSON(http.StatusOK, jwtToken)
}

func Login(c *gin.Context){
  if c.BindJSON(&body) != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error":"Failed to retreive body."})
    return
  }

  userFromDb := database.DB.Get(ctx, body.Login).Val()
  fmt.Println(userFromDb)

  if userFromDb == "" {
    c.JSON(http.StatusBadRequest, gin.H{"error":"User with this login does not exist!"})
    return
  }

  var user model.User
  json.Unmarshal([]byte(userFromDb), &user)

  cmp := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password))  
  if cmp != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error":"Passwords do not match!"})
    return
  }

  jwtToken := createJwt(user)
  c.SetCookie("jwt", jwtToken, 60*60*24*30, "", "", true, true)
  c.JSON(http.StatusOK, jwtToken)
}


func createJwt(user model.User) (tokenString string) {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256,
  jwt.MapClaims{
    "usr" : user.Login,
    "exp" : time.Now().Add(time.Hour*24*30).Unix(),
  })

  tokenString, _ = token.SignedString([]byte(os.Getenv("SECRET")))
  return tokenString
}
