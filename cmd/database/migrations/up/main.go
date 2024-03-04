package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.con/albugowy15/api-double-track/internal/pkg/config"
)

var schema = `
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE "admins" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "username" varchar(30) UNIQUE NOT NULL,
  "password" text NOT NULL,
  "email" varchar(100) UNIQUE,
  "phone_number" varchar(15) UNIQUE,
  "school_id" uuid,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);
CREATE TRIGGER set_admins_timestamp
BEFORE UPDATE ON admins
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE "students" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "username" varchar(30) UNIQUE NOT NULL,
  "password" text NOT NULL,
  "fullname" varchar(500) NOT NULL,
  "nisn" varchar(20) UNIQUE NOT NULL,
  "email" varchar(100) UNIQUE,
  "phone_number" varchar(15) UNIQUE,
  "school_id" uuid,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);
CREATE TRIGGER set_students_timestamp
BEFORE UPDATE ON students
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE "schools" (
  "id" uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  "name" varchar(50) NOT NULL,
  "address" text
);

CREATE TABLE "alternatives" (
  "id" SERIAL PRIMARY KEY,
  "alternative" varchar(50) NOT NULL,
  "description" text
);

CREATE TABLE "questionnare_settings" (
  "id" SERIAL PRIMARY KEY,
  "alternative_id" int,
  "school_id" uuid,
  "total_open_jobs" integer,
  "entrepreneurship_opportunity" integer,
  "salary" integer,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);
CREATE TRIGGER set_questionnare_settings_timestamp
BEFORE UPDATE ON questionnare_settings
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE "questions" (
  "id" SERIAL PRIMARY KEY,
  "number" int UNIQUE NOT NULL,
  "code" varchar(20) UNIQUE NOT NULL,
  "question" text NOT NULL,
  "category" varchar(20) NOT NULL,
  "description" text
);

CREATE TABLE "answers" (
  "id" BIGSERIAL PRIMARY KEY,
  "student_id" uuid,
  "question_id" int,
  "answer_str" varchar(20),
  "answer_num" int,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);
CREATE TRIGGER set_answers_timestamp
BEFORE UPDATE ON answers
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE "ahp_to_alternatives" (
  "id" BIGSERIAL PRIMARY KEY,
  "score" decimal NOT NULL,
  "ahp_id" bigint,
  "alternative_id" int,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);
CREATE TRIGGER set_ahp_to_alternatives_timestamp
BEFORE UPDATE ON ahp_to_alternatives
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE "ahp" (
  "id" BIGSERIAL PRIMARY KEY,
  "student_id" uuid,
  "consistency_index" decimal,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);
CREATE TRIGGER set_ahp_timestamp
BEFORE UPDATE ON ahp
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE "topsis_to_alternatives" (
  "id" BIGSERIAL PRIMARY KEY,
  "score" decimal,
  "topsis_id" bigint,
  "alternative_id" int,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);
CREATE TRIGGER set_topsis_to_alternatives_timestamp
BEFORE UPDATE ON topsis_to_alternatives
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();

CREATE TABLE "topsis" (
  "id" BIGSERIAL PRIMARY KEY,
  "student_id" uuid,
  "created_at" TIMESTAMP NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP NOT NULL DEFAULT (now())
);
CREATE TRIGGER set_topsis_timestamp
BEFORE UPDATE ON topsis
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();


ALTER TABLE "students" ADD FOREIGN KEY ("school_id") REFERENCES "schools" ("id");

ALTER TABLE "admins" ADD FOREIGN KEY ("school_id") REFERENCES "schools" ("id");

ALTER TABLE "questionnare_settings" ADD FOREIGN KEY ("alternative_id") REFERENCES "alternatives" ("id");

ALTER TABLE "questionnare_settings" ADD FOREIGN KEY ("school_id") REFERENCES "schools" ("id");

ALTER TABLE "answers" ADD FOREIGN KEY ("student_id") REFERENCES "schools" ("id");

ALTER TABLE "answers" ADD FOREIGN KEY ("question_id") REFERENCES "questions" ("id");

ALTER TABLE "topsis_to_alternatives" ADD FOREIGN KEY ("topsis_id") REFERENCES "topsis" ("id");

ALTER TABLE "topsis_to_alternatives" ADD FOREIGN KEY ("alternative_id") REFERENCES "alternatives" ("id");

ALTER TABLE "ahp_to_alternatives" ADD FOREIGN KEY ("ahp_id") REFERENCES "ahp" ("id");

ALTER TABLE "ahp_to_alternatives" ADD FOREIGN KEY ("alternative_id") REFERENCES "alternatives" ("id");
`

func main() {
	config.LoadConfig(".")
	conf := config.GetConfig()
	connStr := fmt.Sprintf("dbname=%s host=%s port=%s user=%s password=%s sslmode=%s", conf.DbName, conf.DbHost, conf.DbPort, conf.DbUser, conf.DbPass, conf.DbSsl)

	db, err := sqlx.Connect(conf.DbDriver, connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.MustExec(schema)
}
