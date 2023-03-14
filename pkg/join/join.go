package join

import (
	"context"
	connectionPostgresql "ephorservices/pkg/db"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Join struct {
	ConnectionDb *connectionPostgresql.Manager
	RowsData     pgx.Rows
}

func (j *Join) GetJoin(accountId interface{}, sql string) {
	sqlReplace := ""
	lookForSchema := "/schema/"
	schema := strings.Contains(sql, lookForSchema)
	accountSchema := fmt.Sprintf("account%v", accountId)
	lookForMain := "/main/"
	main := strings.Contains(sql, lookForMain)
	if schema == true {
		sqlReplace = strings.Replace(sql, lookForSchema, accountSchema, -1)
	} else if main == true {
		sqlReplace = strings.Replace(sql, lookForMain, "main", -1)
	}
	ctx := context.Background()
	rows, err := j.ConnectionDb.Conn.Query(ctx, sqlReplace)
	if err != nil {
		log.Println(err)
	}
	j.RowsData = rows
	defer rows.Close()
}
