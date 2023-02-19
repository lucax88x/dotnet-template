package yaml

import (
	"dzor/core"

	"gopkg.in/yaml.v3"
)

func Serialize() {
	buf, err := yaml.Marshal(r.Config)
	if err != nil {
		return err
	}
}
