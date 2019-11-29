package main

import (
	"log"
	"fmt"
	"reflect"
	"github.com/gomodule/redigo/redis"
)

//List 用于交易接收存储 LPUSH LPULL LRANGE
//		LPUSH 存 LPULL 取 LRANGE获取一定数量交易 用于公共交易集

//set 用于查找公共交易集
//  LRANGE获取一定交易后 存set 收到其他节点交易集后存一个set   最后4个set 取交集
//  没被选中的放到一个list 下次用
//Hash 用于缓存块等

type RedisConn struct {
	conn redis.Conn
}

//连接
func redisConn(netWork, address string) (redis.Conn, error) {
	conn, err := redis.Dial(netWork, address)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return conn, nil
}

//set string
func (conn *RedisConn) SetString(key, value string) error {
	_, err := conn.conn.Do("SET", key, value)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

//get string
func (conn *RedisConn) GetString(key string) (string, error) {
	value, err := redis.String(conn.conn.Do("GET", key))
	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	return value, nil
}

//set out time
func (conn *RedisConn) SetExpire(key string, time int) error {
	_, err := conn.conn.Do("expire", key, time)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

//批量获取mget、批量设置mset
func (conn *RedisConn) mSet(keyValue ...string) error {
	_, err := conn.conn.Do("MSET", keyValue[0], keyValue[1], keyValue[2], keyValue[3])
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (conn *RedisConn) mGet(key ...string) {
	res, err := conn.conn.Do("MGET", key[0], key[1])
	if err != nil {
		log.Println(err.Error())
	}

	fmt.Printf("type: %s\n", reflect.TypeOf(res))
	fmt.Printf("value: %s\n", res)
}

//列表
func (conn *RedisConn) lPush() {
	_, err := conn.conn.Do("LPUSH", "list1", "Java", "Python", "Golang")
	if err != nil {
		log.Println(err.Error())
	}
}

func (conn *RedisConn) lPop() {
	res, err := conn.conn.Do("LPOP", "list1")
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Printf("%s\n", res)
}


//hash
func (conn *RedisConn) hashSet() {
	_, err := conn.conn.Do("HSET", "student", "name", "小雪")
	if err != nil {
		fmt.Println("haha")
		log.Println(err.Error())
	}
}

func (conn *RedisConn) hashGet() {
	res, err := conn.conn.Do("HGET", "student", "name")
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Printf("%s\n", res)
}

//hash more option
func (conn *RedisConn) hashMSet() {
	_, err := conn.conn.Do("HMSET", "students", "name", "小雪", "age", 6, "sex", "女")
	if err != nil {
		log.Println(err.Error())
	}
}

func (conn *RedisConn) hashMGet() {
	//int64s 转换返回的数组值
	res, err := redis.Strings(conn.conn.Do("HMGET", "students", "sex"))
	if err != nil {
		log.Println(err.Error())
	}

	for k, v := range res {
		fmt.Println(k, v)
	}

	fmt.Println(res)
}

func (conn *RedisConn) poolGetConn() {

}


var Pool redis.Pool

//MaxActive 最大连接数，即最多的tcp连接数，一般建议往大的配置，但不要超过操作系统文件句柄个数（centos下可以ulimit -n查看）。
//MaxIdle 最大空闲连接数，即会有这么多个连接提前等待着，但过了超时时间也会关闭。
//IdleTimeout 空闲连接超时时间，但应该设置比redis服务器超时时间短。否则服务端超时了，客户端保持着连接也没用。
//Wait 这是个很有用的配置。好多东抄抄本抄抄的文章都没有提。如果超过最大连接，是报错，还是等待。

func init() {
	Pool = redis.Pool{
		//最大的激活连接数，同时最多有N个连接
		MaxActive: 20,
		//最大的空闲连接数，即使没有redis连接时依然可以保持N个空闲的连接，而不被清除，随时处于待命状态
		MaxIdle: 10,
		//空闲连接等待时间，超过此时间后，空闲连接将被关闭
		IdleTimeout: 120,
		//
		Wait: true,
		//连接方法
		Dial: func() (redis.Conn, error) {
			return redisConn("tcp", "127.0.0.1:6379")
		},
	}
}
func main() {
	conn, err := redisConn("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Println(err.Error())
	}

	defer conn.Close()

	redisConn := &RedisConn{conn: conn}
	err = redisConn.SetString("name", "JayeWu")
	if err != nil {
		log.Println(err.Error())
	}

	err = redisConn.SetExpire("name", 10)
	if err != nil {
		log.Println(err.Error())
	}

	name, err := redisConn.GetString("name")
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println(name)

	redisConn.mSet("name", "Jaye", "sex", "男")
	redisConn.mGet("name", "sex")

	fmt.Println("########list########")
	redisConn.lPush()
	res, err := redisConn.conn.Do("LRANGE", "list1", 1, 2)
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Printf("%s\n", res)
	redisConn.lPop()

	fmt.Println("$$$$$$$$$$hash$$$$$$$$$$")
	redisConn.hashSet()
	redisConn.hashGet()

	fmt.Println("$$$$$$$$$$hash-M$$$$$$$$$$")
	redisConn.hashMSet()
	redisConn.hashMGet()

	fmt.Println("###########Pool############")
	pconn := Pool.Get()
	_, err = pconn.Do("SET", "home", "102004")
	if err != nil {
		log.Println(err.Error())
	}
	res, err = redis.String(pconn.Do("GET", "home"))
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Println("pool get value", res)
}