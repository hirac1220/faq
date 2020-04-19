package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	language "cloud.google.com/go/language/apiv1"
	languagepb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

type Config struct {
	Env         string `envconfig:"ENV"`
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

// JSON-to-Go: Convert JSON to Go instantly https://mholt.github.io/json-to-go/
type Syntax struct {
	Sentences []struct {
		Text struct {
			Content string `json:"content"`
		} `json:"text"`
	} `json:"sentences"`
	Tokens []struct {
		Text struct {
			Content string `json:"content"`
		} `json:"text,omitempty"`
		PartOfSpeech struct {
			Tag    int `json:"tag"`
			Proper int `json:"proper"`
		} `json:"part_of_speech,omitempty"`
		DependencyEdge struct {
			HeadTokenIndex int `json:"head_token_index"`
			Label          int `json:"label"`
		} `json:"dependency_edge,omitempty"`
		Lemma string `json:"lemma"`
	} `json:"tokens"`
	Language string `json:"language"`
}

type FAQ struct {
	// ID       string   `json:"id"`
	// Question string   `json:"question"`
	QType  string   `json:"qtype"`
	Nouns  []string `json:"noun"`
	Adjs   []string `json:"adj"`
	Verbs  []string `json:"verb"`
	Answer string   `json:"answer"`
}

type Sample struct {
	ID   uint
	Text string
}

type Answer struct {
	ID   uint
	Text string
}

const (
	NOUN = 6
	ADJ  = 1
	VERB = 11
)

func get(w http.ResponseWriter, r *http.Request) {
	// Get all records
	samples := [...]Sample{{ID: 1, Text: "体を動かしたい"}, {ID: 2, Text: "リラックスしたい"}, {ID: 3, Text: "美味しいもの作りたい"}}

	// Encode json
	resp, err := json.Marshal(&samples)
	log.Printf("%s", resp)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Write response
	w.Write(resp)
}

func post(w http.ResponseWriter, r *http.Request) {
	// Read request body
	q := map[string]string{"text": ""}
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&q); err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	log.Printf("%s", q["text"])
	a := GetFAQ(q["text"])
	// a := "https://www.youtube.com/watch?v=NseLsoJcqXc"
	// a := ""
	answer := Answer{
		ID:   1,
		Text: a,
	}
	log.Printf("%v", answer)
	// Encode json
	resp, err := json.Marshal(&answer)
	if err != nil {
		log.Println(err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	// Write response
	w.Write(resp)
}

func GetFAQ(text string) string {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.Env == "local" {
		log.Println("cannot NLP API (it may running local)")
		return ""
	}
	ctx := context.Background()

	// Creates a client.
	client, err := language.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the text to analyze.
	fmt.Printf("{\"text\": \"%v\"}\n", text)
	start := time.Now()
	syntax, err := client.AnalyzeSyntax(ctx, &languagepb.AnalyzeSyntaxRequest{
		Document: &languagepb.Document{
			Source: &languagepb.Document_Content{
				Content: text,
			},
			Type: languagepb.Document_PLAIN_TEXT,
		},
		EncodingType: languagepb.EncodingType_UTF8,
	})
	end := time.Now()
	fmt.Printf("%fsec\n", (end.Sub(start)).Seconds())
	if err != nil {
		log.Fatalf("Failed to analyze syntax text: %v", err)
	}
	jsonsyntax, err := json.Marshal(syntax) // json encode
	resultsyntax := string(jsonsyntax)
	fmt.Println(resultsyntax)

	var syntaxs Syntax
	if err := json.Unmarshal(jsonsyntax, &syntaxs); err != nil {
		log.Fatal(err)
	}

	// perse each tags and contents
	var nouns, adjs, verbs []string
	for _, t := range syntaxs.Tokens {
		fmt.Printf("{\"tag\" : \"%d\"}\n", t.PartOfSpeech.Tag)
		tag := t.PartOfSpeech.Tag
		if tag == NOUN {
			nouns = append(nouns, t.Text.Content)
			fmt.Printf("{\"noun\" : \"%s\"}\n", t.Text.Content)
		} else if tag == ADJ {
			adjs = append(adjs, t.Text.Content)
			fmt.Printf("{\"adj\" : \"%s\"}\n", t.Text.Content)
		} else if tag == VERB {
			verbs = append(verbs, t.Text.Content)
			fmt.Printf("{\"verb\" : \"%s\"}\n", t.Text.Content)
		}
	}
	f := &FAQ{
		Nouns: nouns,
		Adjs:  adjs,
		Verbs: verbs,
	}
	fmt.Printf("%v", f)

	f, err = SearchFAQ(ctx, *f)
	if err != nil {
		fmt.Print("something wrong")
		return ""
	} else if f == nil {
		log.Println("cannot search!")
		return ""
	}

	fmt.Printf("{\"qtype\" : \"%s\"}\n", f.QType)
	fmt.Printf("{\"answer\" : \"%s\"}\n", f.Answer)

	return f.Answer
}

func SearchFAQ(ctx context.Context, faq FAQ) (*FAQ, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	// Connect DB
	db, err := sql.Open("sqlite3", cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Database connection failed: ", err)
	} else {
		log.Printf("db: %v", cfg.DatabaseURL)
	}
	defer db.Close()

	var cmd string
	f := &FAQ{}
	if len(faq.Adjs) == 0 && len(faq.Verbs) == 0 {
		for _, n := range faq.Nouns {
			fmt.Print("no adj and no verbs\n")
			cmd = "SELECT qtype, answer FROM faq " +
				"WHERE (noun LIKE ?) " +
				"ORDER BY id DESC"
			// cmd = "SELECT qtype, answer FROM faq;"
			// rows, err := db.Query(cmd)
			rows, err := db.Query(cmd, n)
			if rows == nil {
				return nil, nil
			}
			if f, err = scanRows(rows); err != nil {
				fmt.Print("scan error")
			}
			if f != nil {
				return f, nil
			}
		}
	} else if len(faq.Nouns) > 0 && len(faq.Adjs) != 0 {
		fmt.Print("check adj\n")
		for _, n := range faq.Nouns {
			for _, a := range faq.Adjs {
				cmd = "SELECT qtype, answer FROM faq " +
					"WHERE (noun LIKE ? OR adj LIKE ?) " +
					"ORDER BY id DESC;"
				rows, err := db.Query(cmd, n, a)
				if rows == nil {
					return nil, nil
				}
				if f, err = scanRows(rows); err != nil {
					fmt.Print("scan error")
				}
				if f != nil {
					return f, nil
				}
			}
		}
	} else if len(faq.Nouns) == 0 && len(faq.Adjs) != 0 {
		fmt.Print("check adj only\n")
		for _, a := range faq.Adjs {
			cmd = "SELECT qtype, answer FROM faq " +
				"WHERE adj LIKE ? " +
				"ORDER BY id DESC;"
			rows, err := db.Query(cmd, a)
			if rows == nil {
				return nil, nil
			}
			if f, err = scanRows(rows); err != nil {
				fmt.Print("scan error")
			}
			if f != nil {
				return f, nil
			}
		}
	} else if len(faq.Nouns) > 0 && len(faq.Verbs) > 0 {
		fmt.Print("check verbs\n")
		for _, n := range faq.Nouns {
			for _, v := range faq.Verbs {
				cmd = "SELECT qtype, answer FROM faq " +
					"WHERE (noun LIKE ? OR verb LIKE ?) " +
					"ORDER BY id DESC;"
				rows, err := db.Query(cmd, n, v)
				if rows == nil {
					return nil, nil
				}
				if f, err = scanRows(rows); err != nil {
					fmt.Print("scan error")
				}
				if f != nil {
					return f, nil
				}
			}
		}
	} else if len(faq.Nouns) == 0 && len(faq.Verbs) > 0 {
		fmt.Print("check verbs only\n")
		for _, v := range faq.Verbs {
			cmd = "SELECT qtype, answer FROM faq " +
				"WHERE verb LIKE ? " +
				"ORDER BY id DESC;"
			rows, err := db.Query(cmd, v)
			if rows == nil {
				return nil, nil
			}
			if f, err = scanRows(rows); err != nil {
				fmt.Print("scan error")
			}
			if f != nil {
				return f, nil
			}
		}
	}
	return nil, err
}

func scanRows(rows *sql.Rows) (*FAQ, error) {
	f := &FAQ{}
	for rows.Next() {
		err := rows.Scan(&f.QType, &f.Answer)
		if err != nil {
			checkError(f, err)
			return nil, err
		}
		fmt.Printf("qtype: %s, answer: %s\n", f.QType, f.Answer)
	}
	err := rows.Err()
	if err != nil {
		checkError(f, err)
		return nil, err
	}
	return f, nil
}

func checkError(f *FAQ, err error) {
	switch {
	case err == sql.ErrNoRows:
		fmt.Printf("Not found")
	case err != nil:
		panic(err)
	default:
		fmt.Printf("qtype: %s, answer: %s\n", f.QType, f.Answer)
	}
}
