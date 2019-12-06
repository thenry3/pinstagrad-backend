package database

import (
	"context"
	"log"

	"firebase.google.com/go/db"
)

func filterByPhotographerName(ctx context.Context, reference *db.Ref, photographer string) map[string]Photo {
	v := map[string]Photo{}
	if err := reference.Child("all").OrderByChild("photographer").StartAt(photographer).EndAt(photographer).Get(ctx, &v); err != nil {
		log.Fatalln("Error reading value:", err)
	}
	return v
}

func filterByTags(ctx context.Context, reference *db.Ref, tagType ...string) map[string]Photo {
	v := map[string]Photo{}
	reference.Child("all")
	/*for _, tag := range tagType {
		if err := reference.Child("all").OrderByChild("photographer").StartAt(photographer).EndAt(photographer).Get(ctx, &v); err != nil {
			log.Fatalln("Error reading value:", err)
		}
	}*/
	return v
}
