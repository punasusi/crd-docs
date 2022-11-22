package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path"

	_ "github.com/lib/pq"

	"github.com/crdsdev/doc/pkg/crd"
	"github.com/crdsdev/doc/pkg/models"
)

const (
	crdArgCount = 6
	userEnv     = "PG_USER"
	passwordEnv = "PG_PASS"
	hostEnv     = "PG_HOST"
	portEnv     = "PG_PORT"
	dbEnv       = "PG_DB"
)

func main() {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", os.Getenv(userEnv), os.Getenv(passwordEnv), os.Getenv(hostEnv), os.Getenv(portEnv), os.Getenv(dbEnv))
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		fmt.Println(err)
		return
	}
	file := "crd.yaml"
	tagID := 96
	b, err := os.ReadFile(file)
	if err != nil {
		fmt.Print(err)
	}
	crder, err := crd.NewCRDer(b, crd.StripLabels(), crd.StripAnnotations(), crd.StripConversion())
	if err != nil || crder.CRD == nil {
		fmt.Print(err)
		return
	}
	cbytes, err := json.Marshal(crder.CRD)
	if err != nil {
		fmt.Print(err)
		return
	}
	repoCRDs := make(map[string]models.RepoCRD)

	repoCRDs[crd.PrettyGVK(crder.GVK)] = models.RepoCRD{
		Path:     crd.PrettyGVK(crder.GVK),
		Filename: path.Base(file),
		Group:    crder.GVK.Group,
		Version:  crder.GVK.Version,
		Kind:     crder.GVK.Kind,
		CRD:      cbytes,
	}
	allArgs := make([]interface{}, 0, len(repoCRDs)*crdArgCount)
	for _, crd := range repoCRDs {
		allArgs = append(allArgs, crd.Group, crd.Version, crd.Kind, tagID, crd.Filename, crd.CRD)
	}
	// fmt.Println(buildInsert("INSERT INTO crds(\"group\", version, kind, tag_id, filename, data) VALUES ", crdArgCount, len(repoCRDs)) + "ON CONFLICT DO NOTHING")
	// tagstmt := "INSERT INTO tags(name, repo, time) VALUES ($1, $2, $3)"
	teststmt := "INSERT INTO crds(\"group\", version, kind, tag_id, filename, data) VALUES ($1,$2,$3,$4,$5,$6)ON CONFLICT DO NOTHING"
	// _, err = db.Exec(tagstmt, "test", "test", "2022-10-04 07:17:42")
	// if err != nil {
	// 	fmt.Println(err)
	// }
	_, err = db.Exec(teststmt, allArgs...)
	if err != nil {
		fmt.Println(err)
	}
	// db.Exec(buildInsert("INSERT INTO crds(\"group\", version, kind, tag_id, filename, data) VALUES ", crdArgCount, len(repoCRDs))+"ON CONFLICT DO NOTHING", allArgs...)
	// rows, _ := db.Query(`SELECT * FROM "public"."crds" ORDER BY "tag_id","group","version","kind" LIMIT 300 OFFSET 0;`)
	// for rows.Next() {
	// 	fmt.Println(rows.Scan())
	// }
}

func buildInsert(query string, argsPerInsert, numInsert int) string {
	absArg := 1
	for i := 0; i < numInsert; i++ {
		query += "("
		for j := 0; j < argsPerInsert; j++ {
			query += "$" + fmt.Sprint(absArg)
			if j != argsPerInsert-1 {
				query += ","
			}
			absArg++
		}
		query += ")"
		if i != numInsert-1 {
			query += ","
		}
	}
	return query
}
