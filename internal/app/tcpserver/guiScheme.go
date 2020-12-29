package tcpserver

// GuiScheme ...
type GuiScheme struct {
	Rooms int    // count of rooms
	Users int    // count of users
	Port  string // addr of listener
}

// newGuiScheme ....
func newGuiScheme(c *Config) *GuiScheme {
	return &GuiScheme{
		Rooms: c.Rooms,
		Users: c.Conns,
		Port:  c.Port,
	}
}
