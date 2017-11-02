// Copyright 2017 Elliott Polk. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package config

import (
	"encoding/base64"

	"github.com/elliottpolk/confgr/pgp"
	"github.com/elliottpolk/confgr/uuid"

	"github.com/pkg/errors"
)

//  Encrypt will take the provided token, convert it to a base64 string, attempt
//  to encrypt the config.Value, and return the final version of the token used
//  for encryption. If no token is provided, a UUID (v4) will be generated.
func (cfg *Config) Encrypt(t string) (string, error) {
	token := t
	if len(token) < 1 {
		if token = uuid.GetV4(); len(token) < 1 {
			return "", errors.New("UUID produced an empty string")
		}
	}

	token = base64.StdEncoding.EncodeToString([]byte(token))

	cypher, err := pgp.Encrypt([]byte(token), []byte(cfg.Value))
	if err != nil {
		return token, errors.Wrap(err, "unable to encrypt value with token")
	}

	cfg.Value = string(cypher)
	return token, nil
}

//  Decrypt will take the provided token and attempt to decrypt the config.Value.
//  No additional processing on the token will occur.
func (cfg *Config) Decrypt(t string) error {
	txt, err := pgp.Decrypt([]byte(t), []byte(cfg.Value))
	if err != nil {
		return errors.Wrap(err, "unable to decrypt value with provided token")
	}

	cfg.Value = string(txt)
	return nil
}