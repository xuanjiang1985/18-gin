package main

import (
	"bytes"
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/go-ini/ini"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var mainwin *ui.Window
var configPath string
var savePath string
var msgEntry *ui.MultilineEntry

func myUi() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	hbox1 := ui.NewHorizontalBox()
	hbox1.SetPadded(true)

	input1 := ui.NewEntry()
	input2 := ui.NewEntry()

	x, y := getPosition(configPath)
	input1.SetText(x)
	input2.SetText(y)

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	entryForm.Append("网址", input1, true)
	entryForm.Append("并发", input2, false)

	hbox.Append(entryForm, true)

	storeBtn := ui.NewButton("保存")
	startBtn := ui.NewButton("启动")
	stopBtn := ui.NewButton("停止")
	clearBtn := ui.NewButton("清除")
	hbox1.Append(storeBtn, false)
	hbox1.Append(startBtn, false)
	hbox1.Append(stopBtn, false)
	hbox1.Append(clearBtn, false)

	storeBtn.OnClicked(func(*ui.Button) {
		store(input1, input2)
	})

	clearBtn.OnClicked(func(*ui.Button) {
		msgEntry.SetText("")
	})

	startBtn.OnClicked(func(this *ui.Button) {
		sendMsg(msgEntry, "脚本启动")
		getUrlData()
		this.Disable()
	})

	stopBtn.OnClicked(func(*ui.Button) {
		sendMsg(msgEntry, "脚本停止")
		startBtn.Enable()
		storeBtn.Enable()
	})

	vbox.Append(hbox, false)
	vbox.Append(hbox1, false)
	vbox.Append(ui.NewHorizontalSeparator(), false)

	// 下部分
	hbox2 := ui.NewHorizontalBox()
	msgEntry = ui.NewMultilineEntry()
	msgEntry.SetReadOnly(true)
	hbox2.Append(msgEntry, true)

	vbox.Append(hbox2, true)
	return vbox
}

func setupUI() {
	mainwin = ui.NewWindow("网站压测", 350, 480, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	//tab := ui.NewTab()
	mainwin.SetChild(myUi())
	mainwin.SetMargined(true)

	mainwin.Show()
}

func init() {
	args := os.Args
	if strings.Contains(args[0], "go-build") {
		configPath = "./config/setting.ini"
		savePath = "./saveData"
	} else {
		configPath, _ = mainDir()
		configPath += "/config/setting.ini"
		savePath += "/saveData"
	}
}

func main() {
	ui.Main(setupUI)
}

func store(input1, input2 *ui.Entry) bool {
	xStr := input1.Text()
	yStr := input2.Text()

	if xStr == "" || yStr == "" {
		sendMsg(msgEntry, "不能为空")
		return false
	}

	qty, err := strconv.Atoi(yStr)
	if err != nil {
		sendMsg(msgEntry, "必须是整数")
		return false
	}

	if qty <= 0 || qty > 100 {
		sendMsg(msgEntry, "并发数量在1~100之间 ")
		return false
	}

	if savePosition(xStr, yStr, configPath) {
		sendMsg(msgEntry, "保存成功")
		return true
	} else {
		sendMsg(msgEntry, "保存失败")
		return false
	}
}

func getUrlData() {
	url, qty := getPosition(configPath)

	qtyInt, err := strconv.Atoi(qty)
	if err != nil {
		sendMsg(msgEntry, "并发数量必须是整数")
	}

	if qtyInt <= 0 || qtyInt > 100 {
		sendMsg(msgEntry, "并发数量在1~100之间 ")
		return
	}

	// 并发请求http get
	for i := 1; i <= qtyInt; i++ {
		go getEachUrl(i, url)
	}

}

// 下载网页
func getEachUrl(i int, url string) {
	defer func() {
		if err := recover(); err != nil {
			sendMsg(msgEntry, "错误退出"+strconv.Itoa(i)+"次")
		}
	}()
	str := "正在发起第" + strconv.Itoa(i) + "次请求"
	ui.QueueMain(func() {
		sendMsg(msgEntry, str)
	})
	// http.get
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		str = "第" + strconv.Itoa(i) + "次请求报错"
	}
	defer resp.Body.Close()

	out, err := os.Create(savePath + "/" + strconv.Itoa(i) + ".html")
	content, err := ioutil.ReadAll(resp.Body)
	_, err = io.Copy(out, bytes.NewReader(content))
	if err != nil {
		str = "第" + strconv.Itoa(i) + "次保存报错"
	}
	val := time.Since(start)
	str = "第" + strconv.Itoa(i) + "次请求耗时" + val.String()
	ui.QueueMain(func() {
		sendMsg(msgEntry, str)
	})
}

func getPosition(configPath string) (x, y string) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		fmt.Printf("Fail to read file: %v %s", err, configPath)
		os.Exit(1)
	}

	return cfg.Section("server").Key("host").String(), cfg.Section("server").Key("qty").String()
}

func savePosition(x, y, configPath string) bool {
	cfg, err := ini.Load(configPath)
	if err != nil {
		fmt.Printf("Fail to read file: %v %s", err, configPath)
		os.Exit(1)
	}
	cfg.Section("server").Key("host").SetValue(x)
	cfg.Section("server").Key("qty").SetValue(y)
	err = cfg.SaveTo(configPath)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func mainDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	return dir, nil
}

func sendMsg(box *ui.MultilineEntry, msg string) {
	timeStr := time.Now().Format("01-02 15:04:05")
	//hastext := box.Text()
	// if utf8.RuneCountInString(hastext) > 1000 {
	// 	box.SetText(timeStr + " " + msg + "\n")
	// } else {
	box.Append(timeStr + " " + msg + "\n")
	//}
}

// func webService() {
// 	router := gin.Default()
// 	// Query string parameters are parsed using the existing underlying request object.
// 	// The request responds to a url matching:  /welcome?firstname=Jane&lastname=Doe
// 	router.GET("/info", func(c *gin.Context) {

// 		msg := c.Query("msg") // shortcut for c.Request.URL.Query().Get("lastname")
// 		ui.QueueMain(func() {
// 			sendMsg(msgEntry, msg)
// 		})
// 	})
// 	router.Run(":8081")
// }
