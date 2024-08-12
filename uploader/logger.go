package uploader

import (
    "log"
    "os"
)

func LogData(info interface{}) {
	fileInfo, err := openLogFile("./indexer.log")
    if err != nil {
        log.Fatal(err)
    }
    infoLog := log.New(fileInfo, "[indexer]", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
    infoLog.Println(info)
}

func openLogFile(path string) (*os.File, error) {
    logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
    if err != nil {
        return nil, err
    }
    return logFile, nil
}