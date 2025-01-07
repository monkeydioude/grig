package model

type Branch struct {
	Branch  string   `json:"branch"`
	Actions []Action `json:"actions"`
	parent  IndexBuilder
	Indexer
}

func (Branch) GetName() string {
	return "branches"
}

func (c *Branch) SetParent(p IndexBuilder) {
	c.parent = p
}

func (c Branch) GetParent() IndexBuilder {
	return c.parent
}

func (c *Branch) FillBaseData() {
	c.Actions = make([]Action, 1)
	c.Actions[0].FillBaseData()
	c.Actions[0].SetParent(c)
}

func (c *Branch) InitParent() {
	c.SetParent((&Deployment{}))
}

func NewBranch(index int) *Branch {
	br := Branch{}
	br.SetIndex(index)
	br.FillBaseData()
	br.InitParent()
	return &br
}
