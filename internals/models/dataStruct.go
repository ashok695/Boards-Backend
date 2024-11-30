package models 
import(
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type BoardDataStruct struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	Category string             `json:"category" bson:"category"`
	Tasks []TaskDataStruct      `json:"tasks" bson:"tasks"`
}
type TaskDataStruct struct {
	ID             primitive.ObjectID   `json:"_id" bson:"_id"`
	Title          string               `json:"title" bson:"title"`
	Status         interface{}          `json:"status" bson:"status"`
	AssignedTo     interface{}          `json:"assignedTo" bson:"assignedTo"`
	RefBoardID     []primitive.ObjectID `json:"refBoardID" bson:"refBoardID"`
	PlannedFrom    interface{}          `json:"plannedFrom" bson:"plannedFrom"`
	PlannedTo      interface{}          `json:"plannedTo" bson:"plannedTo"`
	ActualStart    interface{}          `json:"startedOn" bson:"startedOn"`
	ActualEnd      interface{}          `json:"completedOn" bson:"completedOn"`
	RevisedStart   interface{}          `json:"revisedStartDate" bson:"revisedStartDate"`
	RevisedEnd     interface{}          `json:"revisedEndDate" bson:"revisedEndDate"`
	Participants   interface{}          `json:"participants" bson:"participants"`
	Tags           interface{}          `json:"tags" bson:"tags"`
	Role           interface{}          `json:"role" bson:"role"`
	TaskType       interface{}          `json:"type" bson:"type"`
	WorkStream     interface{}          `json:"workstream" bson:"workstream"`
}
type AssignedToStruct struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	FullName string             `json:"fullName" bson:"fullName"`
	Email    string             `json:"email" bson:"email"`
}
type StatusDataStruct struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Category string             `json:"category" bson:"category"`
	Status   string             `json:"status" bson:"status"`
	WorkItem string             `json:"workItem" bson:"workItem"`
}
type Mapstruct struct {
	AssignedToMap map[primitive.ObjectID]AssignedToStruct
	StatusMap map[primitive.ObjectID]StatusDataStruct
}
type RWTStruct struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	Name string `json:"name" bson:"name"`
	Type string `json:"__type" bson:"__type"`
}
type TaskBasedOnStatus struct {
	ID primitive.ObjectID `json:"_id" bson:"_id"`
	Category string `json:"category" bson:"category"`
	Status string `json:"status" bson:"status"`
	TaskCount int `json:"taskCount" bson:"taskCount"`
	PercentOfTotalTask int `json:"PercentOfTotalTask" bson:"PercentOfTotalTask"`
}