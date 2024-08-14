package main

import (
	"encoding/json"
	"fmt"
	"github.com/glebarez/sqlite"
	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var db *gorm.DB

var (
	port int
)

var (
	Version   = "unknown version"
	Commit    = "unknown commit"
	BuildTime = "unknown time"
)

type User struct {
	Username      string
	Password      string
	ClientVersion string
	ServerVersion string
	RemoteAddr    string
	LocalAddr     string

	CreatedAt time.Time `json:"createdAt" swaggerignore:"true"`
}

var cmd = &cobra.Command{
	Use:     "hole",
	Version: Version,

	Run: func(cmd *cobra.Command, args []string) {
		server()
	},
}

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "read records from database",

	Run: func(cmd *cobra.Command, args []string) {
		read()
	},
}

func init() {
	cmd.AddCommand(readCmd)
	cmd.SetVersionTemplate(fmt.Sprintf("Hole %s[%s] %s %s with %s %s\n", Version, Commit, runtime.GOOS, runtime.GOARCH, runtime.Version(), BuildTime))
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logrus.Fatal(err)
	}
	cmd.PersistentFlags().IntVarP(&port, "port", "p", 2222, "ssh server port")
	driver := sqlite.Open(fmt.Sprintf("file:%s/hole.db?cache=shared", dir))

	client, err := gorm.Open(driver, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		logrus.Fatal(err)
	}

	db = client

	if err := db.AutoMigrate(&User{}); err != nil {
		logrus.Fatal(err)
	}

	if dab, err := client.DB(); err != nil {
		logrus.Fatal(err)
	} else {
		dab.SetMaxIdleConns(10)
		dab.SetMaxOpenConns(40)
		dab.SetConnMaxLifetime(0)
		dab.SetConnMaxIdleTime(0)
	}

}

func (u *User) String() string {
	res, _ := json.Marshal(u)
	return string(res)
}

func read() {
	var res []User

	if err := db.Find(&res).Error; err != nil {
		logrus.Fatal(err)
	}
	for _, v := range res {
		fmt.Println(v.String())
	}
}

func server() {
	server := &ssh.Server{
		Version: "OpenSSH_8.0",
		PasswordHandler: func(ctx ssh.Context, password string) bool {
			user := &User{
				Username:      ctx.User(),
				Password:      password,
				ClientVersion: ctx.ClientVersion(),
				ServerVersion: ctx.ServerVersion(),
				RemoteAddr:    ctx.RemoteAddr().String(),
				LocalAddr:     ctx.LocalAddr().String(),
			}
			if err := db.Create(user).Error; err != nil {
				logrus.Errorf("%s insert failed: %s", user.String(), err)
			}
			return false
		},
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		logrus.Fatal(err)
	}

	if err := server.Serve(listener); err != nil {
		logrus.Fatal(err)
	}
	defer server.Close()
}

func main() {
	if err := cmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
