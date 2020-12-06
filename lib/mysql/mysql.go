package mysql

import (
	"database/sql"
	"fmt"

	"github.com/cloudflare/cfssl/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/ssbc/common"
)

type User struct {
	ID   int64          `db:"id"`
	Name sql.NullString `db:"name"` //由于在mysql的users表中name没有设置为NOT NULL,所以name可能为null,在查询过程中会返回nil，如果是string类型则无法接收nil,但sql.NullString则可以接收nil值
	Age  int            `db:"age"`
}

const (
	USERNAME = "root"
	PASSWORD = "123456"
	NETWORK  = "tcp"
	SERVER   = "127.0.0.1"
	PORT     = 3306
	DATABASE = "ssbc"
)

var DB *sql.DB

func InitDB() (DB *sql.DB) {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	// 不会校验账号密码是否正确
	tmp, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Infof("Open mysql failed,err:%v\n", err)
		return
	}

	DB = tmp
	return DB
}

func QueryAllBlocks(DB *sql.DB) {
	block := new(common.Block)
	rows, err := DB.Query("select * from block")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
		return
	}
	for rows.Next() {
		err = rows.Scan(&block.Id, &block.Pid, &block.PrevHash, &block.Hash, &block.MerkleRoot, &block.TxCount, &block.Signature, &block.Timestamp)
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return
		}
		fmt.Print(*block)
	}
}

func queryOne(DB *sql.DB) {
	user := new(User)
	row := DB.QueryRow("select * from block where id=?", 1)
	if err := row.Scan(&user.ID, &user.Name, &user.Age); err != nil {
		fmt.Printf("scan failed, err:%v", err)
		return
	}
	fmt.Println(*user)
}

//查询多行
func queryMulti(DB *sql.DB) {
	user := new(User)
	rows, err := DB.Query("select * from block where id > ?", 1)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
		return
	}
	for rows.Next() {
		err = rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return
		}
		fmt.Print(*user)
	}

}
func InsertBlock(block common.Block) {
	insertData(DB, block)
}

//插入数据
func insertData(DB *sql.DB, block common.Block) {
	result, err := DB.Exec("insert INTO block(bIndex,Timestamp,BPM,Hash,PreHash) values(?,?,?,?,?)", block.Id, block.Timestamp, block.Signature, block.Hash, block.PrevHash)
	if err != nil {
		fmt.Printf("Insert failed,err:%v", err)
		return
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		fmt.Printf("Get lastInsertID failed,err:%v", err)
		return
	}
	fmt.Println("LastInsertID:", lastInsertID)
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return
	}
	fmt.Println("RowsAffected:", rowsaffected)
}

//更新数据
func updateData(DB *sql.DB) {
	result, err := DB.Exec("UPDATE users set age=? where id=?", "30", 3)
	if err != nil {
		fmt.Printf("Insert failed,err:%v", err)
		return
	}
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return
	}
	fmt.Println("RowsAffected:", rowsaffected)
}

//删除数据
func deleteData(DB *sql.DB) {
	result, err := DB.Exec("delete from users where id=?", 1)
	if err != nil {
		fmt.Printf("Insert failed,err:%v", err)
		return
	}
	rowsaffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return
	}
	fmt.Println("RowsAffected:", rowsaffected)
}

func CloseDB() error {

	return DB.Close()
}
