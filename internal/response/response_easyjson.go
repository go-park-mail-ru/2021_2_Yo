// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package response

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6ff3ac1dDecodeBackendInternalResponse(in *jlexer.Lexer, out *UserResponseBody) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "name":
			out.Name = string(in.String())
		case "surname":
			out.Surname = string(in.String())
		case "description":
			out.About = string(in.String())
		case "imgUrl":
			out.ImgUrl = string(in.String())
		case "email":
			out.Mail = string(in.String())
		case "password":
			out.Password = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncodeBackendInternalResponse(out *jwriter.Writer, in UserResponseBody) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	if in.Name != "" {
		const prefix string = ",\"name\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Name))
	}
	if in.Surname != "" {
		const prefix string = ",\"surname\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Surname))
	}
	if in.About != "" {
		const prefix string = ",\"description\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.About))
	}
	if in.ImgUrl != "" {
		const prefix string = ",\"imgUrl\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.ImgUrl))
	}
	if in.Mail != "" {
		const prefix string = ",\"email\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Mail))
	}
	if in.Password != "" {
		const prefix string = ",\"password\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Password))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserResponseBody) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeBackendInternalResponse(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserResponseBody) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeBackendInternalResponse(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserResponseBody) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeBackendInternalResponse(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserResponseBody) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeBackendInternalResponse(l, v)
}
func easyjson6ff3ac1dDecodeBackendInternalResponse1(in *jlexer.Lexer, out *UserListResponseBody) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "users":
			if in.IsNull() {
				in.Skip()
				out.Users = nil
			} else {
				in.Delim('[')
				if out.Users == nil {
					if !in.IsDelim(']') {
						out.Users = make([]UserResponseBody, 0, 0)
					} else {
						out.Users = []UserResponseBody{}
					}
				} else {
					out.Users = (out.Users)[:0]
				}
				for !in.IsDelim(']') {
					var v1 UserResponseBody
					(v1).UnmarshalEasyJSON(in)
					out.Users = append(out.Users, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncodeBackendInternalResponse1(out *jwriter.Writer, in UserListResponseBody) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"users\":"
		out.RawString(prefix[1:])
		if in.Users == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.Users {
				if v2 > 0 {
					out.RawByte(',')
				}
				(v3).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v UserListResponseBody) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeBackendInternalResponse1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v UserListResponseBody) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeBackendInternalResponse1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *UserListResponseBody) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeBackendInternalResponse1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *UserListResponseBody) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeBackendInternalResponse1(l, v)
}
func easyjson6ff3ac1dDecodeBackendInternalResponse2(in *jlexer.Lexer, out *SubscribedResponseBody) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "result":
			out.Result = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncodeBackendInternalResponse2(out *jwriter.Writer, in SubscribedResponseBody) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"result\":"
		out.RawString(prefix[1:])
		out.Bool(bool(in.Result))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SubscribedResponseBody) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeBackendInternalResponse2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SubscribedResponseBody) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeBackendInternalResponse2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SubscribedResponseBody) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeBackendInternalResponse2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SubscribedResponseBody) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeBackendInternalResponse2(l, v)
}
func easyjson6ff3ac1dDecodeBackendInternalResponse3(in *jlexer.Lexer, out *Response) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "status":
			out.Status = HttpStatus(in.Int())
		case "message":
			out.Message = string(in.String())
		case "body":
			if m, ok := out.Body.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Body.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Body = in.Interface()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncodeBackendInternalResponse3(out *jwriter.Writer, in Response) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"status\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Status))
	}
	if in.Message != "" {
		const prefix string = ",\"message\":"
		out.RawString(prefix)
		out.String(string(in.Message))
	}
	if in.Body != nil {
		const prefix string = ",\"body\":"
		out.RawString(prefix)
		if m, ok := in.Body.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Body.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Body))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Response) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeBackendInternalResponse3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Response) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeBackendInternalResponse3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Response) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeBackendInternalResponse3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Response) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeBackendInternalResponse3(l, v)
}
func easyjson6ff3ac1dDecodeBackendInternalResponse4(in *jlexer.Lexer, out *NotificationResponseBody) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "type":
			out.Type = string(in.String())
		case "userId":
			out.UserId = string(in.String())
		case "userName":
			out.UserName = string(in.String())
		case "userSurname":
			out.UserSurname = string(in.String())
		case "userImgUrl":
			out.UserImgUrl = string(in.String())
		case "eventId":
			out.EventId = string(in.String())
		case "eventTitle":
			out.EventTitle = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncodeBackendInternalResponse4(out *jwriter.Writer, in NotificationResponseBody) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix[1:])
		out.String(string(in.Type))
	}
	{
		const prefix string = ",\"userId\":"
		out.RawString(prefix)
		out.String(string(in.UserId))
	}
	{
		const prefix string = ",\"userName\":"
		out.RawString(prefix)
		out.String(string(in.UserName))
	}
	{
		const prefix string = ",\"userSurname\":"
		out.RawString(prefix)
		out.String(string(in.UserSurname))
	}
	{
		const prefix string = ",\"userImgUrl\":"
		out.RawString(prefix)
		out.String(string(in.UserImgUrl))
	}
	if in.EventId != "" {
		const prefix string = ",\"eventId\":"
		out.RawString(prefix)
		out.String(string(in.EventId))
	}
	if in.EventTitle != "" {
		const prefix string = ",\"eventTitle\":"
		out.RawString(prefix)
		out.String(string(in.EventTitle))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NotificationResponseBody) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeBackendInternalResponse4(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NotificationResponseBody) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeBackendInternalResponse4(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NotificationResponseBody) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeBackendInternalResponse4(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NotificationResponseBody) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeBackendInternalResponse4(l, v)
}
func easyjson6ff3ac1dDecodeBackendInternalResponse5(in *jlexer.Lexer, out *NotificationListResponseBody) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "notifications":
			if in.IsNull() {
				in.Skip()
				out.Notifications = nil
			} else {
				in.Delim('[')
				if out.Notifications == nil {
					if !in.IsDelim(']') {
						out.Notifications = make([]NotificationResponseBody, 0, 0)
					} else {
						out.Notifications = []NotificationResponseBody{}
					}
				} else {
					out.Notifications = (out.Notifications)[:0]
				}
				for !in.IsDelim(']') {
					var v4 NotificationResponseBody
					(v4).UnmarshalEasyJSON(in)
					out.Notifications = append(out.Notifications, v4)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncodeBackendInternalResponse5(out *jwriter.Writer, in NotificationListResponseBody) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"notifications\":"
		out.RawString(prefix[1:])
		if in.Notifications == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v5, v6 := range in.Notifications {
				if v5 > 0 {
					out.RawByte(',')
				}
				(v6).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v NotificationListResponseBody) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeBackendInternalResponse5(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v NotificationListResponseBody) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeBackendInternalResponse5(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *NotificationListResponseBody) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeBackendInternalResponse5(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *NotificationListResponseBody) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeBackendInternalResponse5(l, v)
}
func easyjson6ff3ac1dDecodeBackendInternalResponse6(in *jlexer.Lexer, out *FavouriteResponseBody) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "result":
			out.Result = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncodeBackendInternalResponse6(out *jwriter.Writer, in FavouriteResponseBody) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"result\":"
		out.RawString(prefix[1:])
		out.Bool(bool(in.Result))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v FavouriteResponseBody) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeBackendInternalResponse6(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v FavouriteResponseBody) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeBackendInternalResponse6(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *FavouriteResponseBody) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeBackendInternalResponse6(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *FavouriteResponseBody) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeBackendInternalResponse6(l, v)
}
func easyjson6ff3ac1dDecodeBackendInternalResponse7(in *jlexer.Lexer, out *EventResponseBody) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "title":
			out.Title = string(in.String())
		case "description":
			out.Description = string(in.String())
		case "text":
			out.Text = string(in.String())
		case "city":
			out.City = string(in.String())
		case "category":
			out.Category = string(in.String())
		case "viewed":
			out.Viewed = int(in.Int())
		case "imgUrl":
			out.ImgUrl = string(in.String())
		case "tag":
			if in.IsNull() {
				in.Skip()
				out.Tag = nil
			} else {
				in.Delim('[')
				if out.Tag == nil {
					if !in.IsDelim(']') {
						out.Tag = make([]string, 0, 4)
					} else {
						out.Tag = []string{}
					}
				} else {
					out.Tag = (out.Tag)[:0]
				}
				for !in.IsDelim(']') {
					var v7 string
					v7 = string(in.String())
					out.Tag = append(out.Tag, v7)
					in.WantComma()
				}
				in.Delim(']')
			}
		case "date":
			out.Date = string(in.String())
		case "geo":
			out.Geo = string(in.String())
		case "address":
			out.Address = string(in.String())
		case "authorid":
			out.AuthorID = string(in.String())
		case "favourite":
			out.IsVisited = bool(in.Bool())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncodeBackendInternalResponse7(out *jwriter.Writer, in EventResponseBody) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ID != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"title\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Title))
	}
	{
		const prefix string = ",\"description\":"
		out.RawString(prefix)
		out.String(string(in.Description))
	}
	{
		const prefix string = ",\"text\":"
		out.RawString(prefix)
		out.String(string(in.Text))
	}
	{
		const prefix string = ",\"city\":"
		out.RawString(prefix)
		out.String(string(in.City))
	}
	{
		const prefix string = ",\"category\":"
		out.RawString(prefix)
		out.String(string(in.Category))
	}
	{
		const prefix string = ",\"viewed\":"
		out.RawString(prefix)
		out.Int(int(in.Viewed))
	}
	{
		const prefix string = ",\"imgUrl\":"
		out.RawString(prefix)
		out.String(string(in.ImgUrl))
	}
	{
		const prefix string = ",\"tag\":"
		out.RawString(prefix)
		if in.Tag == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v8, v9 := range in.Tag {
				if v8 > 0 {
					out.RawByte(',')
				}
				out.String(string(v9))
			}
			out.RawByte(']')
		}
	}
	{
		const prefix string = ",\"date\":"
		out.RawString(prefix)
		out.String(string(in.Date))
	}
	{
		const prefix string = ",\"geo\":"
		out.RawString(prefix)
		out.String(string(in.Geo))
	}
	{
		const prefix string = ",\"address\":"
		out.RawString(prefix)
		out.String(string(in.Address))
	}
	{
		const prefix string = ",\"authorid\":"
		out.RawString(prefix)
		out.String(string(in.AuthorID))
	}
	{
		const prefix string = ",\"favourite\":"
		out.RawString(prefix)
		out.Bool(bool(in.IsVisited))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EventResponseBody) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeBackendInternalResponse7(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventResponseBody) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeBackendInternalResponse7(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventResponseBody) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeBackendInternalResponse7(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventResponseBody) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeBackendInternalResponse7(l, v)
}
func easyjson6ff3ac1dDecodeBackendInternalResponse8(in *jlexer.Lexer, out *EventListResponseBody) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "events":
			if in.IsNull() {
				in.Skip()
				out.Events = nil
			} else {
				in.Delim('[')
				if out.Events == nil {
					if !in.IsDelim(']') {
						out.Events = make([]EventResponseBody, 0, 0)
					} else {
						out.Events = []EventResponseBody{}
					}
				} else {
					out.Events = (out.Events)[:0]
				}
				for !in.IsDelim(']') {
					var v10 EventResponseBody
					(v10).UnmarshalEasyJSON(in)
					out.Events = append(out.Events, v10)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncodeBackendInternalResponse8(out *jwriter.Writer, in EventListResponseBody) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"events\":"
		out.RawString(prefix[1:])
		if in.Events == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v11, v12 := range in.Events {
				if v11 > 0 {
					out.RawByte(',')
				}
				(v12).MarshalEasyJSON(out)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EventListResponseBody) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeBackendInternalResponse8(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventListResponseBody) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeBackendInternalResponse8(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventListResponseBody) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeBackendInternalResponse8(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventListResponseBody) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeBackendInternalResponse8(l, v)
}
func easyjson6ff3ac1dDecodeBackendInternalResponse9(in *jlexer.Lexer, out *EventIDResponseBody) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncodeBackendInternalResponse9(out *jwriter.Writer, in EventIDResponseBody) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v EventIDResponseBody) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeBackendInternalResponse9(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v EventIDResponseBody) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeBackendInternalResponse9(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *EventIDResponseBody) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeBackendInternalResponse9(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *EventIDResponseBody) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeBackendInternalResponse9(l, v)
}
func easyjson6ff3ac1dDecodeBackendInternalResponse10(in *jlexer.Lexer, out *CitiesResponseBody) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "cities":
			if in.IsNull() {
				in.Skip()
				out.Cities = nil
			} else {
				in.Delim('[')
				if out.Cities == nil {
					if !in.IsDelim(']') {
						out.Cities = make([]string, 0, 4)
					} else {
						out.Cities = []string{}
					}
				} else {
					out.Cities = (out.Cities)[:0]
				}
				for !in.IsDelim(']') {
					var v13 string
					v13 = string(in.String())
					out.Cities = append(out.Cities, v13)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6ff3ac1dEncodeBackendInternalResponse10(out *jwriter.Writer, in CitiesResponseBody) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"cities\":"
		out.RawString(prefix[1:])
		if in.Cities == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v14, v15 := range in.Cities {
				if v14 > 0 {
					out.RawByte(',')
				}
				out.String(string(v15))
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CitiesResponseBody) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6ff3ac1dEncodeBackendInternalResponse10(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CitiesResponseBody) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6ff3ac1dEncodeBackendInternalResponse10(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CitiesResponseBody) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6ff3ac1dDecodeBackendInternalResponse10(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CitiesResponseBody) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6ff3ac1dDecodeBackendInternalResponse10(l, v)
}