package response

import (
	error2 "backend/internal/error"
	models "backend/internal/models"
	log "backend/pkg/logger"
	"errors"
	json "github.com/mailru/easyjson"
	"io"
	"net/http"
	"strings"

	"github.com/asaskevich/govalidator"
	"github.com/go-sanitize/sanitize"
)

var (
	ErrJSONDecoding   = errors.New("data decoding error")
	ErrValidation     = errors.New("data validation error")
	ErrSanitizing     = errors.New("data sanitizing error")
	ErrSanitizerError = errors.New("internal sanitizing package error")
)

func ValidateAndSanitize(object interface{}) error {
	s, err := sanitize.New()
	if err != nil {
		return ErrSanitizerError
	}
	err = s.Sanitize(object)
	if err != nil {
		return ErrSanitizing
	}
	valid, err := govalidator.ValidateStruct(object)
	if err != nil || !valid {
		return ErrValidation
	}
	return nil
}

func GetUserFromRequest(r io.Reader) (*models.User, error) {
	message := logMessage + "GetUserFromRequest:"
	_ = message
	userInput := new(UserResponseBody)
	err := json.UnmarshalFromReader(r, userInput)
	//err := json.NewDecoder(r).Decode(userInput)
	if err != nil {
		return nil, ErrJSONDecoding
	}
	err = ValidateAndSanitize(userInput)
	if err != nil {
		return nil, err
	}
	result := &models.User{
		Name:     userInput.Name,
		Surname:  userInput.Surname,
		Mail:     userInput.Mail,
		Password: userInput.Password,
		About:    userInput.About,
	}
	return result, nil
}

func GetUsersIdFromRequest(r io.Reader) ([]string, error) {
	message := logMessage + "GetUsersIdFromRequest:"
	_ = message
	usersIdInput := new(UsersIdResponseBody)
	err := json.UnmarshalFromReader(r, usersIdInput)
	if err != nil {
		return nil, ErrJSONDecoding
	}
	err = ValidateAndSanitize(usersIdInput)
	if err != nil {
		return nil, err
	}
	result := usersIdInput.UsersId
	return result, nil

}

func MakeUserResponseBody(u *models.User) UserResponseBody {
	return UserResponseBody{
		ID:       u.ID,
		Name:     u.Name,
		Surname:  u.Surname,
		About:    u.About,
		ImgUrl:   u.ImgUrl,
		Mail:     u.Mail,
		Password: u.Password,
	}
}

func MakeUserListResponseBody(users []*models.User) UserListResponseBody {
	result := make([]UserResponseBody, len(users))
	for i := 0; i < len(users); i++ {
		result[i] = MakeUserResponseBody(users[i])
	}
	return UserListResponseBody{
		Users: result,
	}
}

func GetEventFromRequest(r io.Reader) (*models.Event, error) {
	eventInput := new(EventResponseBody)
	err := json.UnmarshalFromReader(r, eventInput)

	if err != nil {
		return nil, ErrJSONDecoding
	}
	err = ValidateAndSanitize(eventInput)
	if err != nil {
		return nil, err
	}
	result := &models.Event{
		ID:          eventInput.ID,
		Title:       eventInput.Title,
		Description: eventInput.Description,
		Text:        eventInput.Text,
		City:        eventInput.City,
		Category:    eventInput.Category,
		Viewed:      eventInput.Viewed,
		ImgUrl:      eventInput.ImgUrl,
		Tag:         eventInput.Tag,
		Date:        eventInput.Date,
		Geo:         eventInput.Geo,
		Address:     eventInput.Address,
	}
	return result, nil
}

func MakeEventResponseBody(e *models.Event) EventResponseBody {
	return EventResponseBody{
		ID:          e.ID,
		Title:       e.Title,
		Description: e.Description,
		Text:        e.Text,
		City:        e.City,
		Category:    e.Category,
		Viewed:      e.Viewed,
		ImgUrl:      e.ImgUrl,
		Tag:         e.Tag,
		Date:        e.Date,
		Geo:         e.Geo,
		Address:     e.Address,
		AuthorID:    e.AuthorId,
		IsVisited:   e.IsVisited,
	}
}

func MakeEventListResponseBody(events []*models.Event) EventListResponseBody {
	result := make([]EventResponseBody, len(events))
	for i := 0; i < len(events); i++ {
		result[i] = MakeEventResponseBody(events[i])
	}
	return EventListResponseBody{
		Events: result,
	}
}

func MakeNotificationResponseBody(n *models.Notification) NotificationResponseBody {
	return NotificationResponseBody{
		Type:        n.Type,
		Seen:        n.Seen,
		UserId:      n.UserId,
		UserName:    n.UserName,
		UserSurname: n.UserSurname,
		UserImgUrl:  n.UserImgUrl,
		EventId:     n.EventId,
		EventTitle:  n.EventTitle,
	}
}

func MakeNotificationListResponseBody(notifications []*models.Notification) NotificationListResponseBody {
	result := make([]NotificationResponseBody, len(notifications))
	for i := 0; i < len(notifications); i++ {
		result[i] = MakeNotificationResponseBody(notifications[i])
	}
	return NotificationListResponseBody{
		Notifications: result,
	}
}

func SendResponse(w http.ResponseWriter, response *Response) {
	message := logMessage + "SendResponse:"
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(response)
	if err != nil {
		log.Error(message+"err =", err)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Error(message+"err =", err)
	}
}

func refactorError(err error) (error, HttpStatus) {
	if err == nil {
		return nil, http.StatusOK
	}
	errStr := err.Error()
	if strings.Contains(errStr, "user already exists") {
		return error2.ErrUserExists, http.StatusConflict
	}
	if strings.Contains(errStr, "user not found") {
		return error2.ErrUserNotFound, http.StatusNotFound
	}
	if strings.Contains(errStr, "internal DB server error") {
		return error2.ErrPostgres, http.StatusInternalServerError
	}
	if strings.Contains(errStr, "cookie") {
		return error2.ErrCookie, http.StatusUnauthorized
	}
	if strings.Contains(errStr, "required data is empty") {
		return error2.ErrEmptyData, http.StatusBadRequest
	}
	if strings.Contains(errStr, "cant cast string to int") {
		return error2.ErrAtoi, http.StatusBadRequest
	}
	if strings.Contains(errStr, "user is not allowed to do this") {
		return error2.ErrNotAllowed, http.StatusForbidden
	}
	if strings.Contains(errStr, "no rows in a query result") {
		return error2.ErrNoRows, http.StatusNotFound
	}
	if strings.Contains(errStr, "Error while dialing dial tcp") {
		return error2.ErrInternal, http.StatusInternalServerError
	}
	if strings.Contains(errStr, "session was not found") {
		return error2.ErrSessionNotFound, http.StatusUnauthorized
	}
	return err, http.StatusBadRequest
}

func CheckIfNoError(w *http.ResponseWriter, err error, msg string) bool {
	if err != nil {
		log.Error(msg+"err = ", err)
	}
	errRefactored, status := refactorError(err)
	if err != nil {
		log.Error(msg+"refactored err = ", errRefactored)
		SendResponse(*w, StatusResponse(status))
		return false
	}
	return true
}
