package repository

import (
	"context"
	"fmt"
	"learnfirestore/models"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/google/uuid"
)

type GroupRepository struct {
	fs       *firestore.Client
	coll     *firestore.CollectionRef
	collName string
}

func NewGroupRepository(fs *firestore.Client) *GroupRepository {
	return &GroupRepository{
		fs:       fs,
		coll:     fs.Collection("groups"),
		collName: "groups",
	}
}

func (r *GroupRepository) WithCollection(collectionName string) *GroupRepository {
	return &GroupRepository{
		fs:       r.fs,
		coll:     r.fs.Collection(collectionName),
		collName: collectionName,
	}
}

func (r GroupRepository) GetGroupWithID(ctx context.Context, id string) (models.Group, error) {
	snap, err := r.coll.Doc(id).Get(ctx)
	if err != nil {
		return models.Group{}, fmt.Errorf("failed to get group: %w", err)
	}

	var group models.Group
	if err := snap.DataTo(&group); err != nil {
		return models.Group{}, fmt.Errorf("failed to return group: %w", err)
	}

	return group, nil
}

// IF REQUIRED CAN CHANGE CHANNEL TYPE
// type GroupChanges struct {
// 	Before models.Group
// 	After  models.Group
// }

func (r GroupRepository) WatchGroups(ctx context.Context,
	documentAdded chan models.Group,
	documentModified chan models.GroupChanges,
	documentRemoved chan models.Group,
) error {
	snap := r.coll.Snapshots(ctx)
	tempGroup := make(map[string]models.Group)

	fmt.Println("Start watching Groups Collection...")

	for {
		qs, err := snap.Next()
		if err != nil {
			panic(err)
		}

		for _, v := range qs.Changes {
			// REMEMBER: YOU MUST CONVERT THIS TO STRUCTURE DO NOT USE .DATA
			var currentGroup models.Group
			if err := v.Doc.DataTo(&currentGroup); err != nil {
				panic(err)
			}
			currentGroup.ID = v.Doc.Ref.ID

			switch v.Kind {
			case firestore.DocumentAdded:
				documentAdded <- currentGroup
				tempGroup[currentGroup.ID] = currentGroup
			case firestore.DocumentModified:
				previousGroup, exists := tempGroup[currentGroup.ID]
				if !exists {
					continue
				} else {
					GroupChanges := models.GroupChanges{
						Before: previousGroup,
						After:  currentGroup,
					}
					documentModified <- GroupChanges
				}
				tempGroup[currentGroup.ID] = currentGroup
			case firestore.DocumentRemoved:
				documentRemoved <- currentGroup
				delete(tempGroup, currentGroup.ID)
			}
		}
	}
}

func (r GroupRepository) CreateGroup(ctx context.Context, name string, description string, member models.Member) (models.Group, error) {
	groupRef := r.coll.NewDoc()
	group := models.Group{
		ID:          groupRef.ID,
		Name:        name,
		Description: description,
		Members:     map[string]models.Member{},
		MembersID:   []string{member.ID},
		Tasks:       []models.Tasks{},
		Comments:    []models.Comments{},
		CreatedAt:   time.Now(),
	}
	group.Members[member.ID] = models.Member{
		ID:   member.ID,
		Name: member.Name,
	}

	if _, err := groupRef.Set(ctx, group); err != nil {
		return models.Group{}, fmt.Errorf("failed to create Group: %w", err)
	}
	return group, nil
}

func (r GroupRepository) UpdateGroup(ctx context.Context, id string, group *models.Group) error {
	if id == "" {
		return fmt.Errorf("id group cannot be empty")
	}
	if group == nil {
		return fmt.Errorf("group not found")
	}

	if _, err := r.coll.Doc(id).Set(ctx, group, firestore.MergeAll); err != nil {
		return fmt.Errorf("failed to update Group: %w", err)
	}
	return nil
}

