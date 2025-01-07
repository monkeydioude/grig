package model

type Action struct {
	Action   string    `json:"action"`
	Commands []Command `json:"commands"`
	parent   IndexBuilder
	Indexer
}

func (Action) GetName() string {
	return "actions"
}

func (c *Action) SetParent(p IndexBuilder) {
	c.parent = p
}

func (c Action) GetParent() IndexBuilder {
	return c.parent
}

func (c *Action) FillBaseData() {
	c.Commands = make([]Command, 1)
	c.Commands[0].SetParent(c)
}

func (c *Action) InitParent() {
	br := &Branch{}
	br.InitParent()
	c.SetParent(br)
}

func NewAction(index int) *Action {
	act := Action{}
	act.SetIndex(index)
	act.FillBaseData()
	act.InitParent()
	return &act
}
