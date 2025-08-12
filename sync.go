package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type Artikel struct {
	Id            int
	Artikelnummer string
	Artikelname   string
	Specs         string
	Preis         float64
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	for {
		fmt.Print("\033[H\033[2J")
		startTime := time.Now()
		start()

		fmt.Printf("Aussteller Aktualisert, dauer: %v Minuten\nWarte 1 Stunde\n", time.Since(startTime).Minutes())
		time.Sleep(1 * time.Hour)
	}
}

func start() {
	fmt.Println("Read Sage Artikel")
	sage_artikel, err := getSageArtikel()
	if err != nil {
		panic(err)
	}
	fmt.Println("Update Aussteller Database")
	err = updateDB(sage_artikel)
	if err != nil {
		panic(err)
	}
}

func getSageArtikel() ([]Artikel, error) {
	connectionString, ok := os.LookupEnv("SAGE_URL")
	if !ok {
		return nil, fmt.Errorf("failed to read SAGE_URL")
	}

	conn, err := sql.Open("sqlserver", connectionString)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := "select sg_auf_artikel.SG_AUF_ARTIKEL_PK, sg_auf_artikel.ARTNR, sg_auf_artikel.SUCHBEGRIFF, sg_auf_artikel.ZUSTEXT1, sg_auf_vkpreis.PR01 FROM sg_auf_artikel INNER JOIN sg_auf_vkpreis ON sg_auf_artikel.SG_AUF_ARTIKEL_PK = sg_auf_vkpreis.SG_AUF_ARTIKEL_FK"

	rows, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var Sage []Artikel

	for rows.Next() {
		var Id sql.NullInt64
		var Artikelnummer sql.NullString
		var Artikelname sql.NullString
		var Specs sql.NullString
		var Preis sql.NullFloat64
		if err := rows.Scan(&Id, &Artikelnummer, &Artikelname, &Specs, &Preis); err != nil {
			return nil, err
		}
		if Id.Valid && Artikelnummer.Valid && Artikelname.Valid && Specs.Valid && Preis.Valid {
			var tmp Artikel
			tmp.Id = int(Id.Int64)
			tmp.Artikelnummer = Artikelnummer.String
			tmp.Artikelname = Artikelname.String
			tmp.Preis = Preis.Float64
			tmp.Specs = Specs.String
			Sage = append(Sage, tmp)
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return Sage, nil
}

func updateDB(artikel []Artikel) error {
	connstring, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return fmt.Errorf("failed to read DATABASE_URL")
	}

	rawQuery := "INSERT INTO Aussteller (id, Artikelnummer, Artikelname, Specs, Preis) VALUES"

	for idx, item := range artikel {
		if idx == len(artikel)-1 {
			rawQuery = fmt.Sprintf("%s (%d, '%s', '%s', '%s', %.2f)", rawQuery, item.Id, item.Artikelnummer, strings.ReplaceAll(item.Artikelname, "'", "\""), strings.ReplaceAll(item.Specs, "'", "\""), item.Preis)
		} else {
			rawQuery = fmt.Sprintf("%s (%d, '%s', '%s', '%s', %.2f),", rawQuery, item.Id, item.Artikelnummer, strings.ReplaceAll(item.Artikelname, "'", "\""), strings.ReplaceAll(item.Specs, "'", "\""), item.Preis)
		}
	}

	rawQuery = fmt.Sprintf("%s ON DUPLICATE KEY UPDATE Artikelnummer = VALUES(Artikelnummer), Artikelname = VALUES(Artikelname), Specs = VALUES(Specs), Preis = VALUES(Preis);", rawQuery)

	conn, err := sql.Open("mysql", connstring)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.Exec(rawQuery)
	return err
}
