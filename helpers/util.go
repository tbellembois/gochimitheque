package helpers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/csv"
	"io"
	"log"
)

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
