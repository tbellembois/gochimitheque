// Package utils provides useful chemistry related utilities
package utils

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/smtp"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tbellembois/gochimitheque/global"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var (
	atoms = map[string]string{
		"H":  "hydrogen",
		"He": "helium",
		"Li": "lithium",
		"Be": "berylium",
		"B":  "boron",
		"C":  "carbon",
		"N":  "nitrogen",
		"O":  "oxygen",
		"F":  "fluorine",
		"Ne": "neon",
		"Na": "sodium",
		"Mg": "magnesium",
		"Al": "aluminium",
		"Si": "silicon",
		"P":  "phosphorus",
		"S":  "sulfure",
		"Cl": "chlorine",
		"Ar": "argon",
		"K":  "potassium",
		"Ca": "calcium",
		"Sc": "scandium",
		"Ti": "titanium",
		"V":  "vanadium",
		"Cr": "chromium",
		"Mn": "manganese",
		"Fe": "iron",
		"Co": "cobalt",
		"Ni": "nickel",
		"Cu": "copper",
		"Zn": "zinc",
		"Ga": "gallium",
		"Ge": "germanium",
		"As": "arsenic",
		"Se": "sefeniuo",
		"Br": "bromine",
		"Kr": "krypton",
		"Rb": "rubidium",
		"Sr": "strontium",
		"Y":  "yltrium",
		"Zr": "zirconium",
		"Nb": "niobium",
		"Mo": "molybdenum",
		"Tc": "technetium",
		"Ru": "ruthenium",
		"Rh": "rhodium",
		"Pd": "palladium",
		"Ag": "silver",
		"Cd": "cadmium",
		"In": "indium",
		"Sn": "tin",
		"Sb": "antimony",
		"Te": "tellurium",
		"I":  "iodine",
		"Xe": "xenon",
		"Cs": "caesium",
		"Ba": "barium",
		"Hf": "hafnium",
		"Ta": "tantalum",
		"W":  "tungsten",
		"Re": "rhenium",
		"Os": "osmium",
		"Ir": "iridium",
		"Pt": "platinium",
		"Au": "gold",
		"Hg": "mercury",
		"Tl": "thallium",
		"Pb": "lead",
		"Bi": "bismuth",
		"Po": "polonium",
		"At": "astatine",
		"Rn": "radon",
		"Fr": "francium",
		"Ra": "radium",
		"Rf": "rutherfordium",
		"Db": "dubnium",
		"Sg": "seaborgium",
		"Bh": "bohrium",
		"Hs": "hassium",
		"Mt": "meitnerium",
		"Ds": "darmstadtium",
		"Rg": "roentgenium",
		"Cn": "copemicium",
		"La": "lanthanum",
		"Ce": "cerium",
		"Pr": "praseodymium",
		"Nd": "neodymium",
		"Pm": "promethium",
		"Sm": "samarium",
		"Eu": "europium",
		"Gd": "gadolinium",
		"Tb": "terbium",
		"Dy": "dysprosium",
		"Ho": "holmium",
		"Er": "erbium",
		"Tm": "thulium",
		"Yb": "ytterbium",
		"Lu": "lutetium",
		"Ac": "actinium",
		"Th": "thorium",
		"Pa": "protactinium",
		"U":  "uranium",
		"Np": "neptunium",
		"Pu": "plutonium",
		"Am": "americium",
		"Cm": "curium",
		"Bk": "berkelium",
		"Cf": "californium",
		"Es": "einsteinium",
		"Fm": "fermium",
		"Md": "mendelevium",
		"No": "nobelium",
		"Lr": "lawrencium",
		"D":  "deuterium",
	}

	// general formula regex
	formulaRe *regexp.Regexp

	// basic molecule regex (atoms and numbers only)
	basicMolRe *regexp.Regexp

	// (AYZ)n molecule like regex
	oneGroupMolRe *regexp.Regexp
)

// atomByLength is a string slice sorter.
type atomByLength []string

func (s atomByLength) Len() int           { return len(s) }
func (s atomByLength) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s atomByLength) Less(i, j int) bool { return len(s[i]) > len(s[j]) }

func init() {
	// general formula regex
	formulaRe = regexp.MustCompile("[A-Za-z0-9,\\^]+")

	sortedAtoms := make([]string, 0, len(atoms))
	for k := range atoms {
		sortedAtoms = append(sortedAtoms, k)
	}
	// the atom must be sorted by decreasing size
	// to match first Cl before C for example
	sort.Sort(atomByLength(sortedAtoms))

	// building the basic molecule regex
	// (atom1|atom2|...)([1-9]*)
	var buf bytes.Buffer
	buf.WriteString("(")
	for _, a := range sortedAtoms {
		buf.WriteString(a)
		buf.WriteString("|")
	}
	// removing the last |
	buf.Truncate(buf.Len() - 1)
	buf.WriteString(")")
	buf.WriteString("([1-9]+)*")
	basicMolRe = regexp.MustCompilePOSIX(buf.String())

	// building the one group molecule regex
	buf.Reset()
	buf.WriteString("(?:\\(|\\[)")
	buf.WriteString("((?:[")
	for _, a := range sortedAtoms {
		buf.WriteString(a)
		buf.WriteString("|")
	}
	// removing the last |
	buf.Truncate(buf.Len() - 1)
	buf.WriteString("]+[1-9]*)+)")
	buf.WriteString("(?:\\)|\\])")
	buf.WriteString("([1-9]*)")

	oneGroupMolRe = regexp.MustCompile(buf.String())
}

