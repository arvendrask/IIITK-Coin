package main
import(

	"database/sql"
	_"github.com/mattn/go-sqlite3"
	"fmt"
)
type data struct {
	roll_no int
	name string
}
func Ins(d1 data, db *sql.DB){
	statement,_:= db.Prepare("INSERT INTO roll (roll_no, name) VALUES(?,?)")
	statement.Exec(d1.roll_no, d1.name)
}
func main() {
	fmt.Println("Hello")
	database,_:= sql.Open("sqlite3","./mydb.db")
	statement,_:= database.Prepare("CREATE TABLE IF NOT EXISTS roll (roll_no INTEGER, name TEXT)")
	statement.Exec()
	Ins(data{190183, "arvendra"}, database)

	//fmt.Printf("%T",database)
}
