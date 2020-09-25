package session

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	epsagonawswrapper "github.com/epsagon/epsagon-go/wrappers/aws/aws-sdk-go/aws"
)

func NewDefaultAwsSession(cfgs ...*aws.Config) *session.Session {
	return epsagonawswrapper.WrapSession(
		session.Must(session.NewSession(cfgs...)),
	)
}
