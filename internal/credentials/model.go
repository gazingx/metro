package credentials

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/razorpay/metro/internal/common"
	"github.com/razorpay/metro/pkg/encryption"
	"github.com/sethvargo/go-password/password"
)

const (
	// Prefix for all credential keys in the registry
	Prefix = "credentials/"

	// PartSeparator used in constructing the username parts
	PartSeparator = "__"

	// perf10__656f81 : sample username format
	usernameFormat = "%v%v%v"
)

var (
	// CtxKey identified credentials will be populated in ctx with this key
	CtxKey = contextKey("CredentialsCtxKey")
)

type contextKey string

func (c contextKey) String() string {
	return string(c)
}

// Model for a credential
type Model struct {
	common.BaseModel
	Username  string `json:"username"`
	Password  string `json:"password"`
	ProjectID string `json:"project_id"`
	// in future, this model can contain some ACL as well
}

// ICredentials defines getters on a credential
type ICredentials interface {
	GetUsername() string
	GetPassword() string
	GetProjectID() string
}

// Key returns the key for storing credentials in the registry
func (m *Model) Key() string {
	return m.Prefix() + m.Username
}

// Prefix returns the Key prefix
func (m *Model) Prefix() string {
	return common.GetBasePrefix() + Prefix + m.ProjectID + "/"
}

// NewCredential returns a new credential model
func NewCredential(username, password string) ICredentials {
	m := &Model{}
	m.Username = username
	// store Password only after encrypting it
	m.Password, _ = encryption.EncryptAsHexString([]byte(password))
	// extract projectID from username
	m.ProjectID = GetProjectIDFromUsername(username)
	return m
}

// GetUsername returns the credential Username
func (m *Model) GetUsername() string {
	return m.Username
}

// GetPassword returns the decrypted credential Password
func (m *Model) GetPassword() string {
	// decrypt before reading Password
	pwd, _ := encryption.DecryptFromHexString(m.Password)
	return string(pwd)
}

// GetProjectID returns the credential projectID
func (m *Model) GetProjectID() string {
	return m.ProjectID
}

func newUsername(projectID string) string {
	return fmt.Sprintf(usernameFormat, projectID, PartSeparator, uuid.New().String()[:6])
}

func newPassword() string {
	pwd, _ := password.Generate(20, 10, 0, false, true)
	encryptedPwd, _ := encryption.EncryptAsHexString([]byte(pwd))
	return encryptedPwd
}

// GetProjectIDFromUsername returns the projectID for a given username
// extract only if the username was generated by metro, else return an empty string
func GetProjectIDFromUsername(username string) string {
	if usernameRegex.MatchString(username) {
		return strings.Split(username, PartSeparator)[0]
	}
	return ""
}
