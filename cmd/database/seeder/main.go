package main

import (
	"fmt"
	"log"
	"math/rand"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.con/albugowy15/api-double-track/internal/pkg/config"
	"github.con/albugowy15/api-double-track/internal/pkg/utils"
)

type Admin struct {
	Username    string `db:"username"`
	Password    string `db:"password"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phone_number"`
	SchoolId    string `db:"school_id"`
}

type School struct {
	Name string `db:"name"`
}

type Alternative struct {
	Alternative string `db:"alternative"`
	Description string `db:"description"`
}

type Question struct {
	Code        string `db:"code"`
	Question    string `db:"question"`
	Category    string `db:"category"`
	Description string `db:"description"`
	Number      int    `db:"number"`
}

type Student struct {
	Username    string `db:"username"`
	Password    string `db:"password"`
	Fullname    string `db:"fullname"`
	Nisn        string `db:"nisn"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phone_number"`
	SchoolId    string `db:"school_id"`
}

func main() {
	config.LoadConfig(".")
	conf := config.GetConfig()
	connStr := fmt.Sprintf("dbname=%s host=%s port=%s user=%s password=%s sslmode=%s", conf.DbName, conf.DbHost, conf.DbPort, conf.DbUser, conf.DbPass, conf.DbSsl)

	db, err := sqlx.Connect(conf.DbDriver, connStr)
	if err != nil {
		log.Fatal(err)
	}

	// seed schools
	schools := []School{
		{Name: "SMA IPIEMS Surabaya"},
		{Name: "SMA Dharmawanita Surabaya"},
		{Name: "SMA Negeri 1 Ngadirojo"},
		{Name: "SMA Negeri 1 Jenangan"},
		{Name: "SMA Negeri 1 Gondanglegi"},
		{Name: "SMA Negeri 1 Balongpanggang"},
		{Name: "SMA Negeri 1 Turen"},
		{Name: "SMA Negeri 1 Sumbermanjing"},
		{Name: "SMA Negeri 1 Pulung"},
	}
	_, err = db.NamedExec(`INSERT INTO schools (name) VALUES (:name)`, schools)
	if err != nil {
		log.Fatalf("error insert schools: %v", err)
	}

	tx := db.MustBegin()

	// seed alternatives
	alteratives := []Alternative{
		{Alternative: "Multimedia", Description: "Alternative keterampilan Multimedia"},
		{Alternative: "Teknik Elektro", Description: "Alternative keterampilan Teknik Elektro"},
		{Alternative: "Teknik Listrik", Description: "Alternative keterampilan Tata Listrik"},
		{Alternative: "Tata Boga", Description: "Alternative keterampilan Tata Boga"},
		{Alternative: "Tata Busana", Description: "Alternative keterampilan Tata Busana"},
		{Alternative: "Tata Kecantikan", Description: "Alternative keterampilan Tata Kecantikan"},
		{Alternative: "Teknik Kendararaan Ringan/Motor", Description: "Alternative keterampilan Teknik Kendararaan Ringan/Motor"},
	}
	_, err = tx.NamedExec(`INSERT INTO alternatives (alternative, description) VALUES (:alternative, :description)`, alteratives)
	if err != nil {
		log.Fatalf("error insert alternatives: %v", err)
	}

	// questions
	questions := []Question{
		{Code: "MIN_MUL", Number: 1, Question: "Seberapa besar minatmu terhadap multimedia?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa besar minatmu terhadap keterampilan Multimedia. 1=Sangat Tidak Berminat, 4=Sangat Berminat"},
		{Code: "MIN_TEL", Number: 2, Question: "Seberapa besar minatmu terhadap teknik elektro?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa besar minatmu terhadap keterampilan teknik elektro. 1=Sangat Tidak Berminat, 4=Sangat Berminat"},
		{Code: "MIN_TLI", Number: 3, Question: "Seberapa besar minatmu terhadap teknik listrik?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa besar minatmu terhadap keterampilan teknik listrik. 1=Sangat Tidak Berminat, 4=Sangat Berminat"},
		{Code: "MIN_TBU", Number: 4, Question: "Seberapa besar minatmu terhadap tata busana?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa besar minatmu terhadap keterampilan tata busana. 1=Sangat Tidak Berminat, 4=Sangat Berminat"},
		{Code: "MIN_TBO", Number: 5, Question: "Seberapa besar minatmu terhadap tata boga?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa besar minatmu terhadap keterampilan tata boga. 1=Sangat Tidak Berminat, 4=Sangat Berminat"},
		{Code: "MIN_TKE", Number: 6, Question: "Seberapa besar minatmu terhadap tata kecantikan?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa besar minatmu terhadap keterampilan tata kecantikan. 1=Sangat Tidak Berminat, 4=Sangat Berminat"},
		{Code: "MIN_TKR", Number: 7, Question: "Seberapa besar minatmu terhadap teknik kendaraan ringan/motor?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa besar minatmu terhadap keterampilan teknik kendaraan ringan/motor. 1=Sangat Tidak Berminat, 4=Sangat Berminat"},
		{Code: "FAS_MUL", Number: 8, Question: "Seberapa dukungan fasilitas yang kamu miliki untuk bisa mengikuti keterampilan Multimedia?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa mendukung fasilitas yang kamu miliki terhadap keterampilan Multimedia. 1=Sangat Tidak Mendukung, 4=Sangat Mendukung"},
		{Code: "FAS_TEL", Number: 9, Question: "Seberapa dukungan fasilitas yang kamu miliki untuk bisa mengikuti keterampilan Teknik Elektro?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa mendukung fasilitas yang kamu miliki terhadap keterampilan Teknik Elektro. 1=Sangat Tidak Mendukung, 4=Sangat Mendukung"},
		{Code: "FAS_TLI", Number: 10, Question: "Seberapa dukungan fasilitas yang kamu miliki untuk bisa mengikuti keterampilan Teknik Lisrik?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa mendukung fasilitas yang kamu miliki terhadap keterampilan Teknik Listrik. 1=Sangat Tidak Mendukung, 4=Sangat Mendukung"},
		{Code: "FAS_TBU", Number: 11, Question: "Seberapa dukungan fasilitas yang kamu miliki untuk bisa mengikuti keterampilan Tata Busana?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa mendukung fasilitas yang kamu miliki terhadap keterampilan Tata Busana. 1=Sangat Tidak Mendukung, 4=Sangat Mendukung"},
		{Code: "FAS_TBO", Number: 12, Question: "Seberapa dukungan fasilitas yang kamu miliki untuk bisa mengikuti keterampilan Tata Boga?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa mendukung fasilitas yang kamu miliki terhadap keterampilan Tata Boga. 1=Sangat Tidak Mendukung, 4=Sangat Mendukung"},
		{Code: "FAS_TKE", Number: 13, Question: "Seberapa dukungan fasilitas yang kamu miliki untuk bisa mengikuti keterampilan Tata Kecantikan", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa mendukung fasilitas yang kamu miliki terhadap keterampilan Tata Kecantika. 1=Sangat Tidak Mendukung, 4=Sangat Mendukung"},
		{Code: "FAS_TKR", Number: 14, Question: "Seberapa dukungan fasilitas yang kamu miliki untuk bisa mengikuti keterampilan Teknik Kendararaan Ringan/Motor?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa mendukung fasilitas yang kamu miliki terhadap keterampilan Teknik Kendararaan Ringan/Motor. 1=Sangat Tidak Mendukung, 4=Sangat Mendukung"},
		{Code: "JLP_GAJ", Number: 15, Question: "Dalam memilih bidang keterampilan, seberapa pentingnya jumlah lapangan pekerjaan dengan gaji?", Category: "COMPARISON", Description: "Dari 9 s.d 1/9 dimana 1 s.d 9 menunjukkan jumlah lapangan pekerjaan lebih penting dan 1 s.d 1/9 menunjukkan gaji lebih penting. Sedangkan 1 menunjukkan keduanya sama-sama penting"},
		{Code: "JLP_PEW", Number: 16, Question: "Dalam memilih bidang keterampilan, seberapa pentingnya jumlah lapangan pekerjaan dengan peluang wirausaha?", Category: "COMPARISON", Description: "Dari 9 s.d 1/9 dimana 1 s.d 9 menunjukkan jumlah lapangan pekerjaan lebih penting dan 1 s.d 1/9 menunjukkan peluang wirausaha lebih penting. Sedangkan 1 menunjukkan keduanya sama-sama penting"},
		{Code: "JLP_MIN", Number: 17, Question: "Dalam memilih bidang keterampilan, seberapa pentingnya jumlah lapangan pekerjaan dengan minatmu?", Category: "COMPARISON", Description: "Dari 9 s.d 1/9 dimana 1 s.d 9 menunjukkan jumlah lapangan pekerjaan lebih penting dan 1 s.d 1/9 menunjukkan minat lebih penting. Sedangkan 1 menunjukkan keduanya sama-sama penting"},
		{Code: "JLP_FAS", Number: 18, Question: "Dalam memilih bidang keterampilan, seberapa pentingnya jumlah lapangan pekerjaan dengan dukungan fasilitas yang kamu miliki?", Category: "COMPARISON", Description: "Dari 9 s.d 1/9 dimana 1 s.d 9 menunjukkan jumlah lapangan pekerjaan lebih penting dan 1 s.d 1/9 menunjukkan dukungan fasilitas lebih penting. Sedangkan 1 menunjukkan keduanya sama-sama penting"},
		{Code: "GAJ_PEW", Number: 19, Question: "Dalam memilih bidang keterampilan, seberapa pentingnya gaji dengan peluang wirausaha?", Category: "COMPARISON", Description: "Dari 9 s.d 1/9 dimana 1 s.d 9 menunjukkan gaji lebih penting dan 1 s.d 1/9 menunjukkan peluang wirausaha lebih penting. Sedangkan 1 menunjukkan keduanya sama-sama penting"},
		{Code: "GAJ_MIN", Number: 20, Question: "Dalam memilih bidang keterampilan, seberapa pentingnya gaji dengan minatmu?", Category: "COMPARISON", Description: "Dari 9 s.d 1/9 dimana 1 s.d 9 menunjukkan gaji lebih penting dan 1 s.d 1/9 menunjukkan minat lebih penting. Sedangkan 1 menunjukkan keduanya sama-sama penting"},
		{Code: "GAJ_FAS", Number: 21, Question: "Dalam memilih bidang keterampilan, seberapa pentingnya gaji dengan dukungan fasilitas yang kamu miliki?", Category: "COMPARISON", Description: "Dari 9 s.d 1/9 dimana 1 s.d 9 menunjukkan gaji lebih penting dan 1 s.d 1/9 menunjukkan dukungan fasilitas lebih penting. Sedangkan 1 menunjukkan keduanya sama-sama penting"},
		{Code: "PEW_MIN", Number: 22, Question: "Dalam memilih bidang keterampilan, seberapa pentingnya peluang wirausaha dengan minatmu?", Category: "COMPARISON", Description: "Dari 9 s.d 1/9 dimana 1 s.d 9 menunjukkan peluang wirausaha lebih penting dan 1 s.d 1/9 menunjukkan minat lebih penting. Sedangkan 1 menunjukkan keduanya sama-sama penting"},
		{Code: "PEW_FAS", Number: 23, Question: "Dalam memilih bidang keterampilan, seberapa pentingnya peluang wirausaha dengan dukungan fasilitas yang kamu miliki?", Category: "COMPARISON", Description: "Dari 9 s.d 1/9 dimana 1 s.d 9 menunjukkan peluang wirausaha lebih penting dan 1 s.d 1/9 menunjukkan dukungan fasilitas lebih penting. Sedangkan 1 menunjukkan keduanya sama-sama penting"},
		{Code: "MIN_FAS", Number: 24, Question: "Dalam memilih bidang keterampilan, seberapa pentingnya minat dengan dukungan fasilitas yang kamu miliki?", Category: "COMPARISON", Description: "Dari 9 s.d 1/9 dimana 1 s.d 9 menunjukkan minat lebih penting dan 1 s.d 1/9 menunjukkan dukungan fasilitas lebih penting. Sedangkan 1 menunjukkan keduanya sama-sama penting"},
	}
	_, err = tx.NamedExec(`INSERT INTO questions (code, number, question, category, description) VALUES (:code, :number, :question, :category, :description)`, questions)
	if err != nil {
		log.Fatalf("error insert questions: %v", err)
	}

	type SchooldId struct {
		Id string `db:"id"`
	}
	schoolIds := []SchooldId{}
	db.Select(&schoolIds, "SELECT id from schools LIMIT 3")
	// students
	studentPass, err := utils.HashStr("passwordStudent")
	if err != nil {
		log.Fatalf("error hashing pass: %v", err)
	}
	students := []Student{
		{
			Username:    "bughowy",
			Email:       "bughowy@gmail.com",
			Password:    studentPass,
			Fullname:    "Mohamad Kholid Bughowi",
			Nisn:        "123252532",
			PhoneNumber: "086345193034",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
		{
			Username:    "ahmadthoriq",
			Email:       "ahmadthoriq@gmail.com",
			Password:    studentPass,
			Fullname:    "Ahmad Thoriq",
			Nisn:        "12773663943",
			PhoneNumber: "084495723723",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
		{
			Username:    "dwiassegaf",
			Email:       "dwiassegaf@gmail.com",
			Password:    studentPass,
			Fullname:    "Dwi Assegaf",
			Nisn:        "1344263626",
			PhoneNumber: "076653626642",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
		{
			Username:    "efendi",
			Email:       "efendimanik@gmail.com",
			Password:    studentPass,
			Fullname:    "Efendi Manik",
			Nisn:        "1283774334",
			PhoneNumber: "08423173732",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
		{
			Username:    "gedesam",
			Email:       "gedesam@gmail.com",
			Password:    studentPass,
			Fullname:    "Gede Samudra",
			Nisn:        "1264378882",
			PhoneNumber: "0816763434999",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
		{
			Username:    "gatotbakti",
			Email:       "gatotbakti@gmail.com",
			Password:    studentPass,
			Fullname:    "Gatot Surbakti",
			Nisn:        "74346364343",
			PhoneNumber: "0857285926385",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
		{
			Username:    "harissag",
			Email:       "harissag@gmail.com",
			Password:    studentPass,
			Fullname:    "Haris Saragih",
			Nisn:        "12665577557",
			PhoneNumber: "0865523558834",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
	}
	_, err = tx.NamedExec("INSERT INTO students (username, email, password, fullname, nisn, phone_number, school_id) VALUES (:username, :email, :password, :fullname, :nisn, :phone_number, :school_id)", students)
	if err != nil {
		log.Fatalf("error insert students: %v", err)
	}

	// seed admins
	adminPass, err := utils.HashStr("passwordAdmin")
	if err != nil {
		log.Fatalf("error hashing pass: %v", err)
	}
	admins := []Admin{
		{Username: "admintester1", Password: adminPass, Email: "admintester1@gmail.com", PhoneNumber: "087542845123", SchoolId: schoolIds[rand.Intn(3)].Id},
		{Username: "admintester2", Password: adminPass, Email: "admintester2@gmail.com", PhoneNumber: "085166371256", SchoolId: schoolIds[rand.Intn(3)].Id},
		{Username: "admintester3", Password: adminPass, Email: "admintester3@gmail.com", PhoneNumber: "085441327327", SchoolId: schoolIds[rand.Intn(3)].Id},
	}
	_, err = tx.NamedExec(`INSERT INTO admins (username, password, email, phone_number, school_id) VALUES (:username, :password, :email, :phone_number, :school_id)`, admins)
	if err != nil {
		log.Fatalf("error insert admins: %v", err)
	}
	tx.Commit()

	defer db.Close()
}
