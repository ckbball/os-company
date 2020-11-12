package v1

import (
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
  Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  Email      string             `json:"email,omitempty" bson:"email,omitempty"`
  Password   string             `json:"password,omitempty" bson:"password,omitempty"`
  Firstname  string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
  Lastname   string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
  Phone      string             `json:"phone,omitempty" bson:"phone,omitempty"`
  LastActive int                `json:"lastActive,omitempty" bson:"last_active,omitempty"`
}

type Profile struct {
  Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  UserId       primitive.ObjectID `json:"userid,omitempty" bson:"userid,omitempty"`
  Objective    string             `json:"objective,omitempty" bson:"objective,omitempty"`
  Technologies []string           `json:"technologies,omitempty" bson:"technologies,omitempty"`
  Statements   []string           `json:"statements,omitempty" bson:"statements,omitempty"`
}
