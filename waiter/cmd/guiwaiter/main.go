package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/tymbaca/coffeeshop/waiter/model"
)

const (
	_winWidth  = 800
	_winHeight = 600
)

var _log = widget.NewMultiLineEntry()

func main() {
	a := app.New()
	w := a.NewWindow("Hello World")
	w.Resize(fyne.NewSize(_winWidth, _winHeight))

	_log.Resize(fyne.NewSize(_winWidth, 400))

	lastOrderedData := binding.NewString()
	lastOrderedLabel := widget.NewLabelWithData(lastOrderedData)

	sendButton := widget.NewButton("order coffee", func() {
		order := randomOrder()

		go sendOrder(order)

		lastOrderedData.Set(string(order.Type))
	})

	content := container.NewVBox(lastOrderedLabel, sendButton, _log)

	w.SetContent(content)
	w.ShowAndRun()
}

func randomOrder() model.Order {
	return model.Order{
		Type: model.Cappuccino,
	}
}

func sendOrder(order model.Order) {
	data, err := json.Marshal(order)
	if err != nil {
		logErr(err)
		return
	}

	resp, err := http.Post("http://localhost:8080/order", "application/json", bytes.NewReader(data))
	if err != nil {
		logErr(err)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logErr(err)
		return
	}

	if resp.StatusCode != 200 {
		logErr(fmt.Errorf("code %d: %s", resp.StatusCode, respBody))
		return
	}
}

func log(text string) {
	_log.Append(text + "\n")
}

func logErr(err error) {
	log("ERR: " + err.Error())
}
