package localstorage

import (
	"backend/event"
	"backend/models"
	"strconv"
	"sync"
)

type Repository struct {
	events []*Event
	mutex  *sync.Mutex
}

var eventsDemo = []*Event{
	&Event{
		0,
		"Jusa Tusa",
		"дискотека это тусовка или просто сборище? 8 лет. Дискотека - это когда есть диджей и в этом деле разбираются все и молодежь и взрослые.",
		10,
		"/img/tusa.jpeg",
	},
	&Event{
		1,
		"Джуса Туса",
		"Тусовка весёлая) Реально весело)",
		24,
		"/img/tusa.jpeg",
	},
	&Event{
		2,
		"Тжуса Дуса",
		"Йоу, идешь тусить? Нет? А почему? Пойдём! Туса-Джуса будет!",
		822,
		"/img/tusa.jpeg",
	},
	&Event{
		3,
		"Дуса Тжуса",
		"дискотека это тусовка или просто сборище? 8 лет. Дискотека - это когда есть диджей и в этом деле разбираются все и молодежь и взрослые.",
		3901,
		"/img/tusa.jpeg",
	},
	&Event{
		4,
		"Тужса Туса",
		"дискотека это тусовка или просто сборище? 8 лет. Дискотека - это когда есть диджей и в этом деле разбираются все и молодежь и взрослые.",
		192,
		"/img/tusa.jpeg",
	},
	&Event{
		5,
		"Дуда Туда",
		"дискотека это тусовка или просто сборище? 8 лет. Дискотека - это когда есть диджей и в этом деле разбираются все и молодежь и взрослые.",
		108,
		"/img/tusa.jpeg",
	},
	&Event{
		6,
		"ЙАУУХУ",
		"дискотека это тусовка или просто сборище? 8 лет. Дискотека - это когда есть диджей и в этом деле разбираются все и молодежь и взрослые.",
		1,
		"/img/tusa.jpeg",
	},
	&Event{
		7,
		"Тусняяяк",
		"дискотека это тусовка или просто сборище? 8 лет. Дискотека - это когда есть диджей и в этом деле разбираются все и молодежь и взрослые.",
		1042,
		"/img/tusa.jpeg",
	},
	&Event{
		8,
		"Бауманский посвят",
		"Вот уже несколько десятилетий подряд из года в год соблюдается Ежегодно 30-31 сентрября все. Ясно?",
		8989,
		"/img/tusa.jpeg",
	},
}

func NewRepository() *Repository {
	result := &Repository{
		//events: make([]*Event, 0),
		events: eventsDemo,
		mutex:  new(sync.Mutex),
	}
	return result
}

func (s *Repository) List() ([]*models.Event, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	resultEvents := make([]*models.Event, len(s.events))
	for i := 0; i < len(s.events); i++ {
		resultEvents[i] = toModelEvent(s.events[i])
	}
	return resultEvents, nil
}

func (s *Repository) GetEvent(eventId string) (*models.Event, error) {
	eventIdInt, _ := strconv.Atoi(eventId)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, foundEvent := range s.events {
		if foundEvent.ID == eventIdInt {
			return toModelEvent(foundEvent), nil
		}
	}
	return nil, event.ErrEventNotFound
}
