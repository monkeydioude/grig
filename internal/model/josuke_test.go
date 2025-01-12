package model

import (
	"monkeydioude/grig/internal/errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestICanNOTVerifyAndSanitizeBaseJosukeOnInvalidBaseValues(t *testing.T) {
	j1 := Josuke{
		LogLevel:         "",
		Host:             "",
		Port:             0,
		Store:            "",
		HealthcheckRoute: "",
		Hook:             []Hook{},
		Deployment:       []Deployment{},
		Path:             "",
		FileWriter: func(string, []byte, os.FileMode) error {
			return nil
		},
	}
	assert.ErrorIs(t, j1.VerifyAndSanitize(), errors.ErrModelVerifyInvalidValue)

	j2 := Josuke{
		LogLevel:         "",
		Host:             "salut",
		Port:             -1,
		Store:            "",
		HealthcheckRoute: "",
		Hook:             []Hook{},
		Deployment:       []Deployment{},
		Path:             "",
		FileWriter: func(string, []byte, os.FileMode) error {
			return nil
		},
	}
	assert.ErrorIs(t, j2.VerifyAndSanitize(), errors.ErrModelVerifyInvalidValue)
}

func TestICanNOTVerifyHookOnInvalidValues(t *testing.T) {
	j3 := Josuke{
		LogLevel:         "",
		Host:             "salut",
		Port:             80,
		Store:            "",
		HealthcheckRoute: "",
		Hook: []Hook{
			{
				Name:   "",
				Type:   "",
				Path:   "",
				Secret: "",
			},
		},
		Deployment: []Deployment{},
		Path:       "",
		FileWriter: func(string, []byte, os.FileMode) error {
			return nil
		},
	}
	assert.ErrorIs(t, j3.VerifyAndSanitize(), errors.ErrModelVerifyInvalidValue)
}

func TestICanVerifyHooks(t *testing.T) {
	j4 := Josuke{
		LogLevel:         "",
		Host:             "salut",
		Port:             80,
		Store:            "",
		HealthcheckRoute: "",
		Hook: []Hook{
			{
				Name:   "a1",
				Type:   "b1",
				Path:   "c1",
				Secret: "",
			},
			{
				Name:   "a2",
				Type:   "b2",
				Path:   "c2",
				Secret: "d2",
			},
		},
		Deployment: []Deployment{},
		Path:       "",
		FileWriter: func(string, []byte, os.FileMode) error {
			return nil
		},
	}
	assert.NoError(t, j4.Verify())
}

func TestICanVerifyAWholeJosukeTree(t *testing.T) {
	trial := Josuke{
		LogLevel:         "",
		Host:             "salut",
		Port:             80,
		Store:            "",
		HealthcheckRoute: "",
		Hook: []Hook{
			{
				Name:   "hn1",
				Type:   "ht1",
				Path:   "hp1",
				Secret: "",
			},
			{
				Name:   "hn2",
				Type:   "ht2",
				Path:   "hp2",
				Secret: "hs2",
			},
		},
		Deployment: []Deployment{
			{
				Repo:    "d1_r1",
				ProjDir: "d1_pd1",
				BaseDir: "d1_bd1",
				Branches: []Branch{
					{
						Branch: "d1_br1",
						Actions: []Action{
							{
								Action: "d1_br1_a1",
								Commands: []Command{
									{
										Parts:   []string{"d1_br1_a1_c1", "", "d1_br1_a1_c2"},
										parent:  nil,
										Indexer: Indexer{},
									},
								},
								parent:  nil,
								Indexer: Indexer{},
							},
						},
						parent:  nil,
						Indexer: Indexer{},
					},
					{
						Branch: "d1_br2",
						Actions: []Action{
							{
								Action: "d1_br2_a1",
								Commands: []Command{{
									Parts:   []string{},
									parent:  nil,
									Indexer: Indexer{},
								}},
								parent:  nil,
								Indexer: Indexer{},
							},
						},
						parent:  nil,
						Indexer: Indexer{},
					},
					{
						Branch: "d1_br3",
						Actions: []Action{
							{
								Action: "d1_br3_a1",
								Commands: []Command{
									{
										Parts:   []string{""},
										parent:  nil,
										Indexer: Indexer{},
									},
								},
								parent:  nil,
								Indexer: Indexer{},
							},
							{
								Action: "d1_br3_a2",
								Commands: []Command{
									{
										Parts:   []string{"d1_br3_a2_c1"},
										parent:  nil,
										Indexer: Indexer{},
									},
								},
								parent:  nil,
								Indexer: Indexer{},
							},
						},
						parent:  nil,
						Indexer: Indexer{},
					},
				},
				parent:  nil,
				Indexer: Indexer{},
			},
		},
		Path:       "",
		FileWriter: nil,
	}

	goal := Josuke{
		LogLevel:         "DEBUG",
		Host:             "salut",
		Port:             80,
		Store:            "/tmp",
		HealthcheckRoute: "",
		Hook: []Hook{
			{
				Name:   "hn1",
				Type:   "ht1",
				Path:   "hp1",
				Secret: "",
			},
			{
				Name:   "hn2",
				Type:   "ht2",
				Path:   "hp2",
				Secret: "hs2",
			},
		},
		Deployment: []Deployment{
			{
				Repo:    "d1_r1",
				ProjDir: "d1_pd1",
				BaseDir: "d1_bd1",
				Branches: []Branch{
					{
						Branch: "d1_br1",
						Actions: []Action{
							{
								Action: "d1_br1_a1",
								Commands: []Command{
									{
										Parts:   []string{"d1_br1_a1_c1", "d1_br1_a1_c2"},
										parent:  nil,
										Indexer: Indexer{},
									},
								},
								parent:  nil,
								Indexer: Indexer{},
							},
						},
						parent:  nil,
						Indexer: Indexer{},
					},
					{
						Branch: "d1_br2",
						Actions: []Action{
							{
								Action: "d1_br2_a1",
								Commands: []Command{{
									Parts:   []string{},
									parent:  nil,
									Indexer: Indexer{},
								}},
								parent:  nil,
								Indexer: Indexer{},
							},
						},
						parent:  nil,
						Indexer: Indexer{},
					},
					{
						Branch: "d1_br3",
						Actions: []Action{
							{
								Action: "d1_br3_a1",
								Commands: []Command{
									{
										Parts:   []string{},
										parent:  nil,
										Indexer: Indexer{},
									},
								},
								parent:  nil,
								Indexer: Indexer{},
							},
							{
								Action: "d1_br3_a2",
								Commands: []Command{
									{
										Parts:   []string{"d1_br3_a2_c1"},
										parent:  nil,
										Indexer: Indexer{},
									},
								},
								parent:  nil,
								Indexer: Indexer{},
							},
						},
						parent:  nil,
						Indexer: Indexer{},
					},
				},
				parent:  nil,
				Indexer: Indexer{},
			},
		},
		Path:       "",
		FileWriter: nil,
	}
	assert.NoError(t, trial.VerifyAndSanitize())
	assert.Equal(t, goal, trial)
}
