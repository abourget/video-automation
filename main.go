package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"gopkg.in/yaml.v1"

	"github.com/sclevine/agouti"
)

var data map[string]interface{}
var datas []map[string]interface{}

var videoSize = size{1280, 720}

func main() {

	cnt, err := ioutil.ReadFile("data.yaml")
	failOnError(err)
	failOnError(yaml.Unmarshal(cnt, &datas))

	data = datas[0]

	go serveTemplates()

	launchVideoRecording()
}

func serveTemplates() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, err := template.ParseFiles("template.html")
		failOnError(err)

		w.Header().Set("Content-Type", "text/html")

		err = t.Execute(w, data)
		if err != nil {
			fmt.Fprintln(w, "\n\nERROR PARSING TEMPLATE:", err, "\n\n")
		}
	})
	fmt.Println("Listening on 127.0.0.1:7777")
	failOnError(http.ListenAndServe(":7777", nil))
}

func launchVideoRecording() {
	dr := agouti.ChromeDriver()
	failOnError(dr.Start())

	page, err := dr.NewPage()
	if err != nil {
		log.Fatalln("couldn't get Page")
	}

	failOnError(page.Size(videoSize.W+chrome.Left+chrome.Right, videoSize.H+chrome.Top+chrome.Bottom))

	for _, data = range datas {

		failOnError(page.Navigate("http://localhost:7777/"))

		var screenOffset point
		failOnError(page.RunScript(`
return {x: window.screenX, y: window.screenY};
`, nil, &screenOffset))

		// Wait for .capture to appear
		for {
			visible, err := page.Find(".capture").Visible()
			failOnError(err)
			if visible {
				break
			}
		}

		doneCh := launchFFMPEG(screenOffset.X+chrome.Left, screenOffset.Y+chrome.Top, videoSize.W, videoSize.H, fmt.Sprintf("/tmp/video-automation-%v.mp4", data["slug"]))

		// Wait for .capture to disappear
		for {
			visible, _ := page.Find(".capture").Visible()
			if !visible {
				break
			}
		}

		doneCh <- true

		time.Sleep(1 * time.Second)

	}

	dr.Stop()
}

type rect struct {
	X, Y, W, H int
}
type point struct {
	X, Y int
}
type size struct {
	W, H int
}
type dim struct {
	Top, Bottom, Left, Right int
}

var chrome = dim{77, 5, 5, 5}

func launchFFMPEG(x, y, w, h int, filename string) (done chan bool) {
	done = make(chan bool)

	args := []string{"-video_size", fmt.Sprintf("%dx%d", w, h), "-framerate", "25", "-f", "x11grab", "-draw_mouse", "0", "-i", fmt.Sprintf(":0.0+%d,%d", x, y), "-vcodec", "libx264", "-preset", "veryfast", "-y", filename}
	fmt.Printf("\n\nCommand: ffmpeg %s\n\n", args)
	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdinPipe, err := cmd.StdinPipe()
	failOnError(err)

	go func() {
		failOnError(cmd.Start())
		<-done
		stdinPipe.Write([]byte("q"))
	}()

	return done
}

func failOnError(err error) {
	if err != nil {
		log.Fatalln("Failed", err)
	}
}
