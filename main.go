package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func load_domains() []string {
	file, err := os.Open("./data/domains.txt")
	if err != nil {
		return nil
	}
	defer file.Close()

	var domains []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		domains = append(domains, scanner.Text())
	}
	return domains
}

func validate_guess(partial string, comparison string) bool {
	if len(partial) != len(comparison) {
		return false
	}
	for i, c := range partial {
		if string(c) != string('*') {
			if string(comparison[i]) != string(c) {
				return false
			}
		}
	}
	return true
}

func guess_domain(c *gin.Context) {
	domains := load_domains()
	var possible []string
	for _, val := range domains {
		if validate_guess(string(strings.Split(c.Param("email"), "@")[1]), string(val)) {
			possible = append(possible, string(val))
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"email":            string(c.Param("email")),
		"possible_domains": possible,
	})
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	server := gin.Default()

	server.GET("/guess/:email", guess_domain)

	println("Starting server at http://0.0.0.0:8080")
	println("Make GET requests to /guess/:email")
	println("e.g curl http://localhost:8080/guess/test@g****.com")
	println("Sample Response: {\"email\":\"test@g****.com\",\"possible_domains\":[\"gmail.com\",\"games.com\",\"gawab.com\",\"globo.com\"]}")
	server.Run()
}
