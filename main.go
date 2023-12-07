package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", GetCepHandler)
	fmt.Println("Running server on port:8080")
	http.ListenAndServe(":8080", nil)
}

func GetCepHandler(w http.ResponseWriter, r *http.Request) {
	cOne := make(chan string)
	cTwo := make(chan string)

	cep := r.URL.Query().Get("cep")
	if cep == "" {
		http.Error(w, "CEP is required", http.StatusBadRequest)
		return
	}

	go callBrasilAPI(cep, cOne)
	go callViaCep(cep, cTwo)

	select {
	case callOne := <-cOne:
		fmt.Fprintf(os.Stdout, "Quickest response BrasilAPI: %s \n", callOne)
		fmt.Fprintf(w, "Quickest response BrasilAPI: %s \n", callOne)
	case callTwo := <-cTwo:
		fmt.Fprintf(os.Stdout, "Quickest response ViaCep: %s \n", callTwo)
		fmt.Fprintf(w, "Quickest response ViaCep: %s \n", callTwo)
	case <-time.After(time.Second * 1):
		fmt.Fprintf(os.Stdout, "Timeout reached\n")
		fmt.Fprintf(w, "Timeout reached\n")
	}
}

func callBrasilAPI(cep string, c chan string) {
	req, err := http.NewRequest("GET", "https://brasilapi.com.br/api/cep/v1/"+cep, nil)
	if err != nil {
		fmt.Println("Error creating request")
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error on request")
		return
	}
	defer res.Body.Close()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body")
		return
	}
	c <- string(resp)
}

func callViaCep(cep string, c chan string) {
	req, err := http.NewRequest("GET", "http://viacep.com.br/ws/"+cep+"/json/", nil)
	if err != nil {
		fmt.Println("Error creating request")
		return
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error on request")
		return
	}
	defer res.Body.Close()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error reading response body")
		return
	}
	c <- string(resp)
}
