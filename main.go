package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	fmt.Println("before")
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println("after")
	fmt.Println(q)
	return nil
}

type Execution struct {
	tableName struct{} `pg:"execution, discard_unknown_columns"`

	id           int64         `pg:"id, pk"`
	partnerName  string        `pg:"partner, fk:name"`
	async        bool          `pg:"is_async`
	timestamp    time.Time     `pg:"timestamp, default:now()"`
	loadStatuses []*LoadStatus `pg:"rel:has-many, join_fk:fk_execution"`
}

func NewExecution(partnerName string, async bool, timestamp time.Time) *Execution {
	return &Execution{
		partnerName: partnerName,
		async:       async,
		timestamp:   timestamp,
	}
}

func (execution *Execution) ID() int64 {
	return execution.id
}

func (execution *Execution) PartnerName() string {
	return execution.partnerName
}

func (execution *Execution) Async() bool {
	return execution.async
}

func (execution *Execution) Timestamp() time.Time {
	return execution.timestamp
}

func (execution *Execution) LoadStatuses() []*LoadStatus {
	return execution.loadStatuses
}

type LoadStatus struct {
	tableName struct{} `pg:"load_status, discard_unknown_columns"`

	id          int64      `pg:"id,pk"`
	event       string     `pg:"event"`
	status      string     `pg:"status"`
	description string     `pg:"status"`
	timestamp   time.Time  `pg:"timestamp, default:now()"`
	executionID int64      `pg:"fk_execution"`
	execution   *Execution `pg:"fk:fk_execution"`
}

func NewLoadStatus(event string, status string, description string, executionID int64, timestamp time.Time) *LoadStatus {
	return &LoadStatus{
		event:       event,
		status:      status,
		description: description,
		executionID: executionID,
		timestamp:   timestamp,
	}
}

func (loadStatus *LoadStatus) ID() int64 {
	return loadStatus.id
}

func (loadStatus *LoadStatus) Event() string {
	return loadStatus.event
}

func (loadStatus *LoadStatus) Status() string {
	return loadStatus.status
}

func (loadStatus *LoadStatus) Description() string {
	return loadStatus.description
}

func (loadStatus *LoadStatus) Timestamp() time.Time {
	return loadStatus.timestamp
}

func (loadStatus *LoadStatus) ExecutionID() int64 {
	return loadStatus.executionID
}

func (loadStatus *LoadStatus) Execution() *Execution {
	return loadStatus.execution
}

func ExampleDB_Model() {
	db := pg.Connect(&pg.Options{
		Database: "try_gopg",
		User:     "postgres",
		Password: "secret",
	})

	defer db.Close()

	db.AddQueryHook(dbLogger{})

	err := createSchema(db)
	if err != nil {
		panic(err)
	}

	execution1 := NewExecution("partnerName 1", true, time.Now())
	fmt.Println(execution1)
	_, err = db.Model(execution1).Insert()
	if err != nil {
		panic(err)
	}
}

// createSchema creates database schema for User and Story models.
func createSchema(db *pg.DB) error {
	var u *Execution
	var s *LoadStatus
	models := []interface{}{
		u,
		s,
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp: false,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	ExampleDB_Model()
}
