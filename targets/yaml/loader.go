package yaml

import (
	"github.com/soider/schnur/targets"

	"io/ioutil"
)
import (
	"github.com/ghodss/yaml"
	"github.com/soider/d"
)


type YamlLoader struct {
	Data    []byte
	targets targets.TargetList
	err     error
	file    string
}

func New(filePath string) *YamlLoader {
	l := &YamlLoader{file: filePath}
	return l
}

func (yl *YamlLoader) FromFile(filePath string) {
	data, err := ioutil.ReadFile(filePath)
	yl.Data = data
	if err != nil {
		yl.err = err
		return
	}
	yl.Parse()
}

func (yl *YamlLoader) Parse() {
	yl.targets = make(targets.TargetList, 10)
	d.D(string(yl.Data))
	yl.err = yaml.Unmarshal(yl.Data, &(yl.targets))
}

func (yl *YamlLoader) Targets() (targets.TargetList, error) {
	yl.FromFile(yl.file)
	return yl.targets, yl.err
}
