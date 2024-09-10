package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		// init folder and file
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			panic(err)
		}
		_, err := os.Create(path)
		if err != nil {
			panic(err)
		}
		os.WriteFile(path, []byte(`{ "service_calls": {}, "tokens": {} }`), 0644)
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
	jsonBytes, err := json.MarshalIndent(state, "", "    ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(db.Path, jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}

// claim aquires or waits for a lock on the DB file, reads its
// contents and exposes it via the *DBContent function parameter.
// Afterwrads, its flushes and closes the DB file, releasing
// the mutex lock again.
func (db *DB) claim(write bool, f func(dbc *DBContent) error) error {
	dbc, err := db.readAndLock()
	defer db.Mutex.Unlock()
	if err != nil {
		return err
	}

	if err := f(dbc); err != nil {
		return err
	}

	if write {
		if err := db.save(dbc); err != nil {
			return err
		}
	}

	return nil
}

type DBContent struct {
	ServiceCalls map[string]*ServiceCall  `json:"service_calls"`
	Tokens       map[string]*TokenDetails `json:"tokens"`
}

type ServiceCall struct {
	Service string      `json:"service"`
	Data    interface{} `json:"data"`
}

type TokenDetails struct {
	ServiceCallName string     `json:"service_call_name"`
	Created         *time.Time `json:"created"`
	Expires         *time.Time `json:"expires"`
	UsesMax         int        `json:"uses_max"`
	UsesLeft        int        `json:"uses_left"`
	Comment         *string    `json:"comment"`
}

/// SERVICE CALL FUNCTIONS ///

func (db *DB) GetServiceCallNames() (out []string, err error) {
	err = db.claim(false, func(dbc *DBContent) error {
		for name := range dbc.ServiceCalls {
			out = append(out, name)
		}
		return nil
	})
	return
}

func (db *DB) GetServiceCall(name string) (sc *ServiceCall, err error) {
	err = db.claim(false, func(dbc *DBContent) error {
		scTemp, ok := dbc.ServiceCalls[name]
		if !ok {
			return fmt.Errorf("service call '%s' does not exist", name)
		}
		sc = scTemp
		return nil
	})
	return
}

func (db *DB) AddServiceCall(name string, sc *ServiceCall) (err error) {
	err = db.claim(true, func(dbc *DBContent) error {
		dbc.ServiceCalls[name] = sc
		return nil
	})
	return
}

func (db *DB) DeleteServiceCall(name string) (err error) {
	err = db.claim(true, func(dbc *DBContent) error {
		delete(dbc.ServiceCalls, name)
		return nil
	})
	return
}

/// TOKEN FUNCTIONS ///

func (db *DB) GetTokenNames() (tokenNames []string, err error) {
	err = db.claim(false, func(dbc *DBContent) error {
		for name := range dbc.Tokens {
			tokenNames = append(tokenNames, name)
		}
		return nil
	})
	return
}

func (db *DB) GetAllTokenDetails() (td map[string]*TokenDetails, err error) {
	err = db.claim(false, func(dbc *DBContent) error {
		td = dbc.Tokens
		return nil
	})
	return
}

func (db *DB) GetTokenDetails(token string) (td *TokenDetails, err error) {
	err = db.claim(false, func(dbc *DBContent) error {
		tdTemp, ok := dbc.Tokens[token]
		if !ok {
			return fmt.Errorf("token '%s' does not exist", token)
		}
		td = tdTemp
		return nil
	})
	return
}

func (db *DB) GetTokensByServiceCallName(serviceCallName string) (tokenNames []string, err error) {
	err = db.claim(false, func(dbc *DBContent) error {
		for name, details := range dbc.Tokens {
			if details.ServiceCallName == serviceCallName {
				tokenNames = append(tokenNames, name)
			}
		}
		return nil
	})
	return
}

func (db *DB) AddToken(token string, td *TokenDetails) (err error) {
	err = db.claim(true, func(dbc *DBContent) error {
		dbc.Tokens[token] = td
		return nil
	})
	return
}

func (db *DB) SetUseCountForToken(token string, newUseCount int) (err error) {
	err = db.claim(true, func(dbc *DBContent) error {
		_, ok := dbc.Tokens[token]
		if !ok {
			return fmt.Errorf("token '%s' does not exist", token)
		}
		// directly access the token to preserve references
		dbc.Tokens[token].UsesLeft = newUseCount
		return nil
	})
	return
}

func (db *DB) DecrementUseCountForToken(token string) (err error) {
	err = db.claim(true, func(dbc *DBContent) error {
		tempToken, ok := dbc.Tokens[token]
		if !ok {
			return fmt.Errorf("token '%s' does not exist", token)
		}
		if tempToken.UsesLeft <= 0 {
			return fmt.Errorf("usage limit of token '%s' exceeded", token)
		}
		// directly access the token to preserve references
		dbc.Tokens[token].UsesLeft -= 1
		return nil
	})
	return
}

func (db *DB) ReplenishUseCountForToken(token string) (err error) {
	err = db.claim(true, func(dbc *DBContent) error {
		tempToken, ok := dbc.Tokens[token]
		if !ok {
			return fmt.Errorf("token '%s' does not exist", token)
		}
		// directly access the token to preserve references
		dbc.Tokens[token].UsesLeft = tempToken.UsesMax
		return nil
	})
	return
}

func (db *DB) DeleteToken(token string) (err error) {
	err = db.claim(true, func(dbc *DBContent) error {
		delete(dbc.Tokens, token)
		return nil
	})
	return
}
