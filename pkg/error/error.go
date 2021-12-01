package error

import "errors"

var (
	ErrUserNotFound    = errors.New("Пользователь не найден")
	ErrEmptyData       = errors.New("Отсутствуют необходимые данные")
	ErrPostgres        = errors.New("Проблема с базой данных")
	ErrAtoi            = errors.New("Введённая строка должна быть числовой")
	ErrNotAllowed      = errors.New("Нет прав на совершение действия")
	ErrNoRows          = errors.New("Запрашиваемые данные отсутствуют")
	ErrCookie          = errors.New("Ошибка с получением cookie")
	ErrUserExists      = errors.New("Пользователь уже зарегистрирован")
	ErrAuthService     = errors.New("Проблема на сервисе авторизации")
	ErrInternal        = errors.New("Проблема на стороне сервера")
	ErrSessionNotFound = errors.New("Сессия пользователя не найдена")
	//ErrNotAuthorised = errors.New("Пользователь не авторизован")
)
