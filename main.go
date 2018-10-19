package main

import (
	"fmt"
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	"github.com/go-ini/ini"
	"github.com/go-vgo/robotgo"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var mainwin *ui.Window
var configPath string

func myUi() ui.Control {
	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)

	hbox := ui.NewHorizontalBox()
	hbox.SetPadded(true)

	input1 := ui.NewEntry()
	input2 := ui.NewEntry()

	x, y := getPosition(configPath)
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

func init() {
	args := os.Args
	if strings.Contains(args[0], "go-build") {
		configPath = "./config/setting.ini"
	} else {
		configPath, _ = mainDir()
		configPath += "/config/setting.ini"
	}
}

func main() {

	robotgo.ScrollMouse(10, "up")
	robotgo.MouseClick("left", true)
	robotgo.MoveMouseSmooth(100, 200, 1.0, 100.0)

	ui.Main(setupUI)
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

	if savePosition(xStr, yStr, configPath) {
		ui.MsgBox(mainwin, "消息", "保存成功")
	} else {
		ui.MsgBoxError(mainwin, "错误", "保存失败")
	}
}

func getPosition(configPath string) (x, y string) {
	cfg, err := ini.Load(configPath)
	if err != nil {
		fmt.Printf("Fail to read file: %v %s", err, configPath)
		os.Exit(1)
	}

	return cfg.Section("position").Key("x").String(), cfg.Section("position").Key("y").String()
}

func savePosition(x, y, configPath string) bool {
	cfg, err := ini.Load(configPath)
	if err != nil {
		fmt.Printf("Fail to read file: %v %s", err, configPath)
		os.Exit(1)
	}
	cfg.Section("position").Key("x").SetValue(x)
	cfg.Section("position").Key("y").SetValue(y)
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
