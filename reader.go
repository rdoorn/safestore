package main

import "os"
import "sync/atomic"
import "log"

type CustomReader struct {
    fp   *os.File
    size int64
    read int64
}

func (r *CustomReader) Read(p []byte) (int, error) {
    return r.fp.Read(p)
}

func (r *CustomReader) ReadAt(p []byte, off int64) (int, error) {
    n, err := r.fp.ReadAt(p, off)
    if err != nil {
        return n, err
    }

    // Got the length have read( or means has uploaded), and you can construct your message
    atomic.AddInt64(&r.read, int64(n))

    // I have no idea why the read length need to be div 2,
    // maybe the request read once when Sign and actually send call ReadAt again
    // It works for me
    log.Printf("total read:%d    progress:%d%%\n", r.read/2, int(float32(r.read*100/2)/float32(r.size)))

    return n, err
}

func (r *CustomReader) Seek(offset int64, whence int) (int64, error) {
    return r.fp.Seek(offset, whence)
}
