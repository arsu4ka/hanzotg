package repo

import (
	"database/sql"
	"hanzotg/internal/app/models"
)

type VerificationMessageRepo struct {
	db *sql.DB
}

func NewVerificationMessage(db *sql.DB) *VerificationMessageRepo {
	return &VerificationMessageRepo{db: db}
}

func (m *VerificationMessageRepo) Create(chatId, messageId, aboutUserId string) error {
	_, err := m.db.Exec(
		"INSERT INTO verification_messages (message_id, chat_id, about_user_id, status) VALUES ($1, $2, $3, $4)",
		messageId,
		chatId,
		aboutUserId,
		models.VerificationMessageStatusPending,
	)
	return err
}

func (m *VerificationMessageRepo) UpdateStatuses(aboutUserId string, newStatus models.VerificationMessageStatus) error {
	_, err := m.db.Exec(
		"UPDATE verification_messages SET status = $1 WHERE about_user_id = $2 AND status != $3",
		newStatus,
		aboutUserId,
		newStatus,
	)
	return err
}

func (m *VerificationMessageRepo) GetMessagesAboutUserId(aboutUserId string) ([]*models.VerificationMessage, error) {
	var verificationMessages []*models.VerificationMessage

	// Prepare the SQL query
	query := `
        SELECT message_id, chat_id, about_user_id, status
        FROM verification_messages
        WHERE about_user_id = $1
    `

	// Execute the query
	rows, err := m.db.Query(query, aboutUserId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the result set and populate the verificationMessages slice
	for rows.Next() {
		var msg models.VerificationMessage
		if err := rows.Scan(&msg.MessageId, &msg.ChatId, &msg.AboutUserId, &msg.Status); err != nil {
			return nil, err
		}
		verificationMessages = append(verificationMessages, &msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return verificationMessages, nil
}

func (m *VerificationMessageRepo) IsUserWaitingForPaymentAccept(userId string) (bool, error) {
	query := `
		SELECT COUNT(*) as message_number
		FROM verification_messages
		WHERE status = $1 AND about_user_id = $2
	`

	var pendingMessagesCount int
	err := m.db.QueryRow(query, models.VerificationMessageStatusPending, userId).Scan(&pendingMessagesCount)
	return pendingMessagesCount > 0, err
}
