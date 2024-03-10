package models

import (
	"time"

	"github.com/google/uuid"
)

type Friend struct {
	ID        uint `gorm:"primaryKey"`
	UserID    uint // ID do usuário que está realizando a ação (adicionar/remover amigo)
	FriendID  uuid.UUID
	CreatedAt time.Time // Data e hora da criação da relação de amizade
	UpdatedAt time.Time // Data e hora da última atualização
}
