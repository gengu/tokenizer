package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	packageName = "codec"
)

type config struct {
	url       string
	mapName   string
	filename  string
	localpath string
}

type Tokenizer struct {
	Model Model `json:"model"`
}

type Model struct {
	Vocab map[string]int64 `json:"vocab"`
}

func main() {
	encoding := flag.String("encoding", "", "encoding format. (e.g. cl100k_base)")
	flag.Parse()

	if encoding == nil {
		flag.PrintDefaults()
		os.Exit(1)
	}

	cfg := getConfig(*encoding)

	file, err := os.Create(cfg.filename)
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer file.Close()

	generatePreable(file, *encoding)
	// 如果是从本地加载
	if cfg.localpath != "" {
		// 读本地文件，只拿里面的tokens
		vocab := readVocabularyFromFile(cfg.localpath)
		genVocabularyFromFile(file, cfg.mapName, vocab)
	} else {
		genVocabulary(file, cfg.mapName, cfg.url)
	}
}

func generatePreable(w io.Writer, encoding string) {
	fmt.Fprintf(w, "package %s\n", packageName)
	fmt.Fprintf(w, "//go:generate go run ../internal/cmd/vocab.go -encoding %s\n", encoding)
	fmt.Fprintf(w, "// THIS FILE WAS AUTOMATICALLY GENERATED. DO NOT MODIFY\n")
}

// readVocabularyFromFile read file to tokenizer map
func readVocabularyFromFile(filename string) map[string]int64 {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("error reading file: %v", err)
	}
	var tokenizer Tokenizer
	err = json.Unmarshal(fileContent, &tokenizer)
	if err != nil {
		log.Fatalf("error unmarshalling file: %v", err)
	}
	return tokenizer.Model.Vocab
}

// genVocabularyFromFile generate tokenizer map
func genVocabularyFromFile(w io.Writer, mapName string, m map[string]int64) {
	fmt.Fprintf(w, "var %s vocab = vocab{\n", mapName)
	for k, v := range m {
		fmt.Fprintf(w, "    %s:%s,\n", strconv.Quote(k), strconv.Quote(strconv.FormatInt(v, 10)))
	}
	fmt.Fprintf(w, "}\n\n")
}

func genVocabulary(w io.Writer, mapName string, uri string) {
	resp, err := http.Get(uri)
	if err != nil {
		log.Fatalf("error fetching file: %v", err)
	}
	defer resp.Body.Close()

	fmt.Fprintf(w, "var %s vocab = vocab{\n", mapName)

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")

		if len(parts) != 2 {
			log.Fatalf("invalid line: %s", line)
		}

		word, err := base64.StdEncoding.DecodeString(parts[0])
		if err != nil {
			log.Fatalf("invalid word: %s", parts[0])
		}

		fmt.Fprintf(w, "	%s:%s,\n", strconv.Quote(string(word)), parts[1])
	}

	fmt.Fprintf(w, "}\n\n")
}

func getConfig(encoding string) config {
	switch encoding {
	case "cl100k_base":
		return config{
			mapName:  "cl100kBaseVocab",
			url:      "https://openaipublic.blob.core.windows.net/encodings/cl100k_base.tiktoken",
			filename: "cl100k_base_vocab.go",
		}
	case "r50k_base":
		return config{
			mapName:  "r50kBaseVocab",
			url:      "https://openaipublic.blob.core.windows.net/encodings/r50k_base.tiktoken",
			filename: "r50k_base_vocab.go",
		}
	case "p50k_base":
		return config{
			mapName:  "p50kBaseVocab",
			url:      "https://openaipublic.blob.core.windows.net/encodings/p50k_base.tiktoken",
			filename: "p50k_base_vocab.go",
		}
	case "starcoder":
		return config{
			mapName:   "starcoderVocab",
			localpath: "../internal/resources/starcoder/tokenizer.json",
			filename:  "starcoder_base_vocab.go",
		}
	default:
		log.Fatal("config not found")
		return config{}
	}
}
