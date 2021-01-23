package chat

type ChatMsg interface {
	GetMsg() string
}

func NewStatusMsg(i int, msg string) StatusMsg {
	return StatusMsg{i: i, msg: msg}
}

type StatusMsg struct {
	i   int
	msg string
}

func (sm StatusMsg) GetMsg() string { return sm.msg }

type StreamingStatusMsg struct {
	i   int
	msg string
}

func (sm StreamingStatusMsg) GetMsg() string { return sm.msg }

type ListeningStatusMsg struct {
	i   string
	msg string
}

func (sm ListeningStatusMsg) GetMsg() string { return sm.msg }

type StatusComplexMsg struct {
	i   string
	msg string
}

func (sm StatusComplexMsg) GetMsg() string { return sm.msg }
