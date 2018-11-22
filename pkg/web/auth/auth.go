package web

// import (
// 	"crypto/sha256"
// 	"encoding/hex"
// 	"fmt"
// 	"html/template"
// 	"math/rand"
// 	"net/http"
// 	"strconv"
// 	"strings"
// 	"time"
// )

// // Reference https://austingwalters.com/building-a-web-server-in-go-web-cookies/
// // currently unused

// func SHA(str string) string {

// 	bytes := []byte(str)

// 	// Converts string to sha2
// 	h := sha256.New()                   // new sha256 object
// 	h.Write(bytes)                      // data is now converted to hex
// 	code := h.Sum(nil)                  // code is now the hex sum
// 	codestr := hex.EncodeToString(code) // converts hex to string

// 	return codestr
// }

// func LoginCookie(username string) http.Cookie {

// 	cookieValue := username + ":" + SHA(username+strconv.Itoa(rand.Intn(100000000)))
// 	expire := time.Now().AddDate(0, 0, 1)
// 	return http.Cookie{Name: "SessionID", Value: cookieValue, Expires: expire, HttpOnly: true}
// }

// // Check if the user is logged in, true if they do
// func IsLoggedIn(r *http.Request) bool {

// 	// Obtains cookie from users http.Request
// 	cookie, err := r.Cookie("SessionID")
// 	if err != nil {
// 		fmt.Println(err)
// 		return false
// 	}

// 	// Obtain sessionID from cookie we obtained earlier
// 	sessionID := cookie.Value

// 	// Split the sessionID to Username and ID (username+random)
// 	z := strings.Split(sessionID, ":")
// 	email := z[0]
// 	sessionID = z[1]

// 	// Returns the expectedSessionID from the database
// 	expectedSessionID, errz := lookupSessionID(email)

// 	if errz != "" {
// 		fmt.Println(errz)
// 		return false
// 	}

// 	// If SessionID matches the expected SessionID, it is Good
// 	if sessionID == expectedSessionID {
// 		// If you want to be really secure check IP
// 		return true
// 	}

// 	return false
// }

// // Shows user requested page
// func viewHandler(w http.ResponseWriter, r *http.Request) {

// 	// Parses URL to obtain title of file to add to .body
// 	title := r.URL.Path[len("/"):]

// 	// Load templatized page, given title
// 	p, err := loadPage(title, r)

// 	// If they do not have a cookie force them to login
// 	if err != nil && !cookies.IsLoggedIn(r) {
// 		http.Redirect(w, r, "/home", http.StatusFound)
// 		return

// 		// If they do have a cookie, they are logged in
// 	} else if err != nil {
// 		http.Redirect(w, r, "/login-succeeded", http.StatusFound)
// 		return
// 	}

// 	// Generate template t
// 	t, _ := template.ParseFiles("view.html")

// 	// Write the template attributes of p (from load page) to t
// 	t.Execute(w, p)
// }
