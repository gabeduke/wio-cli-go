package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gabeduke/wio-cli-go/pkg/config"
	"github.com/gabeduke/wio-cli-go/pkg/util"
	"github.com/howeyc/gopass"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net/http"
)

type LoginResponse struct {
	Token  string `json:"token"`
	UserId string `json:"user_id"`
}

type CreateResponse struct {
	Token string `json:"token"`
}

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *credentials) create(logger *log.Entry) (*CreateResponse, error) {
	logger.Debug("creating user")

	c.getEmail(logger)

	err := c.getPassword(logger)
	if err != nil {
		return &CreateResponse{}, err
	}

	d, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	ep := viper.GetString(config.HOST) + "/v1/user/create"

	req, err := http.NewRequest("POST", ep, bytes.NewBuffer(d))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte

	logger.WithField("status", resp.Status).Debug("create")
	logger.WithField("headers", resp.Header).Trace("create")

	var r CreateResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	logger.WithField("token", r.Token).Info("create successful")

	return &r, nil
}

func (r *LoginResponse) login(logger *log.Entry) error {
	var usr credentials

	usr.getEmail(logger)

	err := usr.getPassword(logger)
	if err != nil {
		return err
	}

	d, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	ep := viper.GetString(config.HOST) + "/v1/user/login"
	req, err := http.NewRequest("POST", ep, bytes.NewBuffer(d))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte

	logger.WithField("status", resp.Status).Debug("login")
	logger.WithField("headers", resp.Header).Trace("login")

	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}
	logger.WithField("token", r.Token).WithField("user_id", r.UserId).Info("login successful")

	return nil
}

func (c *credentials) getPassword(logger *log.Entry) error {
	fmt.Printf("Enter password: ")
	password, err := gopass.GetPasswd()

	c.Password = string(password)

	logger.Debugf("Password: %s", password)
	return err
}

func (c *credentials) getEmail(logger *log.Entry) {
	c.Email = viper.GetString("email")

	if c.Email == "" {
		c.Email = util.Prompt("Email Address: ", "")
	}

	logger.Infof("email: %s", c.Email)
}
