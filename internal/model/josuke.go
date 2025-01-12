package model

import (
	"encoding/json"
	"fmt"
	"monkeydioude/grig/internal/consts"
	customErrors "monkeydioude/grig/internal/errors"
	"monkeydioude/grig/internal/service/utils"
	"monkeydioude/grig/pkg/errors"
	"monkeydioude/grig/pkg/trans_types"
	"os"
	"slices"
)

type Indexer struct {
	Index int `json:"-"`
}

func (i Indexer) GetIndex() int {
	return i.Index
}

func (i *Indexer) SetIndex(index int) {
	i.Index = index
}

type Josuke struct {
	LogLevel         string                                  `json:"logLevel"`
	Host             string                                  `json:"host"`
	Port             trans_types.StringInt                   `json:"port"`
	Store            string                                  `json:"store"`
	HealthcheckRoute string                                  `json:"healthcheck_route"`
	Hook             []Hook                                  `json:"hook"`
	Deployment       []Deployment                            `json:"deployment"`
	Path             string                                  `json:"-"`
	FileWriter       func(string, []byte, os.FileMode) error `json:"-"`
}

func (j Josuke) Save() error {
	data, err := json.Marshal(j)
	if err != nil {
		return fmt.Errorf("Josuke.Save(): %w: %w", customErrors.ErrMarshaling, err)
	}
	if j.FileWriter == nil {
		return fmt.Errorf("Josuke.Save(): c.FileWriter: %w", customErrors.ErrNilPointer)
	}
	if err := j.FileWriter(j.Path, data, 0644); err != nil {
		return fmt.Errorf("Josuke.Save(): %w: %w", customErrors.ErrWritingFile, err)
	}
	return nil
}

func (j *Josuke) FillBaseData() {
	if j.HealthcheckRoute == "" {
		j.HealthcheckRoute = "/josuke/healthcheck"
	}
	if len(j.Hook) == 0 {
		j.Hook = make([]Hook, 1)
	}
	if len(j.Deployment) == 0 {
		j.Deployment = make([]Deployment, 1)
	}
}

func (j Josuke) Verify() error {
	if j.Port <= 0 {
		return errors.Wrap(customErrors.ErrNilPointer, "Josuke.Verify: Port")
	}
	if j.Host == "" {
		return errors.Wrap(customErrors.ErrModelVerifyInvalidValue, "Josuke.Verify: Host")
	}
	return nil
}

func (j *Josuke) VerifyAndSanitize() error {
	if err := j.Verify(); err != nil {
		return err
	}
	for _, dep := range j.Deployment {
		if err := dep.Verify(); err != nil {
			return err
		}
	}
	if j.LogLevel == "" {
		j.LogLevel = consts.JOSUKE_DEFAULT_LOG_LEVEL
	}
	if j.Store == "" {
		j.Store = utils.GetDefaultTMPDirectory()
	}
	j.Deployment = slices.DeleteFunc(j.Deployment, func(dep Deployment) bool {
		return dep.VerifyAndSanitize() != nil
	})
	return nil
}