// IsCeNumber returns true if c is a valid ce number
func IsCeNumber(c string) bool {

	var (
		err                error
		checkdigit, checkd int
	)

	if c == "000-000-0" {
		return true
	}

	// compiling regex
	r := regexp.MustCompile("^(?P<groupone>[0-9]{3})-(?P<grouptwo>[0-9]{3})-(?P<groupthree>[0-9]{1})$")
	// finding group names
	n := r.SubexpNames()
	// finding matches
	ms := r.FindAllStringSubmatch(c, -1)
	if len(ms) == 0 {
		return false
	}
	m := ms[0]
	// then building a map of matches
	md := map[string]string{}
	for i, j := range m {
		md[n[i]] = j
	}

	if len(m) > 0 {
		numberpart := md["groupone"] + md["grouptwo"]

		// converting the check digit into int
		if checkdigit, err = strconv.Atoi(string(md["groupthree"])); err != nil {
			return false
		}

		// calculating the check digit
		counter := 1  // loop counter
		currentd := 0 // current processed digit in c

		for i := 0; i < len(numberpart); i++ {
			// converting digit into int
			if currentd, err = strconv.Atoi(string(numberpart[i])); err != nil {
				return false
			}
			checkd += counter * currentd
			counter++
			//fmt.Printf("counter: %d currentd: %d checkd: %d\n", counter, currentd, checkd)
		}
	}

	return checkd%11 == checkdigit
}

// IsCasNumber returns true if c is a valid cas number
func IsCasNumber(c string) bool {

	var (
		err                error
		checkdigit, checkd int
	)

	if c == "0000-00-0" {
		return true
	}

	// compiling regex
	r := regexp.MustCompile("^(?P<groupone>[0-9]{1,7})-(?P<grouptwo>[0-9]{2})-(?P<groupthree>[0-9]{1})$")
	// finding group names
	n := r.SubexpNames()
	// finding matches
	ms := r.FindAllStringSubmatch(c, -1)
	if len(ms) == 0 {
		return false
	}
	m := ms[0]
	// then building a map of matches
	md := map[string]string{}
	for i, j := range m {
		md[n[i]] = j
	}

	if len(m) > 0 {
		numberpart := md["groupone"] + md["grouptwo"]

		// converting the check digit into int
		if checkdigit, err = strconv.Atoi(string(md["groupthree"])); err != nil {
			return false
		}
		//fmt.Printf("checkdigit: %d\n", checkdigit)

		// calculating the check digit
		counter := 1  // loop counter
		currentd := 0 // current processed digit in c

		for i := len(numberpart) - 1; i >= 0; i-- {
			// converting digit into int
			if currentd, err = strconv.Atoi(string(numberpart[i])); err != nil {
				return false
			}
			checkd += counter * currentd
			counter++
			//fmt.Printf("counter: %d currentd: %d checkd: %d\n", counter, currentd, checkd)
		}
	}
	return checkd%10 == checkdigit
}

// LinearToEmpiricalFormula returns the empirical formula from the linear formula f.
// example: [(CH3)2SiH]2NH
//          (CH3)2C[C6H2(Br)2OH]2
func LinearToEmpiricalFormula(f string) string {
	var ef string

	s := "-"
	nf := ""

	// Finding the first (XYZ)n match
	reg := oneGroupMolRe

	for s != "" {
		s = reg.FindString(f)

		// Counting the atoms and rebuilding the molecule string
		m := oneGroupAtomCount(s)
		ms := "" // molecule string
		for k, v := range m {
			ms += k
			if v != 1 {
				ms += fmt.Sprintf("%d", v)
			}
		}

		// Then replacing the match with the molecule string - nf is for "new f"
		nf = strings.Replace(f, s, ms, 1)
		f = nf
	}

	// Counting the atoms
	bAc := basicAtomCount(nf)

	// Sorting the atoms
	// C, H and then in alphabetical order
	var ats []string // atoms
	hasC := false    // C atom present
	hasH := false    // H atom present

	for k := range bAc {
		switch k {
		case "C":
			hasC = true
		case "H":
			hasH = true
		default:
			ats = append(ats, k)
		}
	}
	sort.Strings(ats)

	if hasH {
		ats = append([]string{"H"}, ats...)
	}
	if hasC {
		ats = append([]string{"C"}, ats...)
	}

	for _, at := range ats {
		ef += at
		nb := bAc[at]
		if nb != 1 {
			ef += fmt.Sprintf("%d", nb)
		}
	}

	return ef
}

