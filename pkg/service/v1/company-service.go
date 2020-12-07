package v1

import (
  "context"
  "errors"
  "fmt"

  "golang.org/x/crypto/bcrypt"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/metadata"
  "google.golang.org/grpc/status"

  v1 "github.com/ckbball/os-company/pkg/api/v1"
)

type handler struct {
  repo         repository
  tokenService Authable
}

func NewCompanyServiceServer(repo repository, tokenService Authable) *handler {
  return &handler{
    repo:         repo,
    tokenService: tokenService,
  }
}

func (s *handler) checkAPI(api string) error {
  if len(api) > 0 {
    if apiVersion != api {
      return status.Errorf(codes.Unimplemented,
        "unsupported API version: service implements API version '%s', but asked for '%s'", apiVersion, api)
    }
  }
  return nil
}

func (s *handler) CreateCompany(ctx context.Context, req *v1.UpsertRequest) (*v1.UpsertResponse, error) {
  // check api version
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  // generate hash of password
  hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Company.Password), bcrypt.DefaultCost)
  if err != nil {
    return nil, errors.New(fmt.Sprintf("error hashing password: %v", err))
  }
  req.Company.Password = string(hashedPass)

  id, err := s.repo.Create(req.Company)
  if err != nil {
    return nil, err
  }

  // return
  return &v1.UpsertResponse{
    Api:    apiVersion,
    Status: "Created",
    Id:     id,
    // maybe in future add more data to response about the added user.
  }, nil
}

func (s *handler) Login(ctx context.Context, req *v1.UpsertRequest) (*v1.UpsertResponse, error) {
  // check api version
  if err := s.checkAPI(req.Api); err != nil {
    return nil, err
  }

  // get company from email
  company, err := s.repo.GetByEmail(req.Email)
  if err != nil {
    return nil, err
  }

  // Compare given password to stored hash
  if err = bcrypt.CompareHashAndPassword([]byte(company.Password), []byte(req.Password)); err != nil {
    return nil, err
  }

  intId := company.Id.Hex()

  companyModel := &v1.Company{
    Id:       intId, //
    Email:    company.Email,
    Password: company.Password,
  }

  // generate new token
  token, err := s.tokenService.Encode(companyModel)
  if err != nil {
    return nil, err
  }

  // return
  return &v1.UpsertResponse{
    Api:    apiVersion,
    Status: "Success",
    Token:  token,
    // maybe in future add more data to response about the added company.
  }, nil
}

func (s *handler) GetAuth(ctx context.Context, req *v1.UpsertRequest) (*v1.AuthResponse, error) {

  reqToken := req.Token
  // validate the token company and request company
  claims, err := s.tokenService.Decode(reqToken)
  if err != nil {
    return nil, err
  }

  company, err := s.repo.GetById(claims.Company.Id)
  if err != nil {
    return nil, errors.New("Invalid Token")
  }

  out := exportCompanyModel(company)

  return &v1.AuthResponse{
    Api:     apiVersion,
    Status:  "test",
    Company: out,
    // maybe in future add more data to response about the added company.
  }, nil
}

func (s *handler) GetByEmail(ctx context.Context, req *v1.FindRequest) (*v1.FindResponse, error) {

  // fetch company from repo by email
  company, err := s.repo.GetByEmail(req.Email)
  if err != nil {
    return nil, err
  }

  out := exportCompanyModel(company)

  return &v1.FindResponse{
    Api:     apiVersion,
    Status:  "test",
    Company: out,
    // maybe in future add more data to response about the added company.
  }, nil
}

// this func takes database model of Company and exports it to gRPC message model Company
func exportCompanyModel(company *Company) *v1.Company {
  outId := company.Id.Hex()
  out := &v1.Company{
    Id:         outId,
    LastActive: int32(company.LastActive),
    Name:       company.Name,
    Mission:    company.Mission,
    Location:   company.Location,
    Email:      company.Email,
  }
  return out
}

// this func takes a slice of database model of Companys and exports it to gRPC message model Companys
func exportCompanyModels(companys []*Company) []*v1.Company {
  out := []*v1.Company{}
  for _, element := range companys {
    outId := element.Id.Hex()
    company := &v1.Company{
      Id:         outId,
      LastActive: int32(company.LastActive),
      Firstname:  company.Firstname,
      Lastname:   company.Lastname,
      Phone:      company.Phone,
      Email:      company.Email,
    }
    out = append(out, company)
  }
  return out
}
