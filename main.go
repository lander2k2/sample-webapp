package main

import (
	"database/sql"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"
	_ "github.com/lib/pq"
)

var storage Storage

type Storage struct {
	Db *sql.DB
}

type DbParams struct {
	Name     string
	Host     string
	Port     string
	User     string
	Password string
}

type Planet struct {
	Name     string
	Diameter int
	Distance int
	Moons    []Moon
}

type Moon struct {
	Name string
}

func main() {

	flag.Set("logtostderr", "true")
	flag.Parse()

	glog.Info("Starting sample webapp")

	dbPwd := os.Getenv("DB_PASSWORD")
	if dbPwd == "" {
		glog.Fatal("Failed to retrieve database password from the environment")
	}

	dbParams := DbParams{
		Name:     os.Getenv("DB_NAME"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}

	db, err := dbConnect(&dbParams)
	if err != nil {
		glog.Fatalf("Failed to connect to database: %s", err)
	}

	defer db.Close()
	glog.Info("Connected to database")

	storage.Db = db

	http.HandleFunc("/", handler)
	http.ListenAndServe(":8000", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	glog.Infof("Request received from %s", r.RemoteAddr)

	planets, err := storage.getPlanets()
	if err != nil {
		glog.Errorf("Failed to query planets: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	t, err := template.ParseFiles("index.html")
	if err != nil {
		glog.Errorf("Failed to parse template: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	t.Execute(w, &planets)
	glog.Infof("Returned response to %s", r.RemoteAddr)
}

func dbConnect(params *DbParams) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s connect_timeout=10",
		params.Host,
		params.Port,
		params.Name,
		params.User,
		params.Password,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func (storage *Storage) getPlanets() ([]*Planet, error) {

	stmt := `
SELECT planets.name, planets.diameter, planets.distance, moons.name
FROM planets FULL JOIN moons ON planets.id=moons.planet_id`

	rows, err := storage.Db.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	planets := []*Planet{}
	for rows.Next() {
		planet := &Planet{}
		moon := &Moon{}
		hasMoons := true
		err = rows.Scan(&planet.Name, &planet.Diameter, &planet.Distance, &moon.Name)
		if err != nil {
			if strings.Contains(err.Error(), "Scan error") {
				hasMoons = false
			} else {
				return nil, err
			}
		}
		exists := false
		for i, p := range planets {
			if planet.Name == p.Name {
				planets[i].Moons = append(planets[i].Moons, *moon)
				exists = true
			}
		}
		if !exists {
			if hasMoons {
				planet.Moons = append(planet.Moons, *moon)
				planets = append(planets, planet)
			} else {
				planets = append(planets, planet)
			}
		}
	}

	return planets, nil
}
