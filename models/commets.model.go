package models

import "time"

type Comment struct {
	ID              uint   `gorm:"primaryKey"`
	PostID          uint   // ID da postagem à qual o comentário está associado
	UserID          uint   // ID do usuário que fez o comentário
	Content         string // Conteúdo do comentário
	Likes           int    // Número de curtidas (apoio) para o comentário
	Dislikes        int    // Número de descurtidas (reprovação) para o comentário
	MentionedPostID uint   // ID da publicação mencionada (opcional)
	ReplyToID       uint   // ID do comentário ao qual este comentário está respondendo (opcional)

	ParentCommentID uint      // ID do comentário pai (para respostas a comentários)
	CreatedAt       time.Time // Data e hora de criação do comentário
	UpdatedAt       time.Time // Data e hora de atualização do comentário
}
