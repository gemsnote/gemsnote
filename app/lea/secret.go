package lea

import "errors"

// ShippedAppSecret is the development value that used to ship in conf/app.conf.
// It is public knowledge, so it must never sign sessions or PDF credentials in
// a real deployment.
const ShippedAppSecret = "V85ZzBeTnzpsHyjQX4zukbQ8qqtju9y2aDM55VWxAH9Qop19poekx3xkcDVvrD0y"

// ValidateAppSecret reports whether secret can be used to sign credentials.
// Production deployments fail closed; development only warns so a fresh
// checkout still runs.
func ValidateAppSecret(secret string, production bool) error {
	if secret == "" {
		return errors.New("app.secret is empty; sessions would be unsigned")
	}
	if secret == ShippedAppSecret {
		return errors.New("app.secret still uses the public default value shipped in conf/app.conf")
	}
	if production && len(secret) < 32 {
		return errors.New("app.secret is shorter than 32 characters")
	}
	return nil
}
