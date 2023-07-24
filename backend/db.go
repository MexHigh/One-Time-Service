package main

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"sync"
	"time"
)

type DB struct {
	Path  string
	Mutex sync.Mutex
}

func NewDB(path string) *DB {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		// init file
		_, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		os.WriteFile(path, []byte(`{ "macros": {}, "tokens": {} }`), 0644)
	}

	return &DB{
		Path: path,
	}
}

func (db *DB) readAndLock() (*DBContent, error) {
	db.Mutex.Lock()

	jsonFile, err := os.Open(db.Path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	jsonBytes, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var file DBContent
	if err := json.Unmarshal(jsonBytes, &file); err != nil {
		return nil, err
	}

	return &file, nil
}

func (db *DB) save(state *DBContent) error {
	jsonBytes, err := json.Marshal(state)
	if err != nil {
		return err
	}

	if err := os.WriteFile(db.Path, jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}

// claim aquires or waits for a lock on the DB file, reads
// its contents and exposes it via the *DBContent fuction paramter.
// Afterwrads, its flushes and closes the DB file, releasing
// the mutex lock again.
func (db *DB) claim(f func(dbc *DBContent)) error {
	dbc, err := db.readAndLock()
	defer db.Mutex.Unlock()
	if err != nil {
		return err
	}

	f(dbc)

	if err := db.save(dbc); err != nil {
		return err
	}

	return nil
}

type DBContent struct {
	Macros map[string]*ServiceCall  `json:"macros"`
	Tokens map[string]*TokenDetails `json:"tokens"`
}

type ServiceCall struct {
	Service string      `json:"service"`
	Data    interface{} `json:"data,omitempty"`
}

type TokenDetails struct {
	MacroName string     `json:"macro_name"`
	Created   time.Time  `json:"created"`
	Expires   *time.Time `json:"expires"`
	Comment   string     `json:"comment,omitempty"`
}

func (db *DB) GetMacro(name string) (sc *ServiceCall, err error) {
	err = db.claim(func(dbc *DBContent) {
		sc = dbc.Macros[name]
	})
	return
}

func (db *DB) GetMacroNames() (out []string, err error) {
	err = db.claim(func(dbc *DBContent) {
		for name := range dbc.Macros {
			out = append(out, name)
		}
	})
	return
}

func (db *DB) AddMacro(name string, sc *ServiceCall) (err error) {
	err = db.claim(func(dbc *DBContent) {
		dbc.Macros[name] = sc
	})
	return
}

func (db *DB) DeleteMacro(name string) (err error) {
	err = db.claim(func(dbc *DBContent) {
		delete(dbc.Macros, name)
	})
	return
}

func (db *DB) GetTokenDetails(token string) (td *TokenDetails, err error) {
	err = db.claim(func(dbc *DBContent) {
		td = dbc.Tokens[token]
	})
	return
}

func (db *DB) GetTokens() (tokens map[string]*TokenDetails, err error) {
	err = db.claim(func(dbc *DBContent) {
		tokens = dbc.Tokens
	})
	return
}

func (db *DB) AddToken(token string, td *TokenDetails) (err error) {
	err = db.claim(func(dbc *DBContent) {
		dbc.Tokens[token] = td
	})
	return
}

func (db *DB) DeleteToken(token string) (err error) {
	err = db.claim(func(dbc *DBContent) {
		delete(dbc.Tokens, token)
	})
	return
}
