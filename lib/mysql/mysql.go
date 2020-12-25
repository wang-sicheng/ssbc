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

func InitDB() {
	dsn := fmt.Sprintf("%s:%s@%s(%s:%d)/%s", USERNAME, PASSWORD, NETWORK, SERVER, PORT, DATABASE)
	// 不会校验账号密码是否正确
	tmp, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Infof("Open mysql failed,err:%v\n", err)
		return
	}

	DB = tmp
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

func QueryLastBlock() common.Block {
	block := new(common.Block)
	rows, err := DB.Query("select * from block order by id desc limit 1")
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
		return *block
	}
	for rows.Next() {
		err = rows.Scan(&block.Id, &block.PrevHash, &block.Hash, &block.MerkleRoot, &block.TxCount, &block.Signature, &block.Timestamp)
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return *block
		}
		//fmt.Print(*block)
	}
	return *block
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
func InsertBlock(block common.Block) int {
	result, err := DB.Exec("insert INTO block(prev_hash, hash, tx_count, merkle_root, signature) values(?,?,?,?,?)", block.PrevHash, block.Hash, len(block.TX), block.MerkleRoot, block.Signature)
	if err != nil {
		fmt.Printf("Insert failed,err:%v", err)
		return -1
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		//fmt.Printf("Get lastInsertID failed,err:%v", err)
		return -1
	}
	//fmt.Println("LastInsertID:", lastInsertID)
	_, err = result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return -1
	}
	//fmt.Println("RowsAffected:", rowsaffected)
	return int(lastInsertID)
}

func InsertTransaction(block common.Block) {
	blockId := block.Id
	for _, tx := range block.TX {
		result, err := DB.Exec("insert INTO `transaction`(block_id, sender_address, receiver_address, signature, message, sender_public_key, transfer_amount) values(?,?,?,?,?,?,?)",
			blockId, tx.SenderAddress, tx.ReceiverAddress, tx.Signature, tx.Message, tx.SenderPublicKey, tx.TransferAmount)
		if err != nil {
			fmt.Printf("Insert failed,err:%v", err)
			return
		}
		//lastInsertID, err := result.LastInsertId()
		if err != nil {
			//fmt.Printf("Get lastInsertID failed,err:%v", err)
			return
		}
		//fmt.Println("LastInsertID:", lastInsertID)
		_, err = result.RowsAffected()
		if err != nil {
			fmt.Printf("Get RowsAffected failed,err:%v", err)
			return
		}
		//fmt.Println("RowsAffected:", rowsaffected)
	}
}

func InsertAccount(ac common.Account) {
	result, err := DB.Exec("insert INTO `account`(address, public_key, private_key) values(?,?,?)",
		ac.Address, ac.PublicKey, ac.PrivateKey)
	if err != nil {
		fmt.Printf("Insert failed,err:%v", err)
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		fmt.Printf("Get RowsAffected failed,err:%v", err)
		return
	}
}

func QueryAccountInfo(address string) common.Account{
	ac := new(common.Account)
	rows, err := DB.Query("select * from account where address=?", address)
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		fmt.Printf("Query failed,err:%v", err)
		return *ac
	}
	for rows.Next() {
		err = rows.Scan(&ac.Address, &ac.PrivateKey, &ac.PublicKey)
		if err != nil {
			fmt.Printf("Scan failed,err:%v", err)
			return *ac
		}
	}
	return *ac
}

func CloseDB() error {

	return DB.Close()
}
