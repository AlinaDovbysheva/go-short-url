package storage

import (
	"database/sql"
	"fmt"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
	_ "github.com/lib/pq"
	"io/ioutil"
	"os"
	"path/filepath"

	"encoding/json"
	"errors"
	"github.com/AlinaDovbysheva/go-short-url/internal/app"
)

func NewInPostgre() DBurl {
	// Connect postgres
	db, err := sql.Open("postgres", app.DatabaseDsn)
	if err != nil {
		fmt.Println(err)
		//return err  !!! как отсюда вернуть ошибку, если нужно вернуть структуру DBurl?
	}
	// Ping to connection
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	dir, _ := filepath.Abs(filepath.Dir(os.Args[1]))
	path := filepath.Join(dir, "dbshortnerPG.sql")
	query, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := db.Exec(string(query)); err != nil {
		fmt.Println(err)
	}

	return &InPostgres{db}
}

type InPostgres struct {
	//mapURL map[string]string
	db *sql.DB
}

type URLUid struct {
	Uid      string `json:"-"`
	URLShort string `json:"short_url"`
	URL      string `json:"original_url"`
}

func (m *InPostgres) PingDB() error {
	if err := m.db.Ping(); err != nil {
		return err
	}
	return nil
}

func (m *InPostgres) GetAllURLUid(UID string) ([]byte, error) {

	mUid := make([]*URLUid, 0)

	rows, err := m.db.Query("select t1.url, t1.url_short, t3.user_id from url t1  inner join users_url t2 on t2.url_id = t1.id inner join users t3 on t2.user_id=t3.id where t3.user_id=$1", UID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		bk := new(URLUid)
		err := rows.Scan(&bk.URL, &bk.URLShort, &bk.Uid)
		if err != nil {
			return nil, err
		}
		mUid = append(mUid, bk)
	}

	defer rows.Close()
	data, _ := json.Marshal(mUid)
	fmt.Println("GetAllURLUid data= ", mUid)
	return data, nil
}

func (m *InPostgres) GetURL(shortURL string) (string, error) {
	mUid := new(URLUid)

	err := m.db.QueryRow("select url, url_short from url t where t.url_short=$1", shortURL).
		Scan(&mUid.URL, &mUid.URLShort)
	if err != nil {
		return "", err
	}

	if mUid.URL == "" {
		return "", errors.New("id is absent in db")
	}
	return mUid.URL, nil
}

func (m *InPostgres) PutURL(inputURL string, UID string) (string, error) {

	short := util.RandStringBytes(7)
	var idu int64
	var ids int64
	err := m.db.QueryRow("select id from users where user_id = $1", UID).Scan(&idu)
	if err != nil {
		err = m.db.QueryRow("INSERT INTO users(user_id) VALUES($1) RETURNING id ", UID).Scan(&idu)
		if err != nil {
			fmt.Println("INSERT INTO users= ", err)
			return "", err
		}
	}
	err = m.db.QueryRow("select id,url_short from url where url = $1", inputURL).Scan(&ids, &short)
	if err != nil {
		err = m.db.QueryRow("INSERT INTO url(url,url_short)  VALUES($1,$2)  RETURNING id", inputURL, short).Scan(&ids)
		if err != nil {
			fmt.Println("INSERT INTO url(url,url_short)= ", err)
			return "", err
		}
		_, err = m.db.Exec("INSERT INTO users_url(user_id,url_id)  VALUES($1,$2) ", idu, ids)
		if err != nil {
			fmt.Println("INSERT INTO users_url(user_id,url_id)= ", err)
			return "", err
		}
	}

	return short, nil
}

func (m *InPostgres) Close() error { return m.db.Close() }
