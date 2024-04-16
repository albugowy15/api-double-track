package seeds

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type Alternative struct {
	Alternative string `db:"alternative"`
	Description string `db:"description"`
}

func SeedAlternativeTx(tx *sqlx.Tx) {
	// seed alternatives
	alteratives := []Alternative{
		{Alternative: "Multimedia", Description: "Alternative keterampilan Multimedia"},
		{Alternative: "Teknik Elektro", Description: "Alternative keterampilan Teknik Elektro"},
		{Alternative: "Teknik Listrik", Description: "Alternative keterampilan Tata Listrik"},
		{Alternative: "Tata Busana", Description: "Alternative keterampilan Tata Busana"},
		{Alternative: "Tata Boga", Description: "Alternative keterampilan Tata Boga"},
		{Alternative: "Tata Kecantikan", Description: "Alternative keterampilan Tata Kecantikan"},
		{Alternative: "Teknik Kendararaan Ringan/Motor", Description: "Alternative keterampilan Teknik Kendararaan Ringan/Motor"},
	}
	_, err := tx.NamedExec(`INSERT INTO alternatives (alternative, description) VALUES (:alternative, :description)`, alteratives)
	if err != nil {
		log.Fatalf("error insert alternatives: %v", err)
	}
}
