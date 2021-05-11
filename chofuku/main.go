package chofuku

import (
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"os"
	"path/filepath"
)

type Chofuku struct {
	DB *sql.DB
}
type Duplicate struct {
	Size         int64
	Head100kHash string
	FullHash     string
	Names        []string
}

func New(dir string) (Chofuku, error) {
	//os.Remove("/tmp/chofuku.db")
	//db, err := sql.Open("sqlite3", "/tmp/chofuku.db")
	db, err := sql.Open("sqlite3", "file:chofuku.db?cache=shared&mode=memory")
	if err != nil {
		return Chofuku{}, err
	}
	chofuku := Chofuku{db}
	_, err = db.Exec(`
		CREATE TABLE files (
			name TEXT NOT NULL PRIMARY KEY, 
			size INTEGER NOT NULL,
			head100k_hash TEXT DEFAULT "",
			full_hash TEXT DEFAULT ""
		);
	`)
	if err != nil {
		return chofuku, err
	}
	err = filepath.Walk(dir, func(p string, info os.FileInfo, err error) error {
		IsSymLink := (os.ModeSymlink & info.Mode()) != 0
		if !info.IsDir() && !IsSymLink {
			_, err = db.Exec(`INSERT INTO files (name, size) 
							VALUES (?, ?);`, p, info.Size())
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return chofuku, err
	}
	return chofuku, err
}
func (c *Chofuku) Close() {
	c.DB.Close()
}
func (c *Chofuku) GetDuplicates() ([]Duplicate, error) {
	duplicates := []Duplicate{}
	rows, err := c.DB.Query(`SELECT size, head100k_hash, full_hash, json_group_array(name)
		FROM files GROUP BY size, head100k_hash, full_hash HAVING count(*) > 1;
	`)
	if err != nil {
		fmt.Fprintln(os.Stderr, "group by select error")
		return duplicates, err
	}
	for rows.Next() {
		var size int64
		var head100k_hash, full_hash, json_group_array string
		if err = rows.Scan(&size, &head100k_hash, &full_hash, &json_group_array); err != nil {
			return duplicates, err
		}
		var names []string
		err = json.Unmarshal([]byte(json_group_array), &names)
		if err != nil {
			fmt.Fprintln(os.Stderr, "json parse error")
			return duplicates, err
		}
		duplicates = append(duplicates,
			Duplicate{size,
				head100k_hash,
				full_hash,
				names,
			})
	}
	return duplicates, nil
}
func (c *Chofuku) UpdateHead100k() error {
	duplicates, err := c.GetDuplicates()
	if err != nil {
		return err
	}
	for _, v := range duplicates {
		if v.Size == 0 {
			continue
		}
		for _, n := range v.Names {
			hash, err := read100k(n)
			if err != nil {
				return err
			}
			err = c.updateColumn(n, "head100k_hash", hash)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (c *Chofuku) UpdateFullHash() error {
	duplicates, err := c.GetDuplicates()
	if err != nil {
		return err
	}
	for _, v := range duplicates {
		if v.Size <= 1024*100 {
			continue
		}
		for _, n := range v.Names {
			hash, err := read_all(n)
			if err != nil {
				return err
			}
			err = c.updateColumn(n, "full_hash", hash)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (c *Chofuku) updateColumn(name, column, hash string) error {
	if column != "head100k_hash" && column != "full_hash" {
		return errors.New("invalid column name")
	}
	_, err := c.DB.Exec(`UPDATE files SET `+column+` = ? 
		WHERE name = ?`, hash, name)
	if err != nil {
		return err
	}
	return nil
}
func read100k(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	buf := make([]byte, 100*1024)
	n, err := f.Read(buf)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(buf[:n])), nil
}
func read_all(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
