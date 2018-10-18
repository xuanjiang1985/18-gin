package main

import (
	"18-gin/api/user"
	"18-gin/utils"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
	"log"
	"os"
	"strconv"
)

var mainwin *ui.Window

func myUi() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	input1 := ui.NewEntry()
	input2 := ui.NewEntry()
	x, y := getPosition()
	input1.SetText(x)
	input2.SetText(y)

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	entryForm.Append("横坐标", input1, false)
	entryForm.Append("纵坐标", input2, false)

	hbox.Append(entryForm, false)

	storeBtn := ui.NewButton("保存")
	hbox.Append(storeBtn, false)

	storeBtn.OnClicked(func(*ui.Button) {
		store(input1, input2)
	})

	vbox.Append(hbox, false)
	return vbox
}

func setupUI() {
	mainwin = ui.NewWindow("UI配置", 640, 480, true)
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

func main() {
	go webService()

	ui.Main(setupUI)
}

func webService() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("webService Error: ", err)
		}
	}()

	router := gin.Default()
	router.Use(utils.Cors())
	// This handler will match /user/john but will not match /user/ or /user
	router.GET("/user/:name", user.Name)
	router.Run(":8091")
}

func store(input1, input2 *ui.Entry) {
	xStr := input1.Text()
	yStr := input2.Text()

	if xStr == "" || yStr == "" {
		ui.MsgBoxError(mainwin, "警告", "坐标不能为空")
		return
	}

	_, err := strconv.Atoi(xStr)
	if err != nil {
		ui.MsgBoxError(mainwin, "警告", "坐标必须是整数")
		return
	}

	_, err = strconv.Atoi(yStr)
	if err != nil {
		ui.MsgBoxError(mainwin, "警告", "坐标必须是整数")
		return
	}

	if savePosition(xStr, yStr) {
		ui.MsgBox(mainwin, "消息", "保存成功")
	} else {
		ui.MsgBoxError(mainwin, "错误", "保存失败")
	}
}

func getPosition() (x, y string) {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	return cfg.Section("position").Key("x").String(), cfg.Section("position").Key("y").String()
}

func savePosition(x, y string) bool {
	cfg, err := ini.Load("./config/config.ini")
	if err != nil {
		log.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	cfg.Section("position").Key("x").SetValue(x)
	cfg.Section("position").Key("y").SetValue(y)
	err = cfg.SaveTo("./config/config.ini")
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}
