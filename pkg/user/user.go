package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gabeduke/wio-cli-go/internal"
	"github.com/howeyc/gopass"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"net"
	"net/http"
	"net/url"
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

func (c *credentials) Create(logger *log.Entry) (*CreateResponse, error) {
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

	ep := viper.GetString(internal.HOST) + "/v1/user/Create"

	req, err := http.NewRequest("POST", ep, bytes.NewBuffer(d))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte

	logger.WithField("status", resp.Status).Debug("Create")
	logger.WithField("headers", resp.Header).Trace("Create")

	var r CreateResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	logger.WithField("token", r.Token).Info("Create successful")

	return &r, nil
}

func (r *LoginResponse) Login(logger *log.Entry) error {
	var usr credentials

	usr.getEmail(logger)
	viper.Set(internal.EMAIL, usr.Email)

	err := usr.getPassword(logger)
	if err != nil {
		return err
	}

	d, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	ep := viper.GetString(internal.HOST) + "/v1/user/login"
	req, err := http.NewRequest("POST", ep, bytes.NewBuffer(d))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		return errors.Errorf("Login failed: %v", resp.Status)
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body) // response body is []byte

	logger.WithField("status", resp.Status).Debug("Login")
	logger.WithField("headers", resp.Header).Trace("Login")

	err = json.Unmarshal(body, r)
	if err != nil {
		return err
	}

	viper.Set(internal.TOKEN, r.Token)
	logger.WithField("token", r.Token).WithField("user_id", r.UserId).Info("Login successful")

	return nil
}

func (c *credentials) getPassword(logger *log.Entry) error {
	fmt.Printf("Enter Password: ")
	password, err := gopass.GetPasswd()

	c.Password = string(password)

	logger.Debugf("Password: %s", password)
	return err
}

func (c *credentials) getEmail(logger *log.Entry) {
	c.Email = viper.GetString("Email")

	if c.Email == "" {
		c.Email = internal.Prompt("Email Address: ", "")
	}

	logger.Infof("Email: %s", c.Email)
}

func configure(logger *log.Entry) error {
	logger.Debug("configure called")

	// Prompt for server address
	viper.Set(internal.HOST, internal.Prompt("Enter the server address (eg. https://wio.leetserve.com): ", viper.GetString(internal.HOST)))

	// Prompt for server IP
	mip := internal.Prompt("Enter the server IP address (leave blank to allow discovery): ", "")
	if mip == "" {
		host, err := url.Parse(viper.GetString(internal.HOST))
		if err != nil {
			return errors.Errorf("Error parsing server address: %v", err)
		}
		hostAddr, err := net.LookupIP(host.Hostname())
		if err != nil {
			fmt.Println("Unknown host")
		} else {
			mip = hostAddr[0].String()
		}
	}
	viper.Set(internal.HOST_IP, mip)

	u := LoginResponse{}
	u.Login(logger)

	viper.Set(internal.TOKEN, u.Token)

	logger.Debugf("Wio CLI Configuration: %v", viper.AllSettings())
	logger.WithField("file", viper.ConfigFileUsed()).Info("Wio CLI Configuration file updated")

	return viper.WriteConfigAs(viper.ConfigFileUsed())
}
