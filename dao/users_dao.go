package dao

import (
	"log"

	. "BasicAPI/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UsersDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "users"
)

func (m *UsersDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

func (m *UsersDAO) FindAll() ([]User, error) {
	var users []User
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}

func (m *UsersDAO) FindByUsernamePassword( username string, password string) (User, error) {
	var user User
	err := db.C(COLLECTION).Find(bson.M{"username" : &username, "password" : &password }).One(&user)
	return user, err	
}

func (m *UsersDAO) FindByEmailPassword( email string, password string) (User, error) {
	var user User
	err := db.C(COLLECTION).Find(bson.M{"email" : &email, "password" : &password }).One(&user)
	return user, err	
}


func (m *UsersDAO) FindByUsername( username string) (User, error) {
	var user User
	err := db.C(COLLECTION).Find(bson.M{"username" : &username }).One(&user)
	return user, err	
}

func (m *UsersDAO) FindByEmail( email string) (User, error) {
	var user User
	err := db.C(COLLECTION).Find(bson.M{"email" : &email }).One(&user)
	return user, err	
}


func (m *UsersDAO) FindById(id string) (User, error) {
	var user User
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

func (m *UsersDAO) Insert(user User) error {
	err := db.C(COLLECTION).Insert(&user)
	return err
}

func (m *UsersDAO) Delete(user User) error {
	err := db.C(COLLECTION).Remove(&user)
	return err
}

func (m *UsersDAO) Update(user User) error {
	err := db.C(COLLECTION).UpdateId(user.ID, &user)
	return err
}
