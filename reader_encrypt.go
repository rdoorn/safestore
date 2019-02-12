package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"hash"
	"io"
	"log"
	"os"
	"sync/atomic"
)

type ShaStreamReader struct {
	fp   *os.File
	size int64
	read int64
	iv   []byte
	ctr  cipher.Stream
	hmac hash.Hash
}

func (r *ShaStreamReader) Read(p []byte) (int, error) {
	return r.fp.Read(p)
}

func (r *ShaStreamReader) ReadAt(p []byte, off int64) (int, error) {

	buf := make([]byte, len(p))
	log.Printf("buf: %s (%d)", buf, len(buf))

	n, err := r.fp.ReadAt(buf, off)
	if err != nil {
		return n, err
	}

	r.ctr.XORKeyStream(p, buf[:n])
	r.hmac.Write(p)
	log.Printf("p: %s", p)
	/*r.fp.Rea

	ctr.XORKeyStream(outBuf, buf[:n])
	hmac.Write(outBuf)
	*/
	// Got the length have read( or means has uploaded), and you can construct your message
	atomic.AddInt64(&r.read, int64(n))

	// I have no idea why the read length need to be div 2,
	// maybe the request read once when Sign and actually send call ReadAt again
	// It works for me
	log.Printf("total read:%d    progress:%d%%\n", r.read/2, int(float32(r.read*100/2)/float32(r.size)))

	return n, err
}

func (r *ShaStreamReader) Seek(offset int64, whence int) (int64, error) {
	return r.fp.Seek(offset, whence)
}

const BUFFER_SIZE int = 4096
const IV_SIZE int = 16

func encrypt(filePathIn, filePathOut string, keyAes, keyHmac []byte) error {
	inFile, err := os.Open(filePathIn)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(filePathOut)
	if err != nil {
		return err
	}
	defer outFile.Close()

	iv := make([]byte, IV_SIZE)
	_, err = rand.Read(iv)
	if err != nil {
		return err
	}

	aes, err := aes.NewCipher(keyAes)
	if err != nil {
		return err
	}

	ctr := cipher.NewCTR(aes, iv)
	hmac := hmac.New(sha256.New, keyHmac)

	buf := make([]byte, BUFFER_SIZE)
	for {
		n, err := inFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		outBuf := make([]byte, n)
		ctr.XORKeyStream(outBuf, buf[:n])
		hmac.Write(outBuf)
		outFile.Write(outBuf)

		if err == io.EOF {
			break
		}
	}

	outFile.Write(iv)
	hmac.Write(iv)
	outFile.Write(hmac.Sum(nil))

	return nil
}
