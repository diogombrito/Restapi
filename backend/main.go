package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func main() {
	testdb()
	router.Use(CORS)
	router.POST("/login", login)
	router.PUT("/create", create)
	router.GET("/read", read)
	router.POST("/update", update)
	router.DELETE("/delete", delete)
	log.Fatal(router.Run(":8080"))
}

type loginBan struct {
	uID      uint64
	try      uint64
	time     time.Time
	isBanned bool
}

var listLogin = []loginBan{}

type userType struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

type insertUser struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Age      string `json:"birth"`
	Name     string `json:"name"`
	Family   string `json:"family"`
	Role     string `json:"role"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	Username  string `json:"username"`
	ExpiresIn int64  `json:"expiresIn"`
	Admin     bool   `json:"admin"`
}

type ReadResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Birth    string `json:"birth"`
	Name     string `json:"name"`
	Family   string `json:"family"`
	Admin    string `json:"admin"`
}

func login(c *gin.Context) {
	var u userType
	var resp AuthResponse
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "INVALID_JSON")
		return
	}
	log.Println(u.ID)
	//compare the user from the request, with the one we defined:
	user := selectUsername(u.Username)
	fmt.Println(user)
	fmt.Println("username: ", user.user)
	fmt.Println("password: ", user.pass)
	passBool := CompareHash(u.Password, user.pass)

	u.ID = uint64(selectUsername(u.Username).idP)
	if user.user != u.Username || !passBool {
		logB, indx, isinL := getIDStruct(u.ID, listLogin)
		//fmt.Println("logb: ", logB)
		fmt.Println("listLogin: ", listLogin)
		if isinL {
			fmt.Println("logB.try: ", logB.try)
			if logB.try < 3 {
				logB.try++
				listLogin[indx].try = logB.try
				listLogin[indx].time = time.Now()
				fmt.Println("new timer: ", listLogin[indx].time)
			} else {
				//banning user and show time
				listLogin[indx].isBanned = true
				newtime := time.Now()
				diff := newtime.Sub(listLogin[indx].time)
				fmt.Println("esta banido ", diff)
				//banned expired
				if diff > (30 * time.Second) {
					listLogin[indx].try = 0
					listLogin[indx].time = time.Time{}
					fmt.Println("tirado o ban")

				} else {
					//time left
					timerf := (30 * time.Second) - diff
					fmt.Println("falta: ", timerf)
					fmt.Println(timerf)
					c.JSON(http.StatusUnauthorized, gin.H{
						"error": timerf / 1000000000,
					})
					return
				}
			}

		} else {
			listLogin = append(listLogin, loginBan{
				uID:  u.ID,
				try:  1,
				time: time.Time{},
			})
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "INVALID_PASSWORD",
		})
		return
	}
	token, err := createToken(u.ID)
	if err != nil {
		log.Println("token error")
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	resp.Token = token
	resp.Username = u.Username
	resp.ExpiresIn = time.Now().Add(time.Minute * 30).Unix()
	resp.Admin = isAdminUser(u.Username)
	fmt.Println(resp)
	c.JSON(http.StatusOK, resp)
}

func create(c *gin.Context) {
	var ui insertUser
	///verificacao do json
	if err := c.ShouldBindJSON(&ui); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "INVALID_JSON",
		})
		return
	}
	log.Println("user a inserir: ", ui)
	Token := ""
	//token from header
	if values, _ := c.Request.Header["Authorization"]; len(values) > 0 {
		Token = values[0]
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "INVALID_TOKEN",
		})
		return
	}
	//token verification
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("jdnfksdmfksd"), nil
	})
	log.Println(token)
	if err != nil {
		fmt.Printf("%!", err.Error())
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	//
	userid := claims["user_id"].(float64)
	log.Println("user id: ", userid)
	// admin verification
	if !isAdmin(uint64(userid)) {
		fmt.Println("nao tem permissao")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "NO_PERMISSION",
		})
		return
	}
	// required camps
	if len(ui.Username) == 0 || len(ui.Password) == 0 {
		c.JSON(http.StatusBadRequest, "username e password tem se der preenchidos")
		return
	}
	// encrypt password
	password := GeneratehashAndSalt(ui.Password)
	lastid, inscheck := insertID(ui.Username, ui.Age, ui.Name, ui.Family, password, ui.Role)

	if !inscheck {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "INSERT_ERROR",
		})
		return
	}
	msg := "username inserido " + string(lastid)

	c.JSON(http.StatusOK, msg)
	return
}

func update(c *gin.Context) {
	var ui insertUser
	///json verification
	if err := c.ShouldBindJSON(&ui); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "INVALID_JSON",
		})
		return
	}
	log.Println("user a inserir: ", ui)
	Token := ""
	//take token from header
	if values, _ := c.Request.Header["Authorization"]; len(values) > 0 {
		Token = values[0]
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "INVALID_TOKEN",
		})
		return
	}
	//token verification
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("jdnfksdmfksd"), nil
	})
	log.Println(token)
	if err != nil {
		fmt.Printf("%!", err.Error())
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	// get user id from token
	userid := claims["user_id"].(float64)
	log.Println("user id: ", userid)

	// admin verification
	if !isAdmin(uint64(userid)) {
		fmt.Println("nao tem permissao")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "NO_PERMISSION",
		})
		return
	}
	//password := GeneratehashAndSalt(ui.Password)
	log.Println("user ", ui.ID)
	updated := updateID(int(ui.ID), ui.Username, ui.Name, ui.Family, ui.Role)
	if !updated {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "UPDATE_ERROR",
		})
	}
	log.Println("updated: ", updated)
	c.JSON(http.StatusOK, "user updated")
}

func delete(c *gin.Context) {
	///verificacao do token
	Token := ""
	user := ""
	//Obter o token do header
	if values, _ := c.Request.Header["Authorization"]; len(values) > 0 {
		Token = values[0]
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "INVALID_TOKEN",
		})
		return
	}
	fmt.Println("params: ", c.Request.Header["params"])

	user = c.Request.URL.Query().Get("user")
	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "INVALID_USER",
		})
	}
	fmt.Println("USER A APAGAR : ", user)

	//verificar o token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("jdnfksdmfksd"), nil
	})
	log.Println(token)
	if err != nil {
		fmt.Printf("%!", err.Error())
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	// expiration?
	userid := claims["user_id"].(float64)
	log.Println("user id: ", userid)

	// verificar se e admin
	if !isAdmin(uint64(userid)) {
		fmt.Println("nao tem permissao")
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "NO_PERMISSION",
		})
		return
	}
	if !deleteID(user) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "DELETE_ERROR",
		})

	}
	c.JSON(http.StatusOK, "apagado ")
}

func read(c *gin.Context) {
	//var resp ReadResponse
	Token := ""
	userid := 0.0
	//Obter o token do header
	if values, _ := c.Request.Header["Authorization"]; len(values) > 0 {
		Token = values[0]
		fmt.Println("token ", Token)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "INVALID_TOKEN",
		})
		return
	}
	//verificar o token
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(Token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("jdnfksdmfksd"), nil
	})
	log.Println(token)
	if err != nil {
		fmt.Printf("%!", err.Error())
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	// expiration?

	//obter keys
	userid = claims["user_id"].(float64)
	fmt.Println(userid)
	fmt.Println(int(userid))
	//implementacao do read
	alluseres := selectAll()

	fmt.Println("ALL USERS: ", alluseres)

	resp := []ReadResponse{}

	for i := 0; i < len(alluseres); i++ {
		resp = append(resp, ReadResponse{
			ID:       alluseres[i].idP,
			Username: alluseres[i].user,
			Birth:    alluseres[i].ageP,
			Name:     alluseres[i].nameP,
			Family:   alluseres[i].familyP,
			Admin:    alluseres[i].roleP})
	}

	fmt.Println("resposta: ", resp)

	c.JSON(http.StatusOK, resp)
}

func CORS(c *gin.Context) {

	// First, we add the headers with need to enable CORS
	// Make sure to adjust these headers to your needs
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-Type", "application/json")

	fmt.Println(c)
	// Second, we handle the OPTIONS problem
	if c.Request.Method != "OPTIONS" {

		c.Next()
	} else {

		// Everytime we receive an OPTIONS request,
		// we just return an HTTP 200 Status Code
		// Like this, Angular can now do the real
		// request using any other method than OPTIONS

		c.AbortWithStatus(http.StatusOK)
	}

}