func (r GroupRepository) DeleteGroup(ctx context.Context, id string) error {

	if id == "" {
		return fmt.Errorf("group cannot be empty")
	}

	updates := []firestore.Update{
		{Path: "deletedAt", Value: time.Now()},
	}

	_, err := r.coll.Doc(id).Update(ctx, updates)
	if err != nil {
		return fmt.Errorf("failed to delete group: %w", err)
	}
	return nil
}

func (r GroupRepository) CreateTask(ctx context.Context, groupID string, title string, description string, priority bool, duedate time.Time, createdBy models.Member) (models.Tasks, error) {
	task := models.Tasks{
		UUID:        uuid.NewString(),
		Title:       title,
		Description: description,
		Priority:    priority,
		DueDate:     duedate,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
	}

	groupRef := r.coll.Doc(groupID)
	_, err := groupRef.Update(ctx, []firestore.Update{
		{Path: "tasks", Value: firestore.ArrayUnion(task)},
	})

	if err != nil {
		return models.Tasks{}, fmt.Errorf("failed to create task: %w", err)
	}
	return task, nil
}

func (r GroupRepository) DeleteTask(ctx context.Context, groupID string, taskID string) error {
	groupRef := r.coll.Doc(groupID)
	doc, err := groupRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group: %w", err)
	}

	var group models.Group
	if err := doc.DataTo(&group); err != nil {
		return fmt.Errorf("failed to decode group: %w", err)
	}

	found := false
	for i, task := range group.Tasks {
		if task.UUID == taskID {
			now := time.Now()
			group.Tasks[i].DeletedAt = &now
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task not found")
	}

	_, err = groupRef.Update(ctx, []firestore.Update{
		{Path: "tasks", Value: group.Tasks},
	})
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

func (r GroupRepository) CreateComment(ctx context.Context, groupID string, text string, createdBy models.Member) (models.Comments, error) {
	comment := models.Comments{
		UUID:      uuid.NewString(),
		Text:      text,
		CreatedBy: createdBy,
		CreatedAt: time.Now(),
	}

	groupRef := r.coll.Doc(groupID)
	_, err := groupRef.Update(ctx, []firestore.Update{
		{Path: "comments", Value: firestore.ArrayUnion(comment)},
	})

	if err != nil {
		return models.Comments{}, fmt.Errorf("failed to create commet: %w", err)
	}
	return comment, nil
}

func (r GroupRepository) DeleteComment(ctx context.Context, groupID string, commentID string) error {
	groupRef := r.coll.Doc(groupID)
	doc, err := groupRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group: %w", err)
	}

	var group models.Group
	if err := doc.DataTo(&group); err != nil {
		return fmt.Errorf("failed to decode group: %w", err)
	}

	found := false
	for i, comment := range group.Tasks {
		if comment.UUID == commentID {
			now := time.Now()
			group.Comments[i].DeletedAt = &now
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("comment not found")
	}

	_, err = groupRef.Update(ctx, []firestore.Update{
		{Path: "comment", Value: group.Comments},
	})
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	return nil
}

func (r GroupRepository) AddMember(ctx context.Context, groupID string, newMember models.Member) (models.Member, error) {
	groupRef := r.coll.Doc(groupID)
	_, err := groupRef.Update(ctx, []firestore.Update{
		{Path: "members", Value: firestore.ArrayUnion(newMember)},
	})

	if err != nil {
		return models.Member{}, fmt.Errorf("failed to create task: %w", err)
	}
	return newMember, nil
}

func (r GroupRepository) RemoveMember(ctx context.Context, groupID string, member models.Member) error {
	groupRef := r.coll.Doc(groupID)
	doc, err := groupRef.Get(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch group: %w", err)
	}

	var group models.Group
	if err := doc.DataTo(&group); err != nil {
		return fmt.Errorf("failed to decode group: %w", err)
	}

	value, exists := group.Members[member.ID]
	if exists {
		delete(group.Members, member.ID)
	} else {
		return fmt.Errorf("member '%s' not found", value.Name)
	}

	_, err = groupRef.Update(ctx, []firestore.Update{
		{Path: "Members", Value: group.Members},
	})
	if err != nil {
		return fmt.Errorf("failed to remove member: %w", err)
	}

	return nil
}
