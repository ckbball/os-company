package v1

import (
  "go.mongodb.org/mongo-driver/bson/primitive"
)

type Company struct {
  Id         primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
  Email      string             `json:"email,omitempty" bson:"email,omitempty"`
  Password   string             `json:"password,omitempty" bson:"password,omitempty"`
  Name       string             `json:"name,omitempty" bson:"name,omitempty"`
  Mission    string             `json:"mission,omitempty" bson:"mission,omitempty"`
  Location   string             `json:"location,omitempty" bson:"location,omitempty"`
  LastActive int                `json:"lastActive,omitempty" bson:"last_active,omitempty"`
}
