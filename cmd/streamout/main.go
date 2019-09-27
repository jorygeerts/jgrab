package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"
)

type msg struct {
	Name string    `json:"name"`
	TS   time.Time `json:"timestamp"`
	List []string  `json:"list"`
	Sub  sub       `json:"sub"`
}

type sub struct {
	Name string `json:"name"`
}

func main() {

	log.Println("This is a helper to test jgrab. Run this and pipe the output (stdout, not stderr) into jgrab to see it in action.")

	enc := json.NewEncoder(os.Stdout)
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	str := func() string {
		return String(seededRand.Intn(8)+seededRand.Intn(8), seededRand)
	}

	for i := 0; i < 10; i++ {
		_ = enc.Encode(msg{
			Name: str(),
			TS:   time.Now(),
			List: []string{str(), str()},
			Sub:  sub{Name: str()},
		})
		time.Sleep(time.Second)
	}
}

func StringWithCharset(length int, charset string, seededRand *rand.Rand) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int, seededRand *rand.Rand) string {
	return StringWithCharset(length, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", seededRand)
}
