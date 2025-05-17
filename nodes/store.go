package nodes

import "github.com/jmoiron/sqlx"

func GetAvailableNodes(db *sqlx.DB) ([]Node, error) {
	var nodes []Node
	err := db.Select(&nodes, `SELECT * FROM nodes`)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
