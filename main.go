package main

import (
	"flag"
	"fmt"
	"html/template"
	"net/http"
	_ "net/http/pprof"

	"github.com/media_uploader/core"
	handlers "github.com/media_uploader/handlers"
)

var (
	addr                      = flag.String("addr", "localhost:8080", "HTTP service address")
	perf                      = flag.Bool("perf", false, "Enable performance testing")
	workers                   = flag.Int("workers", 10, "Number of workers")
	chBufferSize              = flag.Int("chBufferSize", 100, "Channel buffer size")
	workerMemoryLimit         = flag.Uint64("workerMemoryLimit", 75, "Worker memory limit")
	useSpawnerWithMemoryLimit = flag.Bool("useSpawnerWithMemoryLimit", true, "Use worker spawner with memory limit")
	enableSimpleInterface     = flag.Bool("enableSimpleInterface", false, "Enable simple interface to upload files")
	saveUploadsTemporarily    = flag.Bool("saveUploadsTemporarily", false, "Save uploaded files temporarily")

	streamTemplate     *template.Template
	fileSelectTemplate *template.Template
)

func main() {
	flag.Parse()

	core.InitializeLogger()

	fmt.Printf("Starting `media_uploader` at http://%s\n", *addr)

	handlers.InitializeWorkerConfig(*workers, *chBufferSize, *workerMemoryLimit)
	handlers.SaveUploadsTemporarily = *saveUploadsTemporarily

	var err error

	fmt.Println("enableSimpleInterface: ", *enableSimpleInterface)
	if *perf {
		go func() {
			err := http.ListenAndServe("localhost:6060", nil)
			if err != nil {
				core.LogError("Failed to start pprof server", err)
			}
		}()
	}

	if *enableSimpleInterface {
		err = parseHTMLTemplates()
		if err != nil {
			core.LogError("Failed to parse HTML templates", err)
			return
		}
	}

	http.HandleFunc("/upload_stream", handlers.StreamHandler)

	if *enableSimpleInterface {
		http.HandleFunc("/stream", stream)
		http.HandleFunc("/file_select", fileSelect)
	}

	go func() {
		err = http.ListenAndServe(*addr, nil)
		if err != nil {
			core.LogFatal("Failed to start server", err)
		}
	}()

	if !*useSpawnerWithMemoryLimit {
		handlers.WorkerPool.Start()
	} else {
		handlers.WorkerSpawner.Start()
	}

	if !*useSpawnerWithMemoryLimit {
		handlers.WorkerPool.Close()
	} else {
		handlers.WorkerSpawner.Close()
	}
}

func parseHTMLTemplates() error {
	var err error

	streamTemplate, err = template.ParseFiles("./static/stream.html")
	if err != nil {
		return err
	}

	fileSelectTemplate, err = template.ParseFiles("./static/file_select.html")
	if err != nil {
		return err
	}

	return nil
}

func stream(w http.ResponseWriter, r *http.Request) {
	streamTemplate.Execute(w, "ws://"+r.Host+"/upload_stream")
}

func fileSelect(w http.ResponseWriter, r *http.Request) {
	fileSelectTemplate.Execute(w, "ws://"+r.Host+"/upload_stream")
}
