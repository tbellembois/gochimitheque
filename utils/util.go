package utils

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/smtp"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/globals"
	. "github.com/tbellembois/gochimitheque/models"
)

// SortEmpiricalFormula returns the sorted f empirical formula.
func SortEmpiricalFormula(f string) (string, error) {
	var (
		err      error
		newf, sp string
	)

	// zero empirical formula
	if f == "XXXX" {
		return f, nil
	}

	// removing spaces
	f = strings.Replace(f, " ", "", -1)

	// if the formula is like abc.def.ghi, spliting it
	splitf := strings.Split(f, ".")
	if len(splitf) == 1 {
		return SortSimpleFormula(f)
	}

	for _, p := range splitf {
		if sp, err = SortSimpleFormula(p); err != nil {
			return "", err
		}
		newf += "." + sp
	}

	return newf, nil
}

// SortSimpleFormula returns the sorted f formula.
func SortSimpleFormula(f string) (string, error) {
	var (
		// 	hasCatom, hasHatom, hasOtherAtom, hasUpperLowerAtom bool
		hasCatom, hasHatom, hasOatom, hasULatom bool
		upperLowerAtoms, otherAtoms             []string
		// 	lastPart                                            string
	)

	// removing spaces
	f = strings.Replace(f, " ", "", -1)

	// checking formula characters
	if !globals.FormulaRe.MatchString(f) {
		return "", errors.New("invalid characters in formula")
	}

	// search atoms with and uppercase followed by lowercase letters like Na or Cl
	// return a list of tuples like:
	// [[Cl Cl Cl] [Na Na Na] [Cl3 Cl3 Cl]]
	// for ClNaHCl3
	// the third member of the tupple is used to detect duplicated atoms
	ULAtomsRe := regexp.MustCompile("((?:^[0-9]+)?([A-Z][a-wy-z]{1,3})[0-9,]*)")
	ula := ULAtomsRe.FindAllStringSubmatch(f, -1)

	// detecting wrong UL atoms
	// counting atoms at the same time and leaving on duplicates
	atomcount := make(map[string]int)
	for _, a := range ula {
		// wrong?
		if _, ok := globals.Atoms[a[2]]; !ok {
			return "", errors.New("wrong UL atom in formula: " + a[2])
		}
		upperLowerAtoms = append(upperLowerAtoms, a[0])
		// duplicate?
		if _, ok := atomcount[a[2]]; !ok {
			atomcount[a[2]] = 0
		} else {
			// atom already present !
			return "", errors.New("duplicate UL atom in formula")
		}
		// removing from formula for the next steps
		f = strings.Replace(f, a[0], "", -1)
	}
	if len(upperLowerAtoms) > 0 {
		hasULatom = true
	}

	// here we should have only one uppercase letter (and digits) per atom for the rest of
	// the formula

	// searching the C atom
	CAtomRe := regexp.MustCompile("((?:^[0-9]+)?(C)[0-9,]*)")
	ca := CAtomRe.FindAllStringSubmatch(f, -1)
	// will return [[C2 C2 C]] for ClNaC2
	// leaving on duplicated C atom
	if len(ca) > 1 {
		return "", errors.New("duplicate C atom in formula")
	}
	if len(ca) == 1 {
		hasCatom = true
		// removing from formula for the next steps
		f = strings.Replace(f, ca[0][0], "", -1)
	}

	// searching the H atom
	HAtomRe := regexp.MustCompile("((?:^[0-9]+)?(H)[0-9,]*)")
	ha := HAtomRe.FindAllStringSubmatch(f, -1)
	// will return [[H2 H2 H]] for ClNaH2
	// leaving on duplicated C atom
	if len(ha) > 1 {
		return "", errors.New("duplicate H atom in formula")
	}
	if len(ha) == 1 {
		hasHatom = true
		// removing from formula for the next steps
		f = strings.Replace(f, ha[0][0], "", -1)
	}

	// searching the other atoms
	OAtomRe := regexp.MustCompile("((?:^[0-9]+)?([A-Z])[0-9,]*)")
	oa := OAtomRe.FindAllStringSubmatch(f, -1)

	// detecting wrong atoms
	// counting atoms at the same time and leaving on duplicates
	atomcount = make(map[string]int)
	for _, a := range oa {
		// wrong?
		if _, ok := globals.Atoms[a[2]]; !ok {
			return "", errors.New("wrong UL atom in formula: " + a[2])
		}
		otherAtoms = append(otherAtoms, a[0])
		// duplicate?
		if _, ok := atomcount[a[2]]; !ok {
			atomcount[a[2]] = 0
		} else {
			// atom already present !
			return "", errors.New("duplicate other atom in formula")
		}
		// removing from formula for the next steps
		f = strings.Replace(f, a[0], "", -1)
	}
	if len(oa) > 0 {
		hasOatom = true
	}

	// if formula is not emty, this is an error
	if len(f) != 0 {
		return "", errors.New("wrong lowercase atoms in formula")
	}

	// rebuilding the formula
	newf := ""
	if hasCatom {
		newf += ca[0][0]
	}
	if hasHatom {
		newf += ha[0][0]
	}
	if hasOatom || hasULatom {
		at := append(otherAtoms, upperLowerAtoms...)
		sort.Strings(at)
		for _, a := range at {
			newf += a
		}
	}

	return newf, nil
}

// RandStringBytes generates a n size random string
func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = globals.LetterBytes[rand.Intn(len(globals.LetterBytes))]
	}
	return string(b)
}

// TestMail send a mail to "to"
func TestMail(to string) error {
	return SendMail(to, "test mail from Chimithèque", "your mail configuration seems ok")
}

