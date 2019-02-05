package model

import (
	"fmt"
	"strings"

	"translator/providers/storage"

	log "github.com/sirupsen/logrus"
)

const (
	querySaveWord = `INSERT INTO words_%s 
							SET word = ? 
                  		ON DUPLICATE KEY UPDATE 
							id=LAST_INSERT_ID(id), 
							word = ?;`

	querySaveRelation = `INSERT IGNORE INTO ru_en (en_id, ru_id) VALUES %s`

	queryList = `SELECT
					en.word w,
					ru.word t
				FROM words_en en
				LEFT JOIN ru_en r ON r.en_id = en.id
				LEFT JOIN words_ru ru ON ru.id = r.ru_id
				WHERE 1=1
				ORDER BY en.id DESC
				LIMIT 50;`
)

type (
	Source struct {
		Lang   string
		WordID int64
	}

	Translates struct {
		Lang     string
		WordsIDs []int64
	}

	// WordMan ...
	WordMan struct {
		storage storage.IStorage
	}

	// Word struct
	Word struct {
		Word      string   `json:"word"`
		Translate []string `json:"translate"`
		Tags      []string `json:"tags,omitempty"`
	}

	ParamsList struct {
		Page, PerPage int
	}
)

// New ...
func New(storage storage.IStorage) *WordMan {
	return &WordMan{
		storage: storage,
	}
}

// List return words lists
func (w *WordMan) List(params ParamsList) ([]Word, error) {

	rows, err := w.storage.MySQL().Query(queryList)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	updateList := func(word, translate string, list []Word) []Word {
		for i, w := range list {
			if w.Word == word {
				list[i].Translate = append(list[i].Translate, translate)
				return list
			}
		}

		list = append(list, Word{
			Word:      word,
			Translate: []string{translate},
		})

		return list
	}

	var list []Word

	for rows.Next() {
		var word, translate string
		if err := rows.Scan(&word, &translate); err != nil {
			log.Errorf("scan list error: %s", err.Error())
			continue
		}

		list = updateList(word, translate, list)
	}

	return list, nil
}

// Save word to storage
func (w *WordMan) Save(lang string, words []string) ([]int64, error) {

	var (
		lids []int64
		err  error
	)

	query := fmt.Sprintf(querySaveWord, lang)

	for _, word := range words {
		result, err := w.storage.MySQL().Exec(query, word, word)
		if err != nil {
			return lids, err
		}

		lid, err := result.LastInsertId()
		if err != nil || lid < 1 {
			return lids, err
		}

		lids = append(lids, lid)
	}

	return lids, err
}

// SaveRelations ...
func (w *WordMan) SaveRelations(source Source, translates Translates) error {

	var subSql []string
	for _, id := range translates.WordsIDs {
		subSql = append(subSql, fmt.Sprintf("(%d, %d)", source.WordID, id))
	}

	query := fmt.Sprintf(querySaveRelation, strings.Join(subSql, ","))
	_, err := w.storage.MySQL().Exec(query)
	if err != nil {
		return err
	}

	return nil
}
