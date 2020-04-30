package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

// Problem is db table
type Problem struct {
	Name      string
	Title     string
	Statement string
	Timelimit float64
	Testhash  string
}

// User is db table
type User struct {
	Name     string
	Passhash string
	Admin    bool
}

// Submission is db table
type Submission struct {
	ID          int32
	ProblemName string
	Problem     Problem `gorm:"foreignkey:ProblemName"`
	Lang        string
	Status      string
	Source      string
	Testhash    string
	MaxTime     int32
	MaxMemory   int64
	JudgePing   time.Time
	JudgeName   string
	JudgeTasked bool
	UserName    sql.NullString
	User        User `gorm:"foreignkey:UserName"`
}

// SubmissionTestcaseResult is db table
type SubmissionTestcaseResult struct {
	Submission int32
	Testcase   string
	Status     string
	Time       int32
	Memory     int64
}

// Task is db table
type Task struct {
	Submission int32
	Priority   int32
}

func fetchSubmission(id int32) (Submission, error) {
	sub := Submission{}
	if err := db.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("name")
		}).
		Preload("Problem", func(db *gorm.DB) *gorm.DB {
			return db.Select("name, title, testhash")
		}).
		Where("id = ?", id).First(&sub).Error; err != nil {
		return Submission{}, errors.New("Submission fetch failed")
	}
	return sub, nil
}

func updateSubmissionRegistration(id int32, judgeName string, register bool) (bool, error) {
	tx := db.Begin()
	sub := Submission{}
	if err := tx.First(&sub, id).Error; err != nil {
		tx.Rollback()
		log.Print(err)
		return false, errors.New("Submission fetch failed")
	}
	registered := sub.JudgeName != "" && sub.JudgePing.Add(time.Minute).After(time.Now())
	if registered && sub.JudgeName != judgeName {
		tx.Rollback()
		return false, tx.Rollback().Error
	}
	myself := registered && sub.JudgeName == judgeName
	if !register && !myself {
		return false, tx.Rollback().Error
	}
	if err := tx.Model(&sub).Updates(map[string]interface{}{
		"judge_name": judgeName,
		"judge_ping": time.Now(),
	}).Error; err != nil {
		tx.Rollback()
		log.Print(err)
		return false, errors.New("Submission update failed")
	}
	if err := tx.Commit().Error; err != nil {
		log.Print(err)
		return false, errors.New("Transaction commit failed")
	}
	return true, nil
}

func dbConnect() *gorm.DB {
	host := getEnv("POSTGRE_HOST", "127.0.0.1")
	port := getEnv("POSTGRE_PORT", "5432")
	user := getEnv("POSTGRE_USER", "postgres")
	pass := getEnv("POSTGRE_PASS", "passwd")

	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=librarychecker password=%s sslmode=disable",
		host, port, user, pass)
	log.Printf("Try connect %s", connStr)
	for i := 0; i < 3; i++ {
		db, err := gorm.Open("postgres", connStr)
		if err != nil {
			log.Printf("Cannot connect db %d/3", i)
			time.Sleep(5 * time.Second)
			continue
		}
		db.AutoMigrate(Problem{})
		db.AutoMigrate(User{})
		db.AutoMigrate(Submission{})
		db.AutoMigrate(SubmissionTestcaseResult{})
		db.AutoMigrate(Task{})
		return db
	}
	log.Fatal("Cannot connect db 3 times")
	return nil
}
