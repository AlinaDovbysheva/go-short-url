package storage

import (
	"database/sql"
	"fmt"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
	_ "github.com/lib/pq"
	"time"

	"encoding/json"
	"errors"
	"github.com/AlinaDovbysheva/go-short-url/internal/app"
)

type strURLMem struct {
	URL            string `json:"original_url"`
	Correlation_id string `json:"correlation_id"`
}

type strURLMemOut struct {
	URL            string `json:"short_url"`
	Correlation_id string `json:"correlation_id"`
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

func NewInPostgre() DBurl {
	// Connect postgres
	db, err := sql.Open("postgres", app.DatabaseDsn)
	if err != nil {
		fmt.Println(err)
		//return err  !!! как отсюда вернуть ошибку, если нужно вернуть структуру DBurl?
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Minute * 2)

	// Ping to connection
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	/*
		// read DB structure from file
		dir, _ := filepath.Abs(filepath.Dir(os.Args[1]))
		path := filepath.Join(dir, "dbshortnerPG.sql")
		query, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println(err)
		}*/
	if _, err := db.Exec(string(BDNew)); err != nil {
		fmt.Println(err)
	}

	return &InPostgres{db}
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
		bk.URLShort = app.BaseURL + `/` + bk.URLShort
		mUid = append(mUid, bk)
	}

	defer rows.Close()
	data, _ := json.Marshal(mUid)
	fmt.Println("GetAllURLUid data= ", string(data))
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

	short := util.RandStringBytes(9)

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
			fmt.Println("INSERT INTO url(url,url_short)= %s , %s ", inputURL, short, err)
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

func (m *InPostgres) PutURLArray(inputURLJSON []byte, UID string) ([]byte, error) {
	var idu int64
	var ids int64
	var valUrl []strURLMem
	var valUrlOut []strURLMemOut

	if err := json.Unmarshal([]byte(inputURLJSON), &valUrl); err != nil {
		panic(err)
	}

	err := m.db.QueryRow("select id from users where user_id = $1", UID).Scan(&idu)
	if err != nil {
		err = m.db.QueryRow("INSERT INTO users(user_id) VALUES($1) RETURNING id ", UID).Scan(&idu)
		if err != nil {
			fmt.Println("INSERT INTO users= ", err)
			return nil, err
		}
	}

	for _, v := range valUrl {
		short := util.RandStringBytes(9)
		inputURL := v.URL
		cor := v.Correlation_id
		err = m.db.QueryRow("select id,url_short from url where url = $1", inputURL).Scan(&ids, &short)
		if err != nil {
			err = m.db.QueryRow("INSERT INTO url(url,url_short)  VALUES($1,$2)  RETURNING id", inputURL, short).Scan(&ids)
			if err != nil {
				fmt.Println("INSERT INTO url(url,url_short)= %s , %s ", inputURL, short, err)
				return nil, err
			}
			_, err = m.db.Exec("INSERT INTO users_url(user_id,url_id)  VALUES($1,$2) ", idu, ids)
			if err != nil {
				fmt.Println("INSERT INTO users_url(user_id,url_id)= ", err)
				return nil, err
			}
		}
		short = app.BaseURL + `/` + short
		valUrlOut = append(valUrlOut, strURLMemOut{short, cor})
	}
	data, _ := json.Marshal(valUrlOut)
	return data, nil
}

func (m *InPostgres) Close() error { return m.db.Close() }

var BDNew = "-- Database: DB_shortner\n\n-- DROP DATABASE IF EXISTS \"DB_shortner\";\n--create extension pgcrypto;\n\nCREATE TABLE IF NOT EXISTS users(\n                      id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,\n                      user_id uuid NOT NULL,\n                      uname VARCHAR ( 255 ) NULL,\n                      last_login TIMESTAMPTZ NOT NULL DEFAULT NOW()\n);\n\nCREATE TABLE IF NOT EXISTS url (\n                     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,\n                     url VARCHAR ( 255 ) UNIQUE NOT NULL,\n                     url_short VARCHAR ( 255 ) UNIQUE NOT NULL,\n                     created_on TIMESTAMP NOT NULL  DEFAULT NOW()\n);\n\n\nCREATE TABLE IF NOT EXISTS users_url (\n                           url_id INT NOT NULL,\n                           user_id INT NOT NULL,\n                           FOREIGN KEY(url_id)\n                               REFERENCES url(id)\n                               ON DELETE SET NULL,\n\n                           FOREIGN KEY (user_id)\n                               REFERENCES users (id)\n);\n\n\n\n"
