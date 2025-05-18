package network

import "github.com/jmoiron/sqlx"

func StoreNetwork(db *sqlx.Tx, network *Network) error {
	_, err := db.NamedExec(`INSERT INTO networks (id, interface_name, ip_address, port, project_id) VALUES (:id, :interface_name, :ip_address, :port, :project_id)`, network)
	if err != nil {
		return err
	}

	return nil
}
