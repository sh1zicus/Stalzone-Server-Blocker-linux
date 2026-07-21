package model

import "strings"

type RowType int

const (
	RowPool RowType = iota
	RowTunnel
)

type Row struct {
	Type RowType

	Pool   int
	Tunnel int // -1 для Pool
}

type State struct {
	Pools []Pool

	Rows []Row

	Cursor int

	Expanded []bool

	Filter string
}

func NewState(pools []Pool) *State {

	s := &State{
		Pools: pools,
		Expanded: make([]bool, len(pools)),
	}

	for i := range s.Expanded {
		s.Expanded[i] = false
	}

	s.Rebuild()

	return s
}

func (s *State) Rebuild() {

	s.Rows = s.Rows[:0]

	filter := strings.ToLower(strings.TrimSpace(s.Filter))

	for pi, p := range s.Pools {

		showPool := filter == ""

		if filter != "" {

			if strings.Contains(strings.ToLower(p.Name), filter) {
				showPool = true
			} else {

				for _, t := range p.Tunnels {

					if strings.Contains(strings.ToLower(t.Name), filter) {
						showPool = true
						break
					}
				}
			}
		}

		if !showPool {
			continue
		}

		s.Rows = append(s.Rows, Row{
			Type:   RowPool,
			Pool:   pi,
			Tunnel: -1,
		})

		if !s.Expanded[pi] {
			continue
		}

		for ti, t := range p.Tunnels {

			if filter != "" &&
				!strings.Contains(strings.ToLower(t.Name), filter) &&
				!strings.Contains(strings.ToLower(p.Name), filter) {
					continue
				}

				s.Rows = append(s.Rows, Row{
					Type:   RowTunnel,
					Pool:   pi,
					Tunnel: ti,
				})
		}
	}

	if len(s.Rows) == 0 {
		s.Cursor = 0
		return
	}

	if s.Cursor >= len(s.Rows) {
		s.Cursor = len(s.Rows) - 1
	}
}

func (s *State) MoveUp() {

	if s.Cursor > 0 {
		s.Cursor--
	}
}

func (s *State) MoveDown() {

	if s.Cursor < len(s.Rows)-1 {
		s.Cursor++
	}
}

func (s *State) Current() *Row {

	if len(s.Rows) == 0 {
		return nil
	}

	return &s.Rows[s.Cursor]
}

func (s *State) Toggle() {

	row := s.Current()

	if row == nil {
		return
	}

	if row.Type == RowPool {

		s.Expanded[row.Pool] = !s.Expanded[row.Pool]
		s.Rebuild()
		return
	}

	t := &s.Pools[row.Pool].Tunnels[row.Tunnel]

	t.Selected = !t.Selected
}

func (s *State) TogglePool() {

	row := s.Current()
	if row == nil || row.Type != RowPool {
		return
	}

	p := &s.Pools[row.Pool]

	all := true

	for _, t := range p.Tunnels {
		if !t.Selected {
			all = false
			break
		}
	}

	for i := range p.Tunnels {
		p.Tunnels[i].Selected = !all
	}
}
