package v1

import (
  "context"
  "time"

  v1 "github.com/ckbball/os-company/pkg/api/v1"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
)

type repository interface {
  Create(*v1.Company) (string, error)
  Update(*v1.Company, string) (int64, int64, error)
  Delete(string) (int64, error)
  GetById(string) (*Company, error)
  GetByEmail(string) (*Company, error)
  GetByName(string) (*Company, error)
  FilterCompanys(*v1.FindRequest) ([]*Company, error)
  UpdateActive(string) (int64, error)
}

type CompanyRepository struct {
  cs *mongo.Collection
}

func NewCompanyRepository(client *mongo.Collection) *CompanyRepository {
  return &CompanyRepository{
    cs: client,
  }
}

func (repository *CompanyRepository) Create(company *v1.Company) (string, error) {
  // add a duplicate email and a duplicate company name check

  insertCompany := bson.D{
    {"email", company.Email},
    {"password", company.Password},
    {"name", company.Name},
    {"mission", company.Mission},
    {"last_active", company.LastActive},
    {"location", company.Location},
  }

  result, err := repository.cs.InsertOne(context.TODO(), insertCompany)

  if err != nil {
    return "", err
  }

  id := result.InsertedID
  w, _ := id.(primitive.ObjectID)

  out := w.Hex()

  return out, err
}

func (repository *CompanyRepository) Update(company *v1.Company, id string) (int64, int64, error) {
  // add a duplicate email and a duplicate companyname check

  primitiveId, _ := primitive.ObjectIDFromHex(id)

  insertCompany := bson.D{
    {"email", company.Email},
    {"password", company.Password},
    {"name", company.Name},
    {"mission", company.Mission},
    {"last_active", company.LastActive},
    {"location", company.Location},
  }

  result, err := repository.cs.UpdateOne(context.TODO(),
    bson.D{
      {"_id", primitiveId},
    },
    bson.D{
      {"$set", insertCompany},
    },
  )

  if err != nil {
    return -1, -1, err
  }

  return result.MatchedCount, result.ModifiedCount, nil
}

func (repository *CompanyRepository) Delete(id string) (int64, error) {
  primitiveId, _ := primitive.ObjectIDFromHex(id)
  filter := bson.D{{"_id", primitiveId}}

  result, err := repository.cs.DeleteOne(context.TODO(), filter)
  if err != nil {
    return -1, err
  }
  return result.DeletedCount, nil
}

func (s *CompanyRepository) GetById(id string) (*Company, error) {
  primitiveId, _ := primitive.ObjectIDFromHex(id)

  var company Company
  err := s.cs.FindOne(context.TODO(), Company{Id: primitiveId}).Decode(&company)
  if err != nil {
    return nil, err
  }

  return &company, nil
}

func (s *CompanyRepository) GetByEmail(email string) (*Company, error) {

  var company Company
  err := s.cs.FindOne(context.TODO(), Company{Email: email}).Decode(&company)
  if err != nil {
    return nil, err
  }

  return &company, nil
}

func (s *CompanyRepository) GetByName(name string) (*Company, error) {

  var company Company
  err := s.cs.FindOne(context.TODO(), Company{Name: name}).Decode(&company)
  if err != nil {
    return nil, err
  }

  return &company, nil
}

func (s *CompanyRepository) UpdateActive(id string) (int64, error) {

  now := time.Now()
  secs := now.Unix()

  var company Company
  err := s.cs.FindOne(context.TODO(), Company{id: id}).Decode(&company)
  if err != nil {
    return nil, err
  }

  primitiveId, _ := primitive.ObjectIDFromHex(id)

  insertCompany := bson.D{
    {"email", company.Email},
    {"password", company.Password},
    {"name", company.Name},
    {"mission", company.Mission},
    {"last_active", secs},
    {"location", company.Location},
  }

  result, err := repository.cs.UpdateOne(context.TODO(),
    bson.D{
      {"_id", primitiveId},
    },
    bson.D{
      {"$set", insertCompany},
    },
  )

  if err != nil {
    return nil, err
  }

  return secs, nil
}
