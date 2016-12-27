package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

/* Generic RPC */

func makeHttpClient(skipCheckCert bool) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipCheckCert},
	}
	client := &http.Client{Transport: tr}
	return client
}

func rpcCall(url string, data url.Values) (reply *bytes.Buffer, err error) {
	client := makeHttpClient(flagNoCheckCertificate)

	resp, err := client.PostForm(url, data)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	strReply := readReply(resp)

	return strReply, err
}

func readReply(reply *http.Response) *bytes.Buffer {
	buffer := bytes.NewBufferString("")
	scanner := bufio.NewScanner(reply.Body)
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)
	i := 0
	for scanner.Scan() {
		buffer.WriteString(scanner.Text())
		i++
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "There was an error with the scanner in attached container", err)
	}
	return buffer
}

/* ruTorrent Methods*/

const JSON_FILE_NAME int = 4
const JSON_FILE_PATH int = 25

type JsonTorrentList struct {
	T   map[string][]string
	Cid int
}

type Torrent struct {
	Hash     string
	FileName string
	FilePath string
}

type RuTorrent struct {
	ApiUrl string
}

func (rto RuTorrent) ListTorrent() []Torrent {
	args := url.Values{}
	args.Add("mode", "list")
	reply, err := rpcCall(rto.ApiUrl, args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	var jsonStruct JsonTorrentList
	json.Unmarshal(reply.Bytes(), &jsonStruct)
	torrents := make([]Torrent, len(jsonStruct.T))

	i := 0
	for k, v := range jsonStruct.T {
		fileName := v[JSON_FILE_NAME]
		filePath := v[JSON_FILE_PATH]

		torrents[i] = Torrent{Hash: k, FileName: fileName, FilePath: filePath}
		i++
	}

	return torrents
}

/* Download Directory */

type DownloadFile struct {
	FileName string
	FilePath string
}

type DownloadDirectory struct {
	Path string
}

func (dir DownloadDirectory) ListFile() []DownloadFile {
	reply, err := filepath.Glob(dir.Path + "/*")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	files := make([]DownloadFile, len(reply))

	i := 0
	for _, f := range reply {
		files[i] = DownloadFile{FileName: filepath.Base(f), FilePath: f}
		i++
	}

	return files
}

/* Utils */

func findTorrent(torrents *[]Torrent, file *DownloadFile) *Torrent {
	for _, t := range *torrents {
		if file.FileName == t.FileName && file.FilePath == t.FilePath {
			return &t
		}
	}
	return nil
}

func findFilesNoCorrespondingTorrent(torrents *[]Torrent, files *[]DownloadFile) *[]DownloadFile {
	res := make([]DownloadFile, 0, 32)

	for _, file := range *files {
		found := findTorrent(torrents, &file)
		if found == nil {
			res = append(res, file)
		}
	}

	return &res
}

/* */

const RTORRENT_CLEANER_VERSION = "0.1.0"

func displayAbout() {
	fmt.Printf("\n")
	fmt.Printf("rtorrent-cleaner - version %s\n", RTORRENT_CLEANER_VERSION)
	fmt.Printf("Licence: GNU GENERAL PUBLIC LICENSE Version 3\n")
	fmt.Printf("Developped by Anthony Pena <https://k49.fr.nf>\n")
	fmt.Printf("\n")
}

var flagRuTorrentUrl string
var flagDownloadDir string
var flagAbsolutePath bool
var flagNoCheckCertificate bool

var flagAbout bool

func initFlag() {

	flag.StringVar(&flagRuTorrentUrl, "ruTorrent", "", "The URL of the ruTorrent instance to use. Exp: https://localhost")
	flag.StringVar(&flagDownloadDir, "dir", "", "The directory where torrent are downloaded. Exp: /home/user/download")
	flag.BoolVar(&flagAbsolutePath, "absolute-path", false, "Enable absolute path output. Repl. Default: file name is displayed.")
	flag.BoolVar(&flagNoCheckCertificate, "no-check-certificate", false, "Disable certificate checking.")

	flag.BoolVar(&flagAbout, "version", false, "Display informations about rtorrent-cleaner.")
	flag.Parse()

	if flagAbout {
		displayAbout()
		os.Exit(0)
	}

	err := false
	if flagRuTorrentUrl == "" {
		fmt.Fprintln(os.Stderr, "ruTorrent parameter is mandatory")
		err = true
	}
	if flagDownloadDir == "" {
		fmt.Fprintln(os.Stderr, "dir parameter is mandatory")
		err = true
	}
	if err {
		os.Exit(1)
	}
}

func main() {

	initFlag()

	apiUrl := flagRuTorrentUrl + "/plugins/httprpc/action.php"
	path := flagDownloadDir

	rto := RuTorrent{ApiUrl: apiUrl}
	torrents := rto.ListTorrent()

	dlDir := DownloadDirectory{Path: path}
	files := dlDir.ListFile()

	fmt.Printf("RuTorrent api used: %s\n", apiUrl)
	fmt.Printf("Download directory used: %s\n", path)
	fmt.Printf("Torrents from ruTorrent: %d\n", len(torrents))
	fmt.Printf("Torrents from directory: %d\n", len(files))

	filesToRemove := findFilesNoCorrespondingTorrent(&torrents, &files)

	fmt.Printf("Files or Directories with not matching torrents : %d\n\n", len(*filesToRemove))
	for _, file := range *filesToRemove {
		if flagAbsolutePath {
			fmt.Printf("%s\n", file.FilePath)
		} else {
			fmt.Printf("%s\n", file.FileName)
		}
	}
}
