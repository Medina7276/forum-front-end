package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"forum/dto"
	"forum/model"
	"io/ioutil"
	"net/http"
	"strings"

	uuid "github.com/satori/go.uuid"
)

var client = &http.Client{}

func CreatePost(post model.Post) (*model.Post, error) {
	if err := checkPost(post); err != nil {
		return nil, err
	}

	requestBody, err := json.Marshal(&post)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post("localhost:8080/api/post", "application/json", bytes.NewBuffer(requestBody)) //send request to backend
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, getErrorFromBody(resp)
	}

	var res model.Post
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func GetPostById(id uuid.UUID) (*dto.PostDto, error) {
	resp, err := client.Get("localhost:8080/api/post/" + id.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, getErrorFromBody(resp)
	}

	var res dto.PostDto
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func GetAllPosts() ([]model.Post, error) {
	resp, err := client.Get("localhost:8080/api/post")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, getErrorFromBody(resp)
	}

	var res []model.Post
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func GetPostsBySubforumId(id uuid.UUID) ([]model.Post, error) {
	resp, err := client.Get("localhost:8080/api/post?subforumid=" + id.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, getErrorFromBody(resp)
	}

	var res []model.Post
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

func checkPost(post model.Post) error {
	if strings.TrimSpace(post.Title) == "" {
		return errors.New("empty title")
	}

	if strings.TrimSpace(post.Content) == "" {
		return errors.New("empty content")
	}

	if post.SubforumID == uuid.Nil {
		return errors.New("nil subforum_id")
	}

	return nil
}

func Login(creds *dto.AuthCredentials) (*dto.UserWithClaims, error) {
	if err := checkCreds(creds); err != nil {
		return nil, err
	}

	requestBody, err := json.Marshal(creds)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post("localhost:8080/login", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, getErrorFromBody(resp)
	}

	var res dto.UserWithClaims
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func checkCreds(creds *dto.AuthCredentials) error {
	if strings.TrimSpace(creds.Username) == "" {
		return errors.New("empty username")
	}

	if strings.TrimSpace(creds.Password) == "" {
		return errors.New("empty password")
	}

	return nil
}

func SignUp(user *model.User) error {
	if err := checkUser(user); err != nil {
		return err
	}

	requestBody, err := json.Marshal(user)
	if err != nil {
		return err
	}

	resp, err := client.Post("localhost:8080/register", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return getErrorFromBody(resp)
	}

	return nil
}

func CreateUser(username, email, password string) (*model.User, error) {
	user := &model.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	return user, checkUser(user)
}

func checkUser(user *model.User) error {
	fmt.Println(user.Username)
	if strings.TrimSpace(user.Username) == "" {
		return errors.New("username is empty")
	}

	if strings.TrimSpace(user.Password) == "" {
		return errors.New("password is empty")
	}

	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is empty")
	}

	// if strings.TrimSpace(user.Name) == "" {
	// 	return errors.New("name is empty")
	// }

	return validatePassword(user.Password)
}

func validatePassword(s string) error {
	if len(s) < 6 {
		return errors.New("Password should be at least 6 characters long")
	}

	legal := ""
	for i := '!'; i <= '~'; i++ {
		legal += string(i)
	}

	for _, r := range s {
		if !strings.ContainsRune(legal, r) {
			return errors.New("Password shold only contain english letters, numbers, or special characters ('.')")
		}
	}

	return nil
}

func CreateComment(comment *model.Comment) error {
	if err := checkComment(comment); err != nil {
		return nil
	}

	requestBody, err := json.Marshal(comment)
	if err != nil {
		return err
	}

	resp, err := client.Post("localhost:8080/api/post", "application/json", bytes.NewBuffer(requestBody)) //send request to backend
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return getErrorFromBody(resp)
	}

	return nil
}

func checkComment(comment *model.Comment) error {
	if strings.TrimSpace(comment.Content) == "" {
		return errors.New("content is empty")
	}

	return nil
}

func Rate(rate model.Like) error {
	client := http.Client{}

	requestBody, err := json.Marshal(&rate)
	if err != nil {
		return err
	}

	resp, err := client.Post("localhost:8080/rate", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return getErrorFromBody(resp)
	}

	return nil
}

func getErrorFromBody(resp *http.Response) error {
	bytes, _ := ioutil.ReadAll(resp.Body)
	return errors.New(string(bytes))
}
