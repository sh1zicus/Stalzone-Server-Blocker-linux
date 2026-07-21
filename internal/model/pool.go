package model

import "time"

type Tunnel struct {
	Name    string `json:"name"`
	Address string `json:"address"`

	Ping     time.Duration `json:"-"`
	PingOK   bool          `json:"-"`
	Selected bool          `json:"-"`
}

type Pool struct {
	Name    string   `json:"name"`
	Region  string   `json:"region"`
	Tunnels []Tunnel `json:"tunnels"`
}

func (p *Pool) SelectedCount() int {
	n := 0
	for _, t := range p.Tunnels {
		if t.Selected {
			n++
		}
	}
	return n
}

func (p *Pool) IPCount() int {
	return len(p.Tunnels)
}

func (t Tunnel) PingString() string {

	if !t.PingOK {
		return "таймаут"
	}

	return t.Ping.Round(time.Millisecond).String()
}

func (t Tunnel) PingMS() int {
	if !t.PingOK {
		return -1
	}
	return int(t.Ping.Milliseconds())
}
