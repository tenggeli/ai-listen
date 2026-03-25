package model

func (s *MySQLStore) ServiceItems() []*ServiceItem {
	rows, err := s.db.Query(`
		SELECT id, name, category, unit, min_price, max_price, status
		FROM service_items
		WHERE status = 1
		ORDER BY sort ASC, id ASC
	`)
	if err != nil {
		return []*ServiceItem{}
	}
	defer rows.Close()

	var items []*ServiceItem
	for rows.Next() {
		item := &ServiceItem{}
		if err := rows.Scan(&item.ID, &item.Name, &item.Category, &item.Unit, &item.MinPrice, &item.MaxPrice, &item.Status); err == nil {
			items = append(items, item)
		}
	}
	return items
}

func (s *MySQLStore) GetServiceItem(id uint64) (*ServiceItem, error) {
	item := &ServiceItem{}
	err := s.db.QueryRow(`
		SELECT id, name, category, unit, min_price, max_price, status
		FROM service_items
		WHERE id = ?
	`, id).Scan(&item.ID, &item.Name, &item.Category, &item.Unit, &item.MinPrice, &item.MaxPrice, &item.Status)
	if err != nil {
		return nil, ErrServiceItemNotFound
	}
	return item, nil
}
