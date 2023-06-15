package nodes

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gabeduke/wio-cli-go/internal"
	"github.com/spf13/viper"
	"io"
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

type Node struct {
	Name        string      `json:"name"`
	NodeKey     string      `json:"node_key"`
	NodeSn      string      `json:"node_sn"`
	Dataxserver interface{} `json:"dataxserver"`
	Board       string      `json:"board"`
	Online      bool        `json:"online"`
}

type ListResp struct {
	Nodes []Node `json:"nodes"`
}

func (l ListResp) String() string {
	b, err := json.MarshalIndent(l, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	return string(b)
}

type deleteResp struct {
	Result string `json:"result"`
}

func (c CreateResp) String() string {
	return fmt.Sprintf("key: %s\nserial name: %s", c.NodeKey, c.NodeSn)
}

func getURIFromConfig() (*url.URL, error) {
	return url.Parse(viper.GetString(internal.HOST))
}

func RegisterNode() error {

	if viper.GetBool("create") {
		if nodeName == "" {
			nodeName = internal.Prompt("Enter a name for your node: ", "")
		}

		if boardType == "" {
			boardTypeStr := internal.Prompt("Enter the board type (node or link): ", "link")
			boardType = boardEnum(boardTypeStr)
		}

		resp, err := CreateNode(nodeName, boardType)
		if err != nil {
			return err
		}

		viper.Set(internal.NODE_KEY, resp.NodeKey)
		viper.Set(internal.NODE_SN, resp.NodeSn)
	}

	fmt.Println("Registering node...")
	fmt.Println("To enter AP mode on the device: hold the `func` button for 5 seconds then connect to the AP from your WIFI network list")
	fmt.Print("Connect to device  then hit RETURN")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()

	p := make([]byte, 2048)
	conn, err := net.Dial("udp", internal.NODE_UDP_ADDR)
	if err != nil {
		return err
	}

	ssid := internal.Prompt("Enter the name of the SSID you want to connect to: ", "")
	pass := internal.Prompt("Enter the password for the SSID: ", "")

	r := fmt.Sprintf("APCFG: %s\t%s\t%s\t%s\t%s\t%s\t\r\n", ssid, pass, viper.GetString("key"), viper.GetString("sn"), viper.GetString(internal.HOST), viper.GetString(internal.HOST_IP))
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

func ListNodes() (ListResp, error) {
	nodes := ListResp{}
	ep, err := getURIFromConfig()
	if err != nil {
		return nodes, err
	}
	ep.Path = "/v1/nodes/list"

	req, err := http.NewRequest("GET", ep.String(), nil)
	req.Header.Add("Authorization", "token "+viper.GetString(internal.TOKEN))
	req.Header.Add("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nodes, err
	}

	if resp.StatusCode != http.StatusOK {
		return nodes, fmt.Errorf("failed to list nodes: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nodes, err
	}
	json.Unmarshal(bodyBytes, &nodes)

	return nodes, nil
}

func CreateNode(name string, boardType boardEnum) (CreateResp, error) {
	var board string
	switch boardType {
	case boardEnumNode:
		board = internal.WIO_NODE_V1_0
	case boardEnumLink:
		board = internal.WIO_LINK_V1_0
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

func DeleteNode(sn string) error {
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
	var deleteResp deleteResp
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
	req.Header.Add("Authorization", "token "+viper.GetString(internal.TOKEN))
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
