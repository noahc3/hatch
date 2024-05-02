package types

import (
	"encoding/json"
	"time"
)

type Variable struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	EnvVariable  string `json:"env_variable"`
	DefaultValue string `json:"default_value"`
	UserViewable bool   `json:"user_viewable"`
	UserEditable bool   `json:"user_editable"`
	Rules        string `json:"rules"`
	Value        string `json:"-"`
}

type Config struct {
	Files   string `json:"files"`
	Startup string `json:"startup"`
	Logs    string `json:"logs"`
	Stop    string `json:"stop"`
}

type Script struct {
	Script     string `json:"script"`
	Container  string `json:"container"`
	Entrypoint string `json:"entrypoint"`
}

type Scripts struct {
	Installation Script `json:"installation"`
}

type Meta struct {
	Version   string `json:"version"`
	UpdateUrl string `json:"update_url"`
}

type CrackProps struct {
	InstallDirPath    string            `json:"-"`
	InstallScriptPath string            `json:"-"`
	EnvVariables      map[string]string `json:"-"`
}

type Egg struct {
	Comment      string            `json:"_comment"`
	Meta         Meta              `json:"meta"`
	ExportedAt   time.Time         `json:"exported_at"`
	Name         string            `json:"name"`
	Author       string            `json:"author"`
	Description  string            `json:"description"`
	Features     []string          `json:"features"`
	DockerImages map[string]string `json:"docker_images"`
	FileDenyList []string          `json:"file_denylist"`
	Startup      string            `json:"startup"`
	Config       Config            `json:"config"`
	Scripts      Scripts           `json:"scripts"`
	Variables    []Variable        `json:"variables"`
	CrackProps   CrackProps        `json:"-"`
}

func (e *Egg) UnmarshalJSON(data []byte) error {
	type Alias Egg

	aux := &struct {
		*Alias
		ExportedAt string `json:"exported_at"`
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	exportedAt, err := time.Parse("2006-01-02T15:04:05Z0700", aux.ExportedAt)
	if err != nil {
		return err
	}

	e.ExportedAt = exportedAt
	return nil
}
