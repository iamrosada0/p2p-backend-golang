package models

import (
	"time"
)

type Post struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      // ID do usuário que criou o post
	Content   string    // Conteúdo do post
	ImageURL  string    // URL da imagem associada ao post (opcional)
	VideoURL  string    // URL do vídeo associado ao post (opcional)
	CreatedAt time.Time // Data e hora de criação do post
	UpdatedAt time.Time // Data e hora de atualização do post
}
