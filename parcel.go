package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

// Addition of a parcel

func (s ParcelStore) Add(p Parcel) (int, error) {

	// Inserting new row into a table

	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return 0, err
	}

	// Getting a last inserted row id

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return int(id), nil

}

// Getting a single parcel by number

func (s ParcelStore) Get(number int) (Parcel, error) {

	p := Parcel{}

	// Selecting a single row

	row := s.db.QueryRow("SELECT number, client, status, address, created_at FROM parcel WHERE number = :number", sql.Named("number", number))

	// Parsing a selection result

	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	if err != nil {
		return Parcel{}, err
	}

	return p, nil
}

// Getting parcels by client

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {

	var res []Parcel

	// Selecting rows

	rows, err := s.db.Query("SELECT number, client, status, address, created_at FROM parcel WHERE client = :client", sql.Named("client", client))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Iterating through selected rows

	for rows.Next() {

		p := Parcel{}

		// Adding an element to resuled splice

		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

		if err != nil {
			return nil, err
		}

		res = append(res, p)

	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

// Setting a parcel status

func (s ParcelStore) SetStatus(number int, status string) error {

	// Executing db update against a parcel by number
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))

	return err

}

// Setting an address of registered parcel

func (s ParcelStore) SetAddress(number int, address string) error {

	// Executing db update against a parcel by number
	_, err := s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number and status = :status",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	return err

}

// Deletion of a registered parcel

func (s ParcelStore) Delete(number int) error {

	// Executing db delete against a parcel by number
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = :number and status = :status",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered))

	return err

}
