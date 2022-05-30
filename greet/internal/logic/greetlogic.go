package logic

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"strings"

	"greet/internal/svc"
	"greet/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GreetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

type User struct {
	Id       int64
	Username string
}

func NewGreetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GreetLogic {
	return &GreetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GreetLogic) Greet(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line
	response := types.Response{Message: fmt.Sprintf("%v - %v", "chenbo", os.Getenv("MYSQL_HOST"))}
	_, err = conDB()
	if err != nil {
		response = types.Response{
			Message: err.Error(),
		}
	} else {
		rows, err := selectData(1510)
		if err != nil {
			response = types.Response{
				Message: err.Error(),
			}
		}
		var temp []string
		for _, v := range rows {
			temp = append(temp, v.Username)
		}
		data := strings.Join(temp, "\n")
		response = types.Response{Message: fmt.Sprintf("查询编号%v的用户姓名%v", req.Name, data)}
	}
	return &response, nil
}

func selectData(id int64) ([]User, error) {
	var users []User
	db, err := conDB()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query("select id,username from user where id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("获取用户信息出现异常错误 %d: %v", id, err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(rows)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.Username); err != nil {
			return nil, fmt.Errorf("获取用户信息出现异常错误 %d: %v", id, err)
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("获取用户信息出现异常错误 %d: %v", id, err)
	}
	return users, nil
}

func conDB() (*sql.DB, error) {
	cfg := mysql.Config{
		User:   os.Getenv("MYSQL_USERNAME"),
		Passwd: os.Getenv("MYSQL_PASSWORD"),
		Net:    "tcp",
		Addr:   fmt.Sprintf("%v:%v", os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_PORT")),
		DBName: os.Getenv("MYSQL_DATABASE"),
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	pingErr := db.Ping()
	if pingErr != nil {
		return nil, fmt.Errorf("无法正常连接数据库")
	}
	return db, nil
}
