package repository

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/lacolaco/activitypub.lacolaco.net/model"
)

const (
	JobsCollectionName = "jobs"
)

type jobRepo struct {
	firestoreClient *firestore.Client
}

func NewJobRepository(firestoreClient *firestore.Client) *jobRepo {
	return &jobRepo{firestoreClient: firestoreClient}
}

func (r *jobRepo) FindByID(ctx context.Context, jobID string) (*model.Job, error) {
	doc, err := r.firestoreClient.Collection(JobsCollectionName).Doc(jobID).Get(ctx)
	if err != nil {
		return nil, err
	}
	var job *model.Job
	if err := doc.DataTo(&job); err != nil {
		return nil, err
	}
	return job, nil
}

func (r *jobRepo) Add(ctx context.Context, job *model.Job) error {
	col := r.firestoreClient.Collection(JobsCollectionName)
	if err := addIfNotExists(ctx, col, job); err != nil {
		return err
	}
	return nil
}

func (r *jobRepo) DeleteByID(ctx context.Context, jobID string) error {
	collection := r.firestoreClient.Collection(JobsCollectionName)
	if err := removeIfExists(ctx, collection, jobID); err != nil {
		return err
	}
	return nil
}
