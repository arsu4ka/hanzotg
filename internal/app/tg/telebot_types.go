package tg

type userRecipient struct {
	UserId string
}

func (r *userRecipient) Recipient() string {
	return r.UserId
}
