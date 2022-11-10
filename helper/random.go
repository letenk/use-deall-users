package helper

import (
	"math/rand"
	"strings"
	"time"
)

// Const alphabet for use random data with string
const alphabet = "abcdefghijklmnopqrstuvwxyz"

// Func init for first run
func init() {
	// Run rand.Seed
	rand.Seed(time.Now().UnixNano())
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	// Get total character on const alphabet
	k := len(alphabet)

	// Loop through n
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generates a random person name
func RandomPerson() string {
	// Create list person name
	persons := []string{"Ari", "Ayu", "Aulia", "Anggi", "Agus", "Ade", "Arya", "Amel", "Andi", "Bayu", "Bagas", "Budi", "Bagus", "Bastian", "Ben", "Chika", "Cinta", "Citra", "Cakra", "Candra", "Darius", "Dimas", "Deo", "Dean", "Dinda", "Dika", "Dodi", "Ernes", "Erwin", "Eka", "Elin", "Elsa", "Ema", "Ela", "Fikri", "Fitri", "Fika", "Fani", "Fina", "Farid", "Fadel", "Galih", "Gading", "Guntur", "Gilang", "Geri", "Gibran", "Hamidah", "Hilda", "Hilmi", "Hisyam", "Haikal", "Harun", "Ita", "Ilham", "Indra", "Ikbal", "Irwan", "Ivan", "Irfan", "Ian", "Joko", "Josua", "Jonathan", "Jeri", "Jefri", "Karin", "Kirana", "Keisya", "Kevin", "Keyla", "Luna", "Lala", "Larisa", "Latif", "Lukman", "Mila", "Monika", "Maya", "Mira", "Malik", "Nila", "Nanda", "Naila", "Nisa", "Niko", "Nida", "Oki", "Okta", "Omar", "Oskar", "Olivia", "Putra", "Putri", "Paul", "Pinkan", "Pedro", "Qiqi", "Qafi", "Qori", "Qamar", "Queen", "Rafi", "Rafa", "Ririn", "Riska", "Rian", "Salsa", "Sinta", "Syifa", "Syahrul", "Samuel", "Tika", "Tristan", "Tobi", "Toni", "Tina", "Ulfa", "Usman", "Ulya", "Utari", "Umi", "Vania", "Virza", "Vincent", "Valdo", "Vino", "Wisnu", "Wulan", "Winda", "William", "Wira", "Yuda", "Yuli", "Yolanda", "Yusron", "Yosep"}

	// Get length person
	n := len(persons)

	// Return person name with rand mucn as length persons
	return persons[rand.Intn(n)]
}
