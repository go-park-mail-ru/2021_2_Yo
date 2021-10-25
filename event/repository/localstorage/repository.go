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
		ID:    1,
		Title: "Jusa Tusa",
		Description: "Небольшое описание мероприятия. " +
			"Да, реально крутая тусовка. Да, говорю. Круто будет, говорю, весело." +
			"Всем ясно? Тусовка. Тусовка. Тусовка. Тусовка. Тусовка. Тусовка. Этот прямоугольник должен" +
			"сжиматься/расширяться в зависимости от длины текста (количества строк).",
		Text: "But I must explain to you how all this mistaken idea of denouncing" +
			"pleasure and praising pain was born and I will give you a complete account of the system," +
			"and expound the actual " +
			"teachings of the great explorer of the truth, the master-builder of human happiness." +
			"No one rejects, dislikes," +
			"or avoids pleasure itself, because it is pleasure, but because those who do not know " +
			"'how to pursue pleasure" +
			"rationally encounter consequences that are extremely painful. Nor again is '" +
			"there anyone who loves or pursues or" +
			"desires to obtain pain of itself, because it is pain, but because occasionally " +
			"'circumstances occur in which" +
			"toil and pain can procure him some great pleasure.',",
		City:     "Москва",
		Category: "Тусовка",
		Viewed:   1220,
		Tag:      []string{"nil", "alco", "hey"},
		Date:     "20.01.19",
		Geo:      "Izmaiilfoofo",
	},
	&Event{
		ID:    2,
		Title: "FFFFFFFFFFFF",
		Description: "Небольшое описание мероприятия. " +
			"Да, реально крутая тусовка. Да, говорю. Круто будет, говорю, весело." +
			"Всем ясно? Тусовка. Тусовка. Тусовка. Тусовка. Тусовка. Тусовка. Этот прямоугольник должен" +
			"сжиматься/расширяться в зависимости от длины текста (количества строк).",
		Text: "But I must explain to you how all this mistaken idea of denouncing" +
			"pleasure and praising pain was born and I will give you a complete account of the system," +
			"and expound the actual " +
			"teachings of the great explorer of the truth, the master-builder of human happiness." +
			"No one rejects, dislikes," +
			"or avoids pleasure itself, because it is pleasure, but because those who do not know " +
			"'how to pursue pleasure" +
			"rationally encounter consequences that are extremely painful. Nor again is '" +
			"there anyone who loves or pursues or" +
			"desires to obtain pain of itself, because it is pain, but because occasionally " +
			"'circumstances occur in which" +
			"toil and pain can procure him some great pleasure.',",
		City:     "Москва",
		Category: "Тусовка",
		Viewed:   1220,
		Tag:      []string{"nil", "alco", "hey"},
		Date:     "20.01.19",
		Geo:      "Izmaiilfoofo",
	},
	&Event{
		ID:    3,
		Title: "ewfhwekfhhlwek",
		Description: "Небольшое описание мероприятия. " +
			"Да, реально крутая тусовка. Да, говорю. Круто будет, говорю, весело." +
			"Всем ясно? Тусовка. Тусовка. Тусовка. Тусовка. Тусовка. Тусовка. Этот прямоугольник должен" +
			"сжиматься/расширяться в зависимости от длины текста (количества строк).",
		Text: "But I must explain to you how all this mistaken idea of denouncing" +
			"pleasure and praising pain was born and I will give you a complete account of the system,",
		City:     "Москва",
		Category: "Тусовка",
		Viewed:   1220,
		Tag:      []string{"nil", "alco", "hey"},
		Date:     "20.01.19",
		Geo:      "Izmaiilfoofo",
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
