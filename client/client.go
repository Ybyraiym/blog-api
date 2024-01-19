package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const serverURL = "http://localhost:8080" // адрес сервера

func main() {
	registerUser()
	authenticateUser()
	getBlogPosts()
}

func registerUser() {
	fmt.Println("Registering a new user...")
	url := serverURL + "/register"

	userData := map[string]string{
		"login":    "testuser",
		"username": "Test User",
		"password": "testpassword",
	}

	sendJSONRequest(url, "POST", userData)
}

func authenticateUser() {
	fmt.Println("Authenticating the user...")
	url := serverURL + "/authenticate"

	authData := map[string]string{
		"login":    "testuser",
		"password": "testpassword",
	}

	response := sendJSONRequest(url, "POST", authData)
	token := response["token"].(string)

	// Use the obtained token for subsequent requests
	getBlogPostsWithToken(token)
}

func getBlogPosts() {
	fmt.Println("Fetching blog posts...")
	url := serverURL + "/api/blogposts"

	sendRequest(url, "GET")
}

func getBlogPostsWithToken(token string) {
	fmt.Println("Fetching blog posts with token...")
	url := serverURL + "/api/blogposts"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response:", string(body))
}

func sendJSONRequest(url, method string, data map[string]string) map[string]interface{} {
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil
	}

	var responseMap map[string]interface{}
	err = json.Unmarshal(responseBody, &responseMap)
	if err != nil {
		fmt.Println("Error unmarshalling JSON response:", err)
		return nil
	}

	fmt.Println("Response:", responseMap)
	return responseMap
}

func sendRequest(url, method string) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// body, _ := io.ReadAll(resp.Body)
	// fmt.Println("Response:", string(body))

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Необработанный ответ:", string(body)) // Вывести необработанный ответ

}
