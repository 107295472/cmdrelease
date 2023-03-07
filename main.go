package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"
)

func main() {

	var ip string         //服务器ip
	var keyfile string    //用户密钥
	var port string       //服务器端口
	var user string       //服务器用户名
	var localfile string  //要上传的文件名
	var remotepath string //远程路径
	var cmd string        //在服务器上执行的命令
	var uptype string     //上传类型
	var xfile string      //排除的文件名
	var pw string
	app := &cli.App{
		Name:    "cmdrelease",  // cli name
		Version: "v1",          // cli version
		Usage:   "linux程序发布工具", // usage
		Flags: []cli.Flag{ // 接受的 flag
			&cli.StringFlag{ // string
				Name:        "ip",          // flag 名称
				Aliases:     []string{"i"}, // 别名
				Value:       "",            // 默认值
				Usage:       "连接的ip",
				Destination: &ip,  // 指定地址，如果没有可以通过 *cli.Context 的 GetString 获取
				Required:    true, // flag 必须设置
			},
			&cli.StringFlag{ // string
				Name:        "keyfile",     // flag 名称
				Aliases:     []string{"k"}, // 别名
				Value:       "",            // 默认值
				Usage:       "密钥文件路径",
				Destination: &keyfile, //
				Required:    false,    // flag 必须设置
			},
			&cli.StringFlag{ // string
				Name:        "port",        // flag 名称
				Aliases:     []string{"p"}, // 别名
				Value:       "22",          // 默认值
				Usage:       "端口默认22",
				Destination: &port, //
				Required:    false, // flag 必须设置
			},
			&cli.StringFlag{ // string
				Name:        "user",        // flag 名称
				Aliases:     []string{"u"}, // 别名
				Value:       "root",        // 默认值
				Usage:       "用户名默认root",
				Destination: &user, //
				Required:    false, // flag 必须设置
			},
			&cli.StringFlag{ // string
				Name: "pw", // flag 名称
				// Aliases:     []string{"pw"}, // 别名
				Value:       "", // 默认值
				Usage:       "密码",
				Destination: &pw,
				Required:    false,
			},
			&cli.StringFlag{ // string
				Name:        "uploadfile",  // flag 名称
				Aliases:     []string{"f"}, // 别名
				Value:       "",            // 默认值
				Usage:       "上传的文件名",
				Destination: &localfile, //
				Required:    true,       // flag 必须设置
			},
			&cli.StringFlag{ // string
				Name:        "remote",      // flag 名称
				Aliases:     []string{"r"}, // 别名
				Value:       "",            // 默认值
				Usage:       "上传到服务器路径",
				Destination: &remotepath, //
				Required:    true,        // flag 必须设置
			},
			&cli.StringFlag{ // string
				Name:        "cmd",         // flag 名称
				Aliases:     []string{"c"}, // 别名
				Value:       "",            // 默认值
				Usage:       "在服务器上执行的命令",
				Destination: &cmd,  //
				Required:    false, // flag 必须设置
			},
			&cli.StringFlag{ // string
				Name:        "uptype",      // flag 名称
				Aliases:     []string{"t"}, // 别名
				Value:       "1",           // 默认值
				Usage:       "文件=1,目录=2",
				Destination: &uptype, //
				Required:    false,   // flag 必须设置
			},
			&cli.StringFlag{ // string
				Name:        "xfile",       // flag 名称
				Aliases:     []string{"x"}, // 别名
				Value:       "",            // 默认值
				Usage:       "排除的文件名",
				Destination: &xfile, //
				Required:    false,  // flag 必须设置
			},
		},
		Action: func(c *cli.Context) error {
			if ip != "" && localfile != "" && remotepath != "" {

				ExeCmd(ip, port, user, pw, keyfile, localfile, remotepath, cmd, uptype, xfile)
			} else {
				fmt.Println("val:", ip, keyfile, localfile, remotepath)
				fmt.Println("参数错误")
			}
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
func GetKey(file string) string {
	// fmt.Println("file:", file)
	b, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("key不存在")
	} // 将读取的yaml文件解析为响应的 struct
	return string(b)
}
func ExeCmd(ip string, port string, user string, pw string, keyfile string, file string, remote string, cmd string, uptype string, xfile string) {

	var keyFiles = []string{}
	// fmt.Println("key:", keystr)
	p, _ := strconv.Atoi(port)
	// if user == "" && keyfile == "" {
	// 	println("user和keyfile至少有一个值")
	// 	return
	// }
	if pw == "" {
		pw = "root123"
	}
	if keyfile != "" {
		keyFiles = []string{keyfile}
	}
	c, err := New(&Config{User: user, Password: pw, Host: ip, Port: p, KeyFiles: keyFiles})
	if err != nil {
		println(err)
	}
	// defer c.SSHSession.Close()
	defer c.Close()
	// fmt.Println("cmd:", cmd)
	// 执行远程命令
	// var err error
	// if cmd != "" {
	// 	var r, err = c.Output(cmd + " k")
	// 	if err != nil {
	// 		fmt.Println(err.Error())
	// 	}
	// 	fmt.Println(string(r))
	// }
	time.Sleep(time.Duration(1) * time.Second)
	if uptype == "1" {
		b := c.IsExist(remote)
		c.UploadFile(file, remote)
		time.Sleep(time.Duration(100) * time.Millisecond)
		if !b {
			_, _ = c.Output("chmod +x " + remote)
		}
	}
	if uptype == "2" {
		c.UploadDir(file, remote, xfile)
	}
	// c.SSHSession.Close()
	// defer c.SSHSession.Close()
	// ns, _ := c.SSHClient.NewSession()
	if cmd != "" {
		time.Sleep(time.Duration(300) * time.Millisecond)
		// c, _ = ssh.New(&ssh.Config{User: user, Host: ip, Port: p, KeyFiles: keyFiles})
		// cmd = strings.Replace(cmd, "\\", "", 1)
		// println("cmd:" + cmd)
		var rust, err = c.Output(cmd)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(string(rust))
	}
	// findfile := "$(ps -ef | grep haocall_amd64 | grep -v 'grep' | head -1 | awk '{print $2}')"
	// if cmd != "" {
	// 	var findp, _ = c.Output(findfile)
	// 	fmt.Println("no start：", string(findp))
	// 	if string(findp) == "" {
	// 		var r, _ = c.Output("uanme")
	// 		fmt.Println(r)
	// 	}
	// }

	// result, _ := c.Output("uname")
	// fmt.Println(string(result))
}
