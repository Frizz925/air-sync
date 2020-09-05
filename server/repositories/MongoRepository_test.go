package repositories

import (
	"air-sync/models"
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoRepository(t *testing.T) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		t.Fatal("Requires MONGODB_URI env to test")
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	require.Nil(t, err)
	db := client.Database("airsync-test")
	require.Nil(t, db.Drop(ctx))

	opts := MongoOptions{db, true}
	sessionRepo := NewSessionMongoRepository(ctx, opts)
	require.Nil(t, sessionRepo.Migrate())
	attachmentRepo := NewAttachmentMongoRepository(ctx, opts)
	require.Nil(t, attachmentRepo.Migrate())

	attachment, err := attachmentRepo.Create(models.CreateAttachment{})
	require.Nil(t, err)

	session, err := sessionRepo.Create()
	require.Nil(t, err)
	_, err = sessionRepo.InsertMessage(session.ID, models.InsertMessage{})
	require.Nil(t, err)

	session, err = sessionRepo.Create()
	require.Nil(t, err)
	insert := models.InsertMessage{}
	insert.AttachmentID = attachment.ID
	_, err = sessionRepo.InsertMessage(session.ID, insert)
	require.Nil(t, err)
	_, err = sessionRepo.InsertMessage(session.ID, insert)
	require.Nil(t, err)
	_, err = sessionRepo.InsertMessage(session.ID, models.InsertMessage{})
	require.Nil(t, err)

	found, err := sessionRepo.Find(session.ID)
	require.Nil(t, err)
	require.Equal(t, session.ID, found.ID)
	require.Equal(t, 3, len(found.Messages))
	require.GreaterOrEqual(t, found.Messages[0].CreatedAt, found.Messages[1].CreatedAt)
	require.Empty(t, found.Messages[0].AttachmentID)
	require.Equal(t, attachment.ID, found.Messages[1].AttachmentID)
}
