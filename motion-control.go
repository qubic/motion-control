package main

import (
    "github.com/go-martini/martini"
    "github.com/martini-contrib/render"
    "os"
    "os/exec"
    "log"
    "net/http"
)

func main() {
    motionPidFile := "/var/run/motion/motion.pid"
    logFile := "/tmp/motion-control.log"
    f, err := os.OpenFile(logFile, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()
    logger := log.New(f, "[martini]", log.LstdFlags)
    m := martini.Classic()
    m.Map(logger)
    m.Use(render.Renderer())
    m.Get("/", func(r render.Render) {
        if _, err := os.Stat(motionPidFile); err == nil {
            r.HTML(200, "status", map[string]interface{}{"motionStatus": "running", "availableAction": "stop"})
        } else {
            r.HTML(200, "status", map[string]interface{}{"motionStatus": "not running", "availableAction": "start"})
        } 
    })
    m.Get("/stop", func(res http.ResponseWriter, req *http.Request) {
        cmd := exec.Command("service", "motion", "stop")
        err := cmd.Run()
        if err != nil {
            logger.Println(err)
        } else {
            logger.Println("Motion has been stopped")
        }
        http.Redirect(res, req, "/", http.StatusSeeOther)
    })
    m.Get("/start", func(res http.ResponseWriter, req *http.Request) {
        cmd := exec.Command("service", "motion", "start")
        err := cmd.Run()
        if err != nil {
            logger.Println(err)
        } else {
            logger.Println("Motion has been started")
        }
        http.Redirect(res, req, "/", http.StatusSeeOther)
    })
    m.Run()
}

