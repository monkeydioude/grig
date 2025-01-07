package model

type Deployment struct {
	Repo     string   `json:"repo"`
	ProjDir  string   `json:"proj_dir"`
	BaseDir  string   `json:"base_dir"`
	Branches []Branch `json:"branches"`
	parent   any
	Indexer
}

func (Deployment) GetName() string {
	return "deployment"
}

func (c *Deployment) SetParent(p IndexBuilder) {
	c.parent = nil
}

func (c Deployment) GetParent() IndexBuilder {
	return nil
}

func (c *Deployment) FillBaseData() {
	c.Branches = make([]Branch, 1)
	c.Branches[0].FillBaseData()
	c.Branches[0].SetParent(c)
}
