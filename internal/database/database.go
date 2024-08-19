package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type Database struct {
	path string
	mu   *sync.Mutex
}

type Structure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users  map[int]User  `json:"users"`
}

type Chirp struct {
	ID       int    `json:"id"`
	Body     string `json:"body"`
	AuthorID int    `json:"author_id"`
}

type User struct {
	ID           int          `json:"id"`
	Email        string       `json:"email"`
	Password     []byte       `json:"password"`
	RefreshToken RefreshToken `json:"refresh_token"`
}

type RefreshToken struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

func NewDatabase(path string) (*Database, error) {
	db := &Database{
		path: path,
		mu:   &sync.Mutex{},
	}

	if err := db.createDBIfNotExists(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) createDBIfNotExists() error {
	_, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		structure := Structure{
			Chirps: map[int]Chirp{},
			Users:  map[int]User{},
		}
		return db.writeDB(structure)
	}
	return err
}

func (db *Database) writeDB(structure Structure) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	data, err := json.Marshal(structure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, data, 0600)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) loadDB() (Structure, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	structure := Structure{}
	data, err := os.ReadFile(db.path)
	if errors.Is(err, os.ErrNotExist) {
		return structure, err
	}

	err = json.Unmarshal(data, &structure)
	if err != nil {
		return structure, err
	}

	return structure, nil
}

func (db *Database) CreateChirp(body string, authorID int) (Chirp, error) {
	chirpsMap, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	id := len(chirpsMap.Chirps) + 1
	chirp := Chirp{
		ID:       id,
		Body:     body,
		AuthorID: authorID,
	}

	chirpsMap.Chirps[id] = chirp
	err = db.writeDB(chirpsMap)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *Database) GetChirps() ([]Chirp, error) {
	chirpsMap, err := db.loadDB()
	if err != nil {
		return nil, err
	}

	chirps := make([]Chirp, 0, len(chirpsMap.Chirps))
	for _, chirp := range chirpsMap.Chirps {
		chirps = append(chirps, chirp)
	}

	return chirps, nil
}

func (db *Database) CreateUser(email string, password []byte) (User, error) {
	usersMap, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range usersMap.Users {
		if user.Email == email {
			return User{}, errors.New("user already exists")
		}

	}

	id := len(usersMap.Users) + 1
	user := User{
		ID:       id,
		Email:    email,
		Password: password,
	}

	usersMap.Users[id] = user
	err = db.writeDB(usersMap)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (db *Database) GetUserByEmail(email string) (User, error) {
	usersMap, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	for _, user := range usersMap.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user not found")
}

func (db *Database) UpdateUser(id int, email string, password []byte) error {
	usersMap, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := usersMap.Users[id]
	if !ok {
		return errors.New("user not found")
	}

	user.Email = email
	user.Password = password

	usersMap.Users[id] = user

	err = db.writeDB(usersMap)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) AddRefreshToken(id int, refreshToken RefreshToken) error {
	usersMap, err := db.loadDB()
	if err != nil {
		return err
	}

	user, ok := usersMap.Users[id]
	if !ok {
		return errors.New("user not found")
	}

	user.RefreshToken = refreshToken

	usersMap.Users[id] = user

	err = db.writeDB(usersMap)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetRefreshToken(token string) (RefreshToken, User, error) {
	usersMap, err := db.loadDB()
	if err != nil {
		return RefreshToken{}, User{}, err
	}

	for _, user := range usersMap.Users {
		if user.RefreshToken.Token == token {
			return user.RefreshToken, user, nil
		}
	}

	return RefreshToken{}, User{}, errors.New("refresh token not found")
}

func (db *Database) RevokeRefreshToken(token string) error {
	usersMap, err := db.loadDB()
	if err != nil {
		return err
	}

	for _, user := range usersMap.Users {
		if user.RefreshToken.Token == token {
			user.RefreshToken = RefreshToken{}
			usersMap.Users[user.ID] = user
			err = db.writeDB(usersMap)
			if err != nil {
				return err
			}
			return nil
		}
	}

	return errors.New("refresh token not found")
}
