package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

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
	println("Start")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	server := os.Getenv("SAGE_SERVER")
	db := os.Getenv("SAGE_DB")
	user := os.Getenv("SAGE_USER")
	password := os.Getenv("SAGE_PASS")
	port, err := strconv.ParseInt(os.Getenv("SAGE_PORT"), 0, 64)
	if err != nil {
		log.Fatal("SAGE_PORT not in .env: ", err)
	}

	mysql_server := os.Getenv("MYSQL_SERVER")
	mysql_port, err := strconv.ParseInt(os.Getenv("MYSQL_PORT"), 0, 64)
	if err != nil {
		log.Fatal("MYSQL_PORT not in .env: ", err)
	}
	mysql_user := os.Getenv("MYSQL_USER")
	mysql_password := os.Getenv("MYSQL_PASS")
	mysql_db := os.Getenv("MYSQL_DB")

	// Read all products from sage
	sage_connstring := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s;port=%d", server, db, user, password, port)
	sage_query := "select sg_auf_artikel.SG_AUF_ARTIKEL_PK, sg_auf_artikel.ARTNR, sg_auf_artikel.SUCHBEGRIFF, sg_auf_artikel.ZUSTEXT1, sg_auf_vkpreis.PR01 FROM sg_auf_artikel INNER JOIN sg_auf_vkpreis ON sg_auf_artikel.SG_AUF_ARTIKEL_PK = sg_auf_vkpreis.SG_AUF_ARTIKEL_FK"

	conn, err := sql.Open("sqlserver", sage_connstring)
	if err != nil {
		log.Fatalf("Open Sage Connection failed: %s", err.Error())
	}

	defer conn.Close()
	println("Read Sage")
	rows, err := conn.Query(sage_query)
	if err != nil {
		log.Fatalf("Sage Query failed: %s", err.Error())
	}
	defer rows.Close()
	println("Sort Sage")
	var Sage []Artikel
	for rows.Next() {
		var Id sql.NullInt64
		var Artikelnummer sql.NullString
		var Artikelname sql.NullString
		var Specs sql.NullString
		var Preis sql.NullFloat64

		if err := rows.Scan(&Id, &Artikelnummer, &Artikelname, &Specs, &Preis); err != nil {
			log.Fatalf("Scan Sage Row failed: %s", err.Error())
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

	mysql_connstring := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", mysql_user, mysql_password, mysql_server, mysql_port, mysql_db)
	mysql_conn, err := sql.Open("mysql", mysql_connstring)
	if err != nil {
		log.Fatal(err)
	}
	defer mysql_conn.Close()
	println("Read MySQL")
	if len(Sage) > 0 {
		for i := range Sage {
			println("Write MySQL with id:", i)
			println("Max::", len(Sage))
			id := Sage[i].Id
			nummer := Sage[i].Artikelnummer
			name := strings.ReplaceAll(Sage[i].Artikelname, "'", "\"")
			spec := strings.ReplaceAll(Sage[i].Specs, "'", "\"")
			price := Sage[i].Preis
			upsert_query := fmt.Sprintf("INSERT INTO Aussteller (id, Artikelnummer, Artikelname, Specs, Preis) VALUES (%d, '%s', '%s','%s','%.2f') ON DUPLICATE KEY UPDATE Artikelname = '%s', Specs = '%s', Preis = '%.2f'", id, nummer, name, spec, price, name, spec, price)
			_, err := mysql_conn.Exec(upsert_query)
			if err != nil {
				log.Print(upsert_query)
				log.Fatal(err)
			}
		}
	}
}
