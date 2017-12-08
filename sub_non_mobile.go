// +build !mobile

package centrifuge

// History allows to extract channel history.
func (s *Sub) History(skip, limit int) ([]Message, int, error) {
	return s.history(skip, limit)
}

// Presence allows to extract channel history.
func (s *Sub) Presence() (map[string]ClientInfo, error) {
	return s.presence()
}
