// package utils provides useful chemistry related utilities
package utils

import (
	"bytes"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

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
	return true
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

	for k, _ := range bAc {
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
