package localstorage

import (
	"backend/models"
	"sync"
)

type RepositoryEventLocalStorage struct {
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
		"дискотека это тусовка или просто сборище? 8 лет. Дискотека - это когда есть диджей и в этом деле разбираются все и молодежь и взрослые.",
		24,
		"/img/tusa.jpeg",
	},
	&Event{
		2,
		"Тжуса Дуса",
		"дискотека это тусовка или просто сборище? 8 лет. Дискотека - это когда есть диджей и в этом деле разбираются все и молодежь и взрослые.",
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
		"дискотека это тусовка или просто сборище? 8 лет. Дискотека - это когда есть диджей и в этом деле разбираются все и молодежь и взрослые.",
		8989,
		"/img/tusa.jpeg",
	},
}

func NewRepositoryEventLocalStorage() *RepositoryEventLocalStorage {
	result := &RepositoryEventLocalStorage{
		//events: make([]*Event, 0),
		events: eventsDemo,
		mutex:  new(sync.Mutex),
	}
	return result
}

func (s *RepositoryEventLocalStorage) List() ([]*models.Event, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	resultEvents := make([]*models.Event, len(s.events))
	for i := 0; i < len(s.events); i++ {
		resultEvents[i] = toModelEvent(s.events[i])
	}
	return resultEvents, nil
}
