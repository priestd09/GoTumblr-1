package main

import(
    "log"
)

func Info(format string, msgs ...interface{}){
    log.Printf(format, msgs...)
}

func Error(format string, msgs ...interface{}){
    log.Printf(format, msgs...)
}

func Warn(format string, msgs ...interface{}){
    log.Printf(format, msgs...)
}
