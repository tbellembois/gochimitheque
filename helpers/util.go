package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/csv"
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"
)

// TimeTrack displays the run time of the function "name"
// from the start time "start"
// use: defer helpers.TimeTrack(time.Now(), "GetProducts")
// at the begining of the function to track
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}

// GetPasswordHash return password hash for the login,
func GetPasswordHash(login string) ([]byte, error) {

	h := hmac.New(sha256.New, []byte("secret"))
	h.Write([]byte(login))

	return h.Sum(nil), nil
}

// CSVToMap takes a reader and returns an array of dictionaries, using the header row as the keys
// credit: https://gist.github.com/drernie/5684f9def5bee832ebc50cabb46c377a
func CSVToMap(reader io.Reader) []map[string]string {
	r := csv.NewReader(reader)
	rows := []map[string]string{}
	var header []string
	for {
		record, err := r.Read()
		log.Debug(fmt.Sprintf("record: %s", record))
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]string{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows
}
