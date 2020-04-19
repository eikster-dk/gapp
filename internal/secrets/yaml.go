package secrets

type YamlDefinition struct {
	Secrets []SecretDefinition `yaml:"secrets"`
}

type SecretDefinition struct {
	Name  string   `yaml:"name"`
	Value string   `yaml:"value"`
	Repos []string `yaml:"repos"`
}
