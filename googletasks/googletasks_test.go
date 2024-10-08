package googletasks

import (
	"testing"

	"google.golang.org/api/tasks/v1"
)

func TestGetAllTasksGoogle(t *testing.T) {
	// taskListId := "VGJNRnByS3dTZk9aVy02MQ"
	taskListTitle := "Test obsidian_tasks"
	taskTitle1 := "Task with subs"
	taskTitle2 := "sub_task1"
	taskTitle3 := "sub_task2"

	allTasksMap := GetAllTasksGoogle(taskListTitle)

	_, ok1 := allTasksMap["|"+taskTitle1]
	_, ok2 := allTasksMap["|"+taskTitle2]
	_, ok3 := allTasksMap["dettaglio2|"+taskTitle3]

	if !ok1 {
		t.Fatalf(`%v missing in %v`, taskTitle1, allTasksMap)
	}
	if !ok2 {
		t.Fatalf(`%v missing in %v`, taskTitle2, allTasksMap)
	}
	if !ok3 {
		t.Fatalf(`%v missing in %v`, taskTitle3, allTasksMap)
	}
}

func TestDoneTaskGoogle(t *testing.T) {
	wantTaskListId := "VGJNRnByS3dTZk9aVy02MQ"
	taskListTitle := "Test obsidian_tasks"
	taskTitle := "Testing task update with Go"
	taskNotes := "Notes are here?"
	taskIdMd := taskNotes + "|" + taskTitle

	taskGoogle := tasks.Task{
		Title: taskTitle,
		Notes: taskNotes,
	}

	taskListId, err := GetTasksListId(taskListTitle)
	if err != nil && taskListId == wantTaskListId {
		t.Fatalf(`Wasn't able to retrieve the TaskListId:
		%v != %v, error: %v`, wantTaskListId, taskListId, err)
	}

	wantTask, err := AddTaskGoogle(taskListId, &taskGoogle)
	if err != nil {
		t.Fatalf(`Got error: %v`, err)
	}

	wantTask.Status = "completed"
	updateTaskGoogle, err := DoneTaskGoogle(taskListId, wantTask.Id, wantTask)

	if updateTaskGoogle.Status != wantTask.Status {
		t.Fatalf(`Status of %v != %v`, updateTaskGoogle.Title, wantTask.Title)
	}

	allTasksMap := GetAllTasksGoogle(taskListTitle)

	if allTasksMap[taskIdMd].Status != wantTask.Status {
		t.Fatalf(`Status of %v != %v`, allTasksMap[taskIdMd].Title, wantTask.Title)
	}
}

func TestAddTaskGoogle(t *testing.T) {
	taskListId := "VGJNRnByS3dTZk9aVy02MQ"
	taskTitle := "Testing task creation in Go"
	taskNotes := "path will be here again"

	taskGoogle := tasks.Task{
		Title: taskTitle,
		Notes: taskNotes,
	}

	retTaskGoogle, err := AddTaskGoogle(taskListId, &taskGoogle)
	if err != nil {
		t.Fatalf(`Got error: %v`, err)
	}
	if taskTitle != retTaskGoogle.Title && taskNotes != retTaskGoogle.Notes {
		t.Fatalf(`title %v != %v or notes %v != %v`, taskTitle, retTaskGoogle.Title, taskNotes, retTaskGoogle.Notes)
	}
}

func TestSetParent(t *testing.T) {
	taskListId := "VGJNRnByS3dTZk9aVy02MQ"

	taskParent := "This is parent test"
	taskChild := "Child of 'This is parent test'"

	taskParentG := tasks.Task{
		Title: taskParent,
	}
	taskParentGoogle, err := AddTaskGoogle(taskListId, &taskParentG)
	if err != nil {
		t.Fatalf(`Got error: %v`, err)
	}
	taskChildG := tasks.Task{
		Title: taskChild,
	}
	taskChildGoogle, err := AddTaskGoogle(taskListId, &taskChildG)
	if err != nil {
		t.Fatalf(`Got error: %v`, err)
	}

	taskMoved, err := SetParentGoogle(taskListId, taskChildGoogle.Id, taskParentGoogle.Id)

	if err != nil && taskMoved.Id != taskChildGoogle.Id {
		t.Fatalf(`Parent not set %v`, err)
	}
}