// SendMail send a mail
func SendMail(to string, subject string, body string) error {

	var (
		e         error
		tlsconfig *tls.Config
		tlsconn   *tls.Conn
		client    *smtp.Client
		smtpw     io.WriteCloser
		n         int64
		message   string
	)

	// build message
	message += fmt.Sprintf("From: %s\r\n", globals.MailServerSender)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "\r\n" + body

	globals.Log.WithFields(logrus.Fields{
		"globals.MailServerAddress":       globals.MailServerAddress,
		"globals.MailServerPort":          globals.MailServerPort,
		"globals.MailServerSender":        globals.MailServerSender,
		"globals.MailServerUseTLS":        globals.MailServerUseTLS,
		"globals.MailServerTLSSkipVerify": globals.MailServerTLSSkipVerify,
		"subject":                         subject,
		"to":                              to}).Debug("sendMail")

	if globals.MailServerUseTLS {
		// tls
		tlsconfig = &tls.Config{
			InsecureSkipVerify: globals.MailServerTLSSkipVerify,
			ServerName:         globals.MailServerAddress,
		}
		if tlsconn, e = tls.Dial("tcp", globals.MailServerAddress+":"+globals.MailServerPort, tlsconfig); e != nil {
			return e
		}
		defer tlsconn.Close()
		if client, e = smtp.NewClient(tlsconn, globals.MailServerAddress+":"+globals.MailServerPort); e != nil {
			return e
		}
	} else {
		if client, e = smtp.Dial(globals.MailServerAddress + ":" + globals.MailServerPort); e != nil {
			return e
		}
	}
	defer client.Close()

	// to && from
	globals.Log.Debug("setting sender")
	if e = client.Mail(globals.MailServerSender); e != nil {
		return e
	}
	globals.Log.Debug("setting recipient")
	if e = client.Rcpt(to); e != nil {
		return e
	}
	// data
	globals.Log.Debug("setting body")
	if smtpw, e = client.Data(); e != nil {
		return e
	}
	defer smtpw.Close()

	// send message
	globals.Log.Debug("sending message")
	buf := bytes.NewBufferString(message)
	if n, e = buf.WriteTo(smtpw); e != nil {
		return e
	}
	globals.Log.WithFields(logrus.Fields{"n": n}).Debug("sendMail")

	// send quit command
	globals.Log.Debug("setting quit command")
	_ = client.Quit()

	return nil
}

// TimeTrack displays the run time of the function "name"
// from the start time "start"
// use: defer utils.TimeTrack(time.Now(), "GetProducts")
// at the begining of the function to track
func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	globals.Log.Debug(fmt.Sprintf("%s took %s", name, elapsed))
}

// GetPasswordHash return password hash for the login,
func GetPasswordHash(login string) ([]byte, error) {

	h := hmac.New(sha256.New, []byte("secret"))
	if _, err := h.Write([]byte(login)); err != nil {
		return nil, err
	}

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
		globals.Log.Debug(fmt.Sprintf("record: %s", record))
		if err == io.EOF {
			break
		}
		if err != nil {
			globals.Log.Fatal(err)
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

// ProductsToCSV returns a file name of the products prs
// exported into CSV
func ProductsToCSV(prs []Product) string {

	header := []string{"product_id",
		"product_name",
		"product_synonyms",
		"product_cas",
		"product_ce",
		"product_specificity",
		"empirical_formula",
		"linear_formula",
		"3D_formula",
		"MSDS",
		"class_of_compounds",
		"physical_state",
		"signal_word",
		"symbols",
		"hazard_statements",
		"precautionary_statements",
		"remark",
		"disposal_comment",
		"restricted?",
		"radioactive?"}

	// create a temp file
	tmpFile, err := ioutil.TempFile(os.TempDir(), "chimitheque-")
	if err != nil {
		globals.Log.Error("cannot create temporary file", err)
	}
	// creates a csv writer that uses the io buffer
	csvwr := csv.NewWriter(tmpFile)
	// write the header
	_ = csvwr.Write(header)
	for _, p := range prs {
		_ = csvwr.Write(p.ProductToStringSlice())
	}

	csvwr.Flush()
	return strings.Split(tmpFile.Name(), "chimitheque-")[1]
}

// StoragesToCSV returns a file name of the products prs
// exported into CSV
func StoragesToCSV(sts []Storage) (string, error) {

	var (
		err     error
		tmpFile *os.File
	)

	header := []string{"storage_id",
		"product_name",
		"product_casnumber",
		"product_specificity",
		"storelocation",
		"quantity",
		"unit",
		"barecode",
		"supplier",
		"creation_date",
		"modification_date",
		"entry_date",
		"exit_date",
		"opening_date",
		"expiration_date",
		"comment",
		"reference",
		"batch_number",
		"to_destroy?",
		"archive?"}

	// create a temp file
	if tmpFile, err = ioutil.TempFile(os.TempDir(), "chimitheque-"); err != nil {
		globals.Log.Error("cannot create temporary file", err)
		return "", err
	}
	// creates a csv writer that uses the io buffer
	csvwr := csv.NewWriter(tmpFile)
	// write the header
	if err = csvwr.Write(header); err != nil {
		globals.Log.Error("cannot write header", err)
		return "", err
	}

	for _, s := range sts {
		if err = csvwr.Write(s.StorageToStringSlice()); err != nil {
			globals.Log.Error("cannot write entry", err)
			return "", err
		}
	}

	csvwr.Flush()

	return strings.Split(tmpFile.Name(), "chimitheque-")[1], nil
}
