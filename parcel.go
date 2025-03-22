package main

import (
	"database/sql"
	"fmt"
	"strconv"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка

	p := Parcel{}

	// заполните объект Parcel данными из таблицы
	var client, status, address, created_at string
	row := s.db.QueryRow("SELECT * FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&number, &client, &status, &address, &created_at)
	if err != nil {
		fmt.Println(err)
		return p, err
	}

	client_int, err := strconv.Atoi(client)
	if err != nil {
		return p, err
	}

	p = Parcel{
		Number:    number,
		Client:    client_int,
		Status:    status,
		Address:   address,
		CreatedAt: created_at,
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = :client", sql.Named("client", client))

	if err != nil {
		fmt.Println(err)
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		var number, client, status, address, created_at string

		err := rows.Scan(&number, &client, &status, &address, &created_at)
		if err != nil {
			fmt.Println(err)
			return res, err
		}

		number_int, err := strconv.Atoi(number)
		if err != nil {
			return res, err
		}

		client_int, err := strconv.Atoi(client)
		if err != nil {
			return res, err
		}

		parcel := Parcel{
			Number:    number_int,
			Client:    client_int,
			Status:    status,
			Address:   address,
			CreatedAt: created_at,
		}

		res = append(res, parcel)
	}

	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return res, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = :status WHERE number = :number",
		sql.Named("status", status),
		sql.Named("number", number))
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered

	var status string

	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&status)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if status == "registered" {
		_, err = s.db.Exec("UPDATE parcel SET address = :address WHERE number = :number",
			sql.Named("address", address),
			sql.Named("number", number))
		if err != nil {
			fmt.Println(err)
		}
	} else {
		//fmt.Println("Смена адреса возможна только для посылок со статусом 'Зарегистрирован'.")
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered

	var status string

	row := s.db.QueryRow("SELECT status FROM parcel WHERE number = :number", sql.Named("number", number))
	err := row.Scan(&status)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if status == "registered" {
		_, err = s.db.Exec("DELETE FROM parcel WHERE number = :number",
			sql.Named("number", number))
		if err != nil {
			fmt.Println(err)
		}
	} else {
		//fmt.Println("Удаление строки возможно только для посылок со статусом 'Зарегистрирован'.")
	}

	return nil
}
