package main

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
	"github.com/peterbourgon/mergemap"
	"gopkg.in/mgo.v2/bson"
	"github.com/badoux/checkmail"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	. "GwGTeamProjectApi/config"
	. "GwGTeamProjectApi/dao"
	. "GwGTeamProjectApi/models"
)

var config = Config{}
var dao = UsersDAO{}

const (

	// For simplicity these files are in the same folder as the app binary.

	privKeyPath = "app.rsa"
	pubKeyPath  = "app.rsa.pub"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initKeys() {
	signBytes, err := ioutil.ReadFile(privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

// GET list of users

func AllUsersEndPoint(w http.ResponseWriter, r *http.Request) {
	users, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, users)
}

//GET a user by Username or Email + Password

func LoginByUsernameEmailPassword(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user User	
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	user, err:= dao.FindByUsernamePassword(user.Username, user.Password)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User Credentials")
		return
	}
	
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString(signKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}

	response := Token{tokenString}
	JsonResponse(response, w)

}


// GET a user by its ID

func FindUserEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}
	respondWithJson(w, http.StatusOK, user)
}

// POST a new user

func CreateUserEndPoint(w http.ResponseWriter, r *http.Request) {

	var user User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Error in request")
		return
	}

	if len([]rune(user.Username)) < 4 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Username must be at least 4 characters")
			return
		
	}
	if checkmail.ValidateFormat(user.Email) != err || checkmail.ValidateHost(user.Email) != err {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Invalid Email Address")
			return
		
	}

	if len([]rune(user.Password)) < 8 {
			w.WriteHeader(http.StatusForbidden)
			fmt.Println("Error logging in")
			fmt.Fprint(w, "Password must be at least 8 characters")
			return
		
	}
	token := jwt.New(jwt.SigningMethodRS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	token.Claims = claims

	tokenString, err := token.SignedString(signKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Error while signing the token")
		fatal(err)
	}

	response := Token{tokenString}

	user.ID = bson.NewObjectId()
	if err := dao.Insert(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJsonToken(w, http.StatusCreated, response, user)
}

// PUT update an existing user

func UpdateUserEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing user

func DeleteUserEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJsonToken(w http.ResponseWriter, code int, response interface{}, payload interface{}) {
	x1, _ := json.Marshal(payload)
	x2, _ := json.Marshal(response)

	var m1, m2 map[string]interface{}
	json.Unmarshal(x1, &m1)
	json.Unmarshal(x2, &m2)

	merged := mergemap.Merge(m1, m2)

	JsonResponse, _ := json.Marshal(merged)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(JsonResponse)

}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {

JsonResponse, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(JsonResponse)

}

// Parse the configuration file 'config.toml', and establish a connection to DB

func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
}

// Define HTTP request routes

func StartServer() {
	r := mux.NewRouter()
	r.HandleFunc("/signup", CreateUserEndPoint).Methods("POST")
	r.HandleFunc("/users", AllUsersEndPoint).Methods("GET")
	r.HandleFunc("/login", LoginByUsernameEmailPassword).Methods("GET")

	r.Handle("/users/{id}", negroni.New(
		negroni.HandlerFunc(ValidateTokenMiddleware),
		negroni.Wrap(http.HandlerFunc(ProtectedHandler)),
	))

	if err := http.ListenAndServe(":3100", r); err != nil {
		log.Fatal(err)
	}
}

func main() {

	initKeys()
	StartServer()
}

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {

	response := Response{"Gained access to protected resource"}
	JsonResponse(response, w)

}


func ValidateTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

	if err == nil {
		if token.Valid {
			next(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, "Token is not valid")
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "Unauthorized access to this resource")
	}

}

func JsonResponse(response interface{}, w http.ResponseWriter) {

	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
