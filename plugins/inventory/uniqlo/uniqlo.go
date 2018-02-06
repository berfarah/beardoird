package uniqlo

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/botopolis/bot"
	slacker "github.com/botopolis/slack"
	"github.com/nlopes/slack"
)

const uniqloStockURL = "https://www.uniqlo.com/on/demandware.store/Sites-UniqloUS-Site/default/Product-GetAvailability?pid=401925COL69SMA001000&Quantity=1"

type store struct {
	mu sync.Mutex
	M  map[string]map[string]Product
}

func (s *store) Add(user string, p Product) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.M[user]; !ok {
		s.M[user] = make(map[string]Product)
	}
	s.M[user][p.SKU()] = p
}

func (s *store) Remove(user string, p Product) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.M[user]; !ok {
		return
	}
	delete(s.M[user], p.SKU())
}

type Plugin struct {
	store *store
}

func (p *Plugin) Load(r *bot.Robot) {
	r.Respond(bot.Regexp(`track[\s\w]*https?://www.uniqlo.com.*?(\d*).html`), func(r bot.Responder) error {
		id := r.Match[1]
		colors, _ := getColors(id)
		// use colors to present color options and sizes to present size options
		p.store.Add(r.User, Product{ID: id, Color: colors[1].Code, Size: XXS})
		return r.Reply("On it, boss!")
	})

	r.Respond(bot.Regexp("stop tracking uniqlo item"), func(r bot.Responder) error {
		return nil
	})

	r.Router.HandleFunc("interaction", func(w http.ResponseWriter, req *http.Request) {
		token := "eivqImx4C6ZjYIaSvPGCa2M1"
		decoder := json.NewDecoder(req.Body)
		var body slack.AttachmentActionCallback
		if err := decoder.Decode(&body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if body.Token != token {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		adapter, ok := r.Chat.(slacker.Adapter)
		if !ok {
			return
		}

		go func() { adapter.Client.UpdateMessage(body.Channel.ID, body.MessageTs, "Done!") }()

		w.WriteHeader(http.StatusOK)
	}).Methods("POST")

	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		products := []Product{
			{ID: "401925", Color: "COL01", Size: XXS},
			{ID: "401925", Color: "COL09", Size: XXS},
			{ID: "401925", Color: "COL17", Size: XXS},
			{ID: "401925", Color: "COL69", Size: XXS},
		}
		for range ticker.C {
			for _, product := range products {
				if !product.Available() {
					continue
				}
				err := r.Chat.Direct(bot.Message{
					User: "nancy",
					Text: product.URL(),
				})
				if err != nil {
					r.Logger.Error(err.Error())
				}
			}
		}
	}()
}
