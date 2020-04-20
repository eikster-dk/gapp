package secrets

import (
	"github.com/eikc/gapp/internal"
	"gopkg.in/yaml.v2"
)

type FileReader interface {
	IsFile(path string) (bool, error)
	ReadFile(path string) ([]byte, error)
	ReadDir(path string) ([]string, error)
}

type parser struct {
	reader FileReader
}

func NewParser(r FileReader) *parser {
	return &parser{
		reader: r,
	}
}

func (p *parser) Parse(path string) (map[string][]internal.Secret, error) {
	isFile, err := p.reader.IsFile(path)
	if err != nil {
		return nil, err
	}

	if isFile {
		secrets, err := p.parseSecret(path)
		if err != nil {
			return nil, err
		}

		sortedSecrets := p.sortSecrets(secrets)
		return sortedSecrets, nil
	} else {
		paths, err := p.reader.ReadDir(path)
		if err != nil {
			return nil, err
		}

		var secretDefinitions []SecretDefinition
		for _, path := range paths {
			ss, err := p.parseSecret(path)
			if err != nil {
				return nil, err
			}

			secretDefinitions = append(secretDefinitions, ss.Secrets...)
		}

		sortedSecrets := p.sortSecrets(YamlDefinition{Secrets: secretDefinitions})

		return sortedSecrets, nil
	}

}

func (p *parser) parseSecret(path string) (YamlDefinition, error) {
	content, err := p.reader.ReadFile(path)
	if err != nil {
		return YamlDefinition{}, err
	}

	var definition YamlDefinition
	err = yaml.Unmarshal(content, &definition)
	if err != nil {
		return YamlDefinition{}, err
	}

	return definition, nil
}

func (p *parser) sortSecrets(definition YamlDefinition) map[string][]internal.Secret {
	mapped := make(map[string][]internal.Secret)
	for _, secret := range definition.Secrets {
		for _, r := range secret.Repos {
			rr := mapped[r]
			rr = append(rr, internal.Secret{
				Name:  secret.Name,
				Value: secret.Value,
			})

			mapped[r] = rr
		}
	}

	return mapped
}