// oneGroupAtomCount returns a count of the atoms of the f formula as a map.
// f must be a formula like (XYZ) (XYZ)n or [XYZ] [XYZ]n.
// example:
// (CH3)2 will return "C":2, "H":6
// CH3CH(NO2)CH3 will return "N":1 "O":2
// CH3CH(NO2)(CH3)2 will return "N":1 "O":2 - process only the first match
func oneGroupAtomCount(f string) map[string]int {
	var (
		// the result map
		c = make(map[string]int)
		r = oneGroupMolRe
	)
	// Looking for non matching molecules.
	if !r.MatchString(f) {
		return nil
	}

	// sl is a list of 3 elements like
	// [[(CH3Na6CCl5H)2 CH3Na6CCl5H 2]]
	sl := r.FindAllStringSubmatch(f, -1)
	basicMol := sl[0][1]
	multiplier, _ := strconv.Atoi(sl[0][2])

	// if there is no multiplier
	if multiplier == 0 {
		multiplier = 1
	}

	// counting the atoms
	aCount := basicAtomCount(basicMol)
	for at, nb := range aCount {
		c[at] = nb * multiplier
	}

	return c
}

// basicAtomCount returns a count of the atoms of the f formula as a map.
// f must be a basic formula with only atoms and numbers.
// example:
// C6H5COC6H4CO2H will return "C1":4, "H":10, "O":3
// CH3CH(NO2)CH3 will return Nil, parenthesis are not allowed
func basicAtomCount(f string) map[string]int {
	var (
		// the result map
		c   = make(map[string]int)
		r   = basicMolRe
		err error
	)
	// Looking for non matching molecules.
	if !r.MatchString(f) {
		return nil
	}

	// sl is a slice like [[Na Na ] [Cl Cl ] [C2 C 2] [Cl3 Cl 3]]
	// for f = NaClC2Cl3
	// [ matchingString capture1 capture2 ]
	// capture1 is the atom
	// capture2 is the its number
	sl := r.FindAllStringSubmatch(f, -1)
	for _, i := range sl {
		atom := i[1]
		var nbAtom int
		if i[2] != "" {
			nbAtom, err = strconv.Atoi(i[2])
			if err != nil {
				return nil
			}
		} else {
			nbAtom = 1
		}
		if _, ok := c[atom]; ok {
			c[atom] = c[atom] + nbAtom
		} else {
			c[atom] = nbAtom
		}
	}
	return c
}

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
	if !formulaRe.MatchString(f) {
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
		if _, ok := atoms[a[2]]; !ok {
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
		if _, ok := atoms[a[2]]; !ok {
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
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
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
	message += fmt.Sprintf("From: %s\r\n", global.MailServerSender)
	message += fmt.Sprintf("To: %s\r\n", to)
	message += fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z))
	message += fmt.Sprintf("Subject: %s\r\n", subject)
	message += "\r\n" + body

	global.Log.WithFields(logrus.Fields{
		"global.MailServerAddress":       global.MailServerAddress,
		"global.MailServerPort":          global.MailServerPort,
		"global.MailServerSender":        global.MailServerSender,
		"global.MailServerUseTLS":        global.MailServerUseTLS,
		"global.MailServerTLSSkipVerify": global.MailServerTLSSkipVerify,
		"subject":                        subject,
		"to":                             to}).Debug("sendMail")

	if global.MailServerUseTLS {
		// tls
		tlsconfig = &tls.Config{
			InsecureSkipVerify: global.MailServerTLSSkipVerify,
			ServerName:         global.MailServerAddress,
		}
		if tlsconn, e = tls.Dial("tcp", global.MailServerAddress+":"+global.MailServerPort, tlsconfig); e != nil {
			return e
		}
		defer tlsconn.Close()
		if client, e = smtp.NewClient(tlsconn, global.MailServerAddress); e != nil {
			return e
		}
	} else {
		if client, e = smtp.Dial(global.MailServerAddress + ":" + global.MailServerPort); e != nil {
			return e
		}
	}
	defer client.Close()

	// to && from
	if e = client.Mail(global.MailServerSender); e != nil {
		return e
	}
	if e = client.Rcpt(to); e != nil {
		return e
	}
	// data
	if smtpw, e = client.Data(); e != nil {
		return e
	}
	defer smtpw.Close()

	// send message
	buf := bytes.NewBufferString(message)
	if n, e = buf.WriteTo(smtpw); e != nil {
		return e
	}
	global.Log.WithFields(logrus.Fields{"n": n}).Debug("sendMail")

	// send quit command
	client.Quit()

	return nil
}
