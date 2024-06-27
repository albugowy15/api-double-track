package seeds

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type Question struct {
	Code        string `db:"code"`
	Question    string `db:"question"`
	Category    string `db:"category"`
	Description string `db:"description"`
	Number      int    `db:"number"`
}

func SeedQuestionsTx(tx *sqlx.Tx) {
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
		{Code: "FAS_TKR", Number: 14, Question: "Seberapa dukungan fasilitas yang kamu miliki untuk bisa mengikuti keterampilan Teknik Kendaraan Ringan/Motor?", Category: "PREFERENCE", Description: "Dari skala 1-4 tentukan seberapa mendukung fasilitas yang kamu miliki terhadap keterampilan Teknik Kendaraan Ringan/Motor. 1=Sangat Tidak Mendukung, 4=Sangat Mendukung"},
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
	_, err := tx.NamedExec(`INSERT INTO questions (code, number, question, category, description) VALUES (:code, :number, :question, :category, :description)`, questions)
	if err != nil {
		log.Fatalf("error insert questions: %v", err)
	}
}
