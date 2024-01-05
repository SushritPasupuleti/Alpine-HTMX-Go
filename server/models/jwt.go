package models

import "github.com/golang-jwt/jwt/v4"

type JWTClaims struct {
	*jwt.RegisteredClaims
	Email       string      `json:"email,omitempty"`
	AppMetadata AppMetadata `json:"app_metadata,omitempty"`
	Issuer      string      `json:"iss,omitempty"`
	Subject     string      `json:"sub,omitempty"`
	Audience    string      `json:"aud,omitempty"`
	Expiration  int64       `json:"exp,omitempty"`
	IssuedAt    int64       `json:"iat,omitempty"`
	JWTID       string      `json:"jti,omitempty"`
}

type AppMetadata struct {
	Authorization Authorization `json:"authorization,omitempty"`
}

type Authorization struct {
	Roles []string `json:"roles,omitempty"`
}

// Create a new JWTClaims object
func (jwtClaims *JWTClaims) AddAppMetadata(appMetadata AppMetadata) {
	jwtClaims.AppMetadata = appMetadata
}

// Add a new Authorization object to the JWTClaims object
func (jwtClaims *JWTClaims) AddAuthorization(authorization Authorization) {
	jwtClaims.AppMetadata.Authorization = authorization
}

// Add a new Roles array to the Authorization object
func (jwtClaims *JWTClaims) AddRoles(roles []string) {
	jwtClaims.AppMetadata.Authorization.Roles = roles
}

// Add a new Role to the Roles array in the Authorization object
func (jwtClaims *JWTClaims) AddRole(role string) {
	jwtClaims.AppMetadata.Authorization.Roles = append(jwtClaims.AppMetadata.Authorization.Roles, role)
}
