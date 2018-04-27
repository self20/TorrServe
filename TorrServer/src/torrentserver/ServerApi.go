package torrentserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"torrentserver/server"
)

func TorrentServerAdd(host, magnet string) (string, error) {
	reqUrl, err := joinUrl(host, "/torrent/add")
	if err != nil {
		return "", err
	}

	jsReq := new(server.TorrentJsonRequest)
	jsReq.Link = magnet

	buf, err := json.Marshal(jsReq)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(reqUrl, "application/json", bytes.NewReader(buf))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", getErrorResponse(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		return string(body), err
	}
	return "", err
}

func TorrentServerGet(host, hash string) (string, error) {
	reqUrl, err := joinUrl(host, "/torrent/get")
	if err != nil {
		return "", err
	}

	jsReq := new(server.TorrentJsonRequest)
	jsReq.Hash = hash

	buf, err := json.Marshal(jsReq)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(reqUrl, "application/json", bytes.NewReader(buf))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", getErrorResponse(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		return string(body), err
	}
	return "", err
}

func TorrentServerRem(host, hash string) (string, error) {
	reqUrl, err := joinUrl(host, "/torrent/rem")
	if err != nil {
		return "", err
	}

	jsReq := new(server.TorrentJsonRequest)
	jsReq.Hash = hash

	buf, err := json.Marshal(jsReq)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(reqUrl, "application/json", bytes.NewReader(buf))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", getErrorResponse(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		return string(body), err
	}
	return "", err
}

func TorrentServerCleanCache(host string) error {
	reqUrl, err := joinUrl(host, "/torrent/cleancache")
	if err != nil {
		return err
	}

	resp, err := http.Post(reqUrl, "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return getErrorResponse(resp)
	}

	return err
}

func TorrentServerList(host string) (string, error) {
	reqUrl, err := joinUrl(host, "/torrent/list")
	if err != nil {
		return "", err
	}

	resp, err := http.Post(reqUrl, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", getErrorResponse(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		return string(body), err
	}
	return "", err
}

func TorrentServerInfo(host, hash string) (string, error) {
	reqUrl, err := joinUrl(host, "/torrent/stat")
	if err != nil {
		return "", err
	}

	jsReq := new(server.TorrentJsonRequest)
	jsReq.Hash = hash

	buf, err := json.Marshal(jsReq)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(reqUrl, "application/json", bytes.NewReader(buf))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", getErrorResponse(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if body != nil {
		return string(body), err
	}
	return "", err
}

func TorrentServerEcho(host string) error {
	reqUrl, err := joinUrl(host, "/echo")
	if err != nil {
		return err
	}

	resp, err := http.Get(reqUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return getErrorResponse(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if body != nil && string(body) == "Ok" {
		return nil
	}
	return err
}

func TorrentServerReadSets(host string) (string, error) {
	reqUrl, err := joinUrl(host, "/settings/read")
	if err != nil {
		return "", err
	}

	resp, err := http.Post(reqUrl, "application/json", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", getErrorResponse(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func TorrentServerWriteSets(host string, sets string) error {
	reqUrl, err := joinUrl(host, "/settings/write")
	if err != nil {
		return err
	}

	resp, err := http.Post(reqUrl, "application/json", bytes.NewBufferString(sets))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return getErrorResponse(resp)
	}

	return nil
}

func joinUrl(base, path string) (string, error) {
	baseUrl, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	pathUrl, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	return baseUrl.ResolveReference(pathUrl).String(), nil
}

func getErrorResponse(resp *http.Response) error {
	buf, _ := ioutil.ReadAll(resp.Body)
	if len(buf) > 0 {
		type jserr struct {
			Message string `json:"message"`
		}
		errMsg := jserr{}
		err := json.Unmarshal(buf, &errMsg)
		if err == nil {
			return errors.New(errMsg.Message)
		}
		return errors.New(string(buf))
	}
	return errors.New(resp.Status)
}
