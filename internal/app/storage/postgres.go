package storage

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlinaDovbysheva/go-short-url/internal/app"
	"github.com/AlinaDovbysheva/go-short-url/internal/app/util"
	_ "github.com/lib/pq"
	"math/big"
	"strings"

	"github.com/jackc/pgx/v4"
)

type strURLMem struct {
	URL           string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

type strURLMemOut struct {
	URL           string `json:"short_url"`
	CorrelationID string `json:"correlation_id"`
}

type InPostgres struct {
	//mapURL map[string]string
	db *pgx.Conn //sql.DB
}

type URLUid struct {
	UID      string `json:"-"`
	URLShort string `json:"short_url"`
	URL      string `json:"original_url"`
	Deleted  bool   `json:"-"`
}

func NewInPostgre() DBurl {
	// Connect postgres
	/*db, err := sql.Open("postgres", app.DatabaseDsn)
	if err != nil {
		fmt.Println(err)
		//return err  !!! как отсюда вернуть ошибку, если нужно вернуть структуру DBurl?
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxIdleTime(time.Second * 30)
	db.SetConnMaxLifetime(time.Minute * 2)*/

	db, err := pgx.Connect(context.Background(), app.DatabaseDsn)
	if err != nil {
		panic(err)
	}

	// Ping to connection
	err = db.Ping(context.Background())
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
	if _, err := db.Exec(context.Background(), string(BDNew)); err != nil {
		fmt.Println(err)
	}

	return &InPostgres{db}
}

func (m *InPostgres) PingDB() error {
	if err := m.db.Ping(context.Background()); err != nil {
		return err
	}
	return nil
}

func (m *InPostgres) GetAllURLUid(UID string) ([]byte, error) {

	mUID := make([]*URLUid, 0)

	rows, err := m.db.Query(context.Background(), "select t1.url, t1.url_short, t3.user_id from url t1  inner join users_url t2 on t2.url_id = t1.id inner join users t3 on t2.user_id=t3.id where t3.user_id=$1", UID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		bk := new(URLUid)
		err := rows.Scan(&bk.URL, &bk.URLShort, &bk.UID)
		if err != nil {
			return nil, err
		}
		bk.URLShort = app.BaseURL + `/` + bk.URLShort
		mUID = append(mUID, bk)
	}
	defer rows.Close()

	if len(mUID) < 1 {
		return nil, errors.New("urls is absent in db")
	}

	data, _ := json.Marshal(mUID)
	return data, nil
}

func (m *InPostgres) GetURL(shortURL string, UID string) (string, error) {
	mUID := new(URLUid)

	err := m.db.QueryRow(context.Background(), "select ur.url, ur.url_short, u.deleted from users_url u "+
		"INNER JOIN url ur on u.url_id=ur.id "+
		"INNER JOIN users us on u.user_id=us.id "+
		"where ur.url_short=$1 and us.user_id=$2 ", shortURL, UID). //and us.user_id=$2
		Scan(&mUID.URL, &mUID.URLShort, &mUID.Deleted)
	if err != nil {
		err = m.db.QueryRow(context.Background(), "select ur.url, ur.url_short, false deleted from url ur "+
			"where ur.url_short=$1 ", shortURL). //and us.user_id=$2
			Scan(&mUID.URL, &mUID.URLShort, &mUID.Deleted)
		if err != nil {
			return "", util.ErrHandler400
		}
	}
	if mUID.Deleted {
		return "", util.ErrHandler410
	}
	return mUID.URL, nil
}

func (m *InPostgres) PutURL(inputURL string, UID string) (string, []byte, error) {
	var idu int64
	var ids int64
	ctx := context.Background()
	err := m.db.QueryRow(ctx, "select id from users where user_id = $1", UID).Scan(&idu)
	if err != nil {
		err = m.db.QueryRow(ctx, "INSERT INTO users(user_id) VALUES($1) RETURNING id ", UID).Scan(&idu)
		if err != nil {
			fmt.Println("INSERT INTO users= ", err)
			return "", nil, err
		}
	}
	var errExist error
	errExist = nil
	short := ""
	_ = m.db.QueryRow(ctx, "select id,url_short from url where url = $1", inputURL).Scan(&ids, &short)
	fmt.Println("1 Select url_short if exist=", short)
	if short == "" {
		ns, _ := rand.Int(rand.Reader, big.NewInt(10000000)) //util.RandStringBytes(24)
		short = ns.String()
		err = m.db.QueryRow(ctx, "INSERT INTO url(url,url_short)  VALUES($1,$2)  RETURNING id", inputURL, short).Scan(&ids)
		if err != nil {
			fmt.Println("INSERT INTO url(url,url_short)=", inputURL, short, err)
			return "", nil, err
		}
		_, err = m.db.Exec(ctx, "INSERT INTO users_url(user_id,url_id)  VALUES($1,$2) ", idu, ids)
		if err != nil {
			fmt.Println("INSERT INTO users_url(user_id,url_id)= ", err)
			return "", nil, err
		}
	} else {
		errExist = util.ErrHandler409
	}
	d := util.StrtoJSON(app.BaseURL + `/` + short)
	return short, d, errExist
}

func (m *InPostgres) PutURLArray(inputURLJSON []byte, UID string) ([]byte, error) {
	var idu int64
	var ids int64
	var valURL []strURLMem
	var valURLOut []strURLMemOut
	ctx := context.Background()
	if err := json.Unmarshal([]byte(inputURLJSON), &valURL); err != nil {
		panic(err)
	}

	err := m.db.QueryRow(ctx, "select id from users where user_id = $1", UID).Scan(&idu)
	if err != nil {
		err = m.db.QueryRow(ctx, "INSERT INTO users(user_id) VALUES($1) RETURNING id ", UID).Scan(&idu)
		if err != nil {
			fmt.Println("INSERT INTO users= ", err)
			return nil, err
		}
	}

	for _, v := range valURL {
		short := ""
		inputURL := v.URL
		cor := v.CorrelationID
		err = m.db.QueryRow(ctx, "select id,url_short from url where url = $1", inputURL).Scan(&ids, &short)
		fmt.Println("2 Select url_short if exist=", short)
		if err != nil {
			ns, _ := rand.Int(rand.Reader, big.NewInt(10000000)) //util.RandStringBytes(24)
			short = ns.String()
			err = m.db.QueryRow(ctx, "INSERT INTO url(url,url_short)  VALUES($1,$2)  RETURNING id", inputURL, short).Scan(&ids)
			if err != nil {
				fmt.Println("INSERT INTO url(url,url_short)= ", inputURL, short, err)
				return nil, err
			}
			_, err = m.db.Exec(ctx, "INSERT INTO users_url(user_id,url_id)  VALUES($1,$2) ", idu, ids)
			if err != nil {
				fmt.Println("INSERT INTO users_url(user_id,url_id)= ", err)
				return nil, err
			}
		}
		short = app.BaseURL + `/` + short
		valURLOut = append(valURLOut, strURLMemOut{short, cor})
	}
	data, _ := json.Marshal(valURLOut)
	return data, nil
}

func (m *InPostgres) DelURLArray(inputURLJSON []byte, UID string) error {
	var idu int64

	ctx := context.Background()
	err := m.db.QueryRow(ctx, "select id from users where user_id = $1", UID).Scan(&idu)
	if err != nil {
		fmt.Println("User not exists in DB UID="+UID, err)
		return err
	}
	vURL := strings.ReplaceAll(string(inputURLJSON), " ", "")
	vURL = strings.ReplaceAll(strings.ReplaceAll(vURL, "[", ""), "]", "")

	valURL := strings.Split(strings.ReplaceAll(vURL, "\"", ""), ",")
	fmt.Println("Split url short ", valURL)
	if len(valURL) > 20 {
		// batch
		batch := &pgx.Batch{}
		for _, v := range valURL {
			batch.Queue("UPDATE users_url set deleted=true "+
				"where url_id=(select id from url where url_short=$1) and user_id=$2", v, idu)
		}
		br := m.db.SendBatch(context.Background(), batch)

		ct, err := br.Exec()
		if err != nil {
			fmt.Println("Not Updated users_url(user_id,url_id) ", err)
		}
		if ct.RowsAffected() != 1 {
			fmt.Println("ct.RowsAffected()", ct.RowsAffected())
		}
		br.Close()

	} else {
		// обновляем одним запросом, списком
		vURL := strings.ReplaceAll(vURL, "\"", "'")
		query := "UPDATE users_url set deleted=true where " +
			"url_id in (select id from url where url_short in (" + vURL +
			") ) and user_id=$1"
		fmt.Println("query =", query)
		if _, err := m.db.Exec(ctx, query, idu); err != nil {
			fmt.Println("Not Updated users_url(user_id,url_id) ", err)
			return err
		}
	}
	return nil
}

func (m *InPostgres) Close() error { return m.db.Close(context.Background()) }

var BDNew = "-- Database: DB_shortner\n\n-- DROP DATABASE IF EXISTS \"DB_shortner\";\n--create extension pgcrypto;\n\nCREATE TABLE IF NOT EXISTS users(\n                      id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,\n                      user_id uuid NOT NULL,\n                      uname VARCHAR ( 255 ) NULL,\n                      last_login TIMESTAMPTZ NOT NULL DEFAULT NOW()\n);\n\nCREATE TABLE IF NOT EXISTS url (\n                     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,\n                     url VARCHAR ( 255 ) UNIQUE NOT NULL,\n                     url_short VARCHAR ( 255 ) UNIQUE NOT NULL,\n                     created_on TIMESTAMP NOT NULL  DEFAULT NOW()\n);\n\n\nCREATE TABLE IF NOT EXISTS users_url (\n                           url_id INT NOT NULL,\n                           user_id INT NOT NULL, deleted bool not null  DEFAULT  false,\n                           FOREIGN KEY(url_id)\n                               REFERENCES url(id)\n                               ON DELETE SET NULL,\n\n                           FOREIGN KEY (user_id)\n                               REFERENCES users (id)\n);\n\n\n\n"
