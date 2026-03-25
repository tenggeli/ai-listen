package model

func (s *MemoryStore) ServiceItems() []*ServiceItem {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*ServiceItem, 0, len(s.serviceItems))
	for _, item := range s.serviceItems {
		copyItem := *item
		result = append(result, &copyItem)
	}
	return result
}

func (s *MemoryStore) GetServiceItem(id uint64) (*ServiceItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.serviceItems[id]
	if !ok {
		return nil, ErrServiceItemNotFound
	}
	copyItem := *item
	return &copyItem, nil
}
