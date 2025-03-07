package parser

import (
	"errors"

	occfg "github.com/opencloud-eu/opencloud/pkg/config"
	"github.com/opencloud-eu/opencloud/pkg/shared"
	"github.com/opencloud-eu/opencloud/services/idp/pkg/config"
	"github.com/opencloud-eu/opencloud/services/idp/pkg/config/defaults"

	"github.com/opencloud-eu/opencloud/pkg/config/envdecode"
)

// ParseConfig loads configuration from known paths.
func ParseConfig(cfg *config.Config) error {
	err := occfg.BindSourcesToStructs(cfg.Service.Name, cfg)
	if err != nil {
		return err
	}

	defaults.EnsureDefaults(cfg)

	// load all env variables relevant to the config in the current context.
	if err := envdecode.Decode(cfg); err != nil {
		// no environment variable set for this config is an expected "error"
		if !errors.Is(err, envdecode.ErrNoTargetFieldsAreSet) {
			return err
		}
	}

	defaults.Sanitize(cfg)

	return Validate(cfg)
}

func Validate(cfg *config.Config) error {
	switch cfg.IDP.IdentityManager {
	case "cs3":
		if cfg.MachineAuthAPIKey == "" {
			return shared.MissingMachineAuthApiKeyError(cfg.Service.Name)
		}
	case "ldap":
		if cfg.Ldap.BindPassword == "" {
			return shared.MissingLDAPBindPassword(cfg.Service.Name)
		}
	}

	return nil
}
