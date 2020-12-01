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
