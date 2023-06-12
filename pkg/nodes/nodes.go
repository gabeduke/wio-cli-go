package nodes

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gabeduke/wio-cli-go/pkg/config"
	"github.com/spf13/viper"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type CreateResp struct {
	NodeKey string `json:"node_key"`
	NodeSn  string `json:"node_sn"`
}

type DeleteResp struct {
	Result string `json:"result"`
}

func (c CreateResp) String() string {
	return fmt.Sprintf("key: %s\nserial name: %s", c.NodeKey, c.NodeSn)
}

func getURIFromConfig() (*url.URL, error) {
	return url.Parse(viper.GetString(config.HOST))
}

func RegisterNode() error {

	fmt.Print("Connect to device then hit RETURN")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	p := make([]byte, 2048)
	conn, err := net.Dial("udp", config.NODE_UDP_ADDR)
	if err != nil {
		return err
	}

	r := fmt.Sprintf("APCFG: %s\t%s\t%s\t%s\t%s\t%s\t\r\n", viper.GetString("ssid"), viper.GetString("pass"), viper.GetString("key"), viper.GetString("sn"), viper.GetString(config.HOST), viper.GetString(config.HOST_IP))
	fmt.Println(r)
	_, err = fmt.Fprintf(conn, r)
	if err != nil {
		return err
	}

	_, err = bufio.NewReader(conn).Read(p)
	if err == nil {
		fmt.Println(string(p))
	} else {
		return err
	}
	conn.Close()

	return nil
}

func ListNodes() (*http.Response, error) {
	ep, err := getURIFromConfig()
	if err != nil {
		return &http.Response{}, err
	}
	ep.Path = "/v1/nodes/list"

	req, err := http.NewRequest("GET", ep.String(), nil)
	req.Header.Add("Authorization", "token "+viper.GetString(config.TOKEN))
	req.Header.Add("Accept", "application/json")
	return http.DefaultClient.Do(req)
}

func CreateNode(name string, boardType boardEnum) (CreateResp, error) {
	var board string
	switch boardType {
	case boardEnumNode:
		board = config.WIO_NODE_V1_0
	case boardEnumLink:
		board = config.WIO_LINK_V1_0
	}

	data := url.Values{
		"name":  {name},
		"board": {board},
	}

	ep, err := getURIFromConfig()
	if err != nil {
		return CreateResp{}, err
	}

	ep.Path = "/v1/nodes/create"

	resp, err := postRequest(data, ep)
	if err != nil {
		return CreateResp{}, err
	}

	var registerResp CreateResp
	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&registerResp)
	if err != nil {
		return CreateResp{}, err
	}

	return registerResp, nil
}

func deleteNode(sn string) error {
	data := url.Values{
		"node_sn": {sn},
	}

	ep, err := getURIFromConfig()
	if err != nil {
		return err
	}

	ep.Path = "/v1/nodes/delete"

	resp, err := postRequest(data, ep)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	var deleteResp DeleteResp
	err = json.NewDecoder(resp.Body).Decode(&deleteResp)
	if err != nil {
		return err
	}

	if deleteResp.Result != "ok" {
		return fmt.Errorf("failed to delete node: %s", sn)
	}

	return nil
}

func postRequest(data url.Values, ep *url.URL) (*http.Response, error) {
	req, err := http.NewRequest("POST", ep.String(), strings.NewReader(data.Encode()))
	if err != nil {
		return &http.Response{}, err
	}

	// set headers
	req.Header.Add("Authorization", "token "+viper.GetString(config.TOKEN))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	// send request with headers
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &http.Response{}, err
	}

	return resp, nil
}
