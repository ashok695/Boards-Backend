// services/board_service.go
package services

import (
    "context"
    "fmt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "time"
	"go.mongodb.org/mongo-driver/bson/primitive"
    "boards/internals/models"
    "boards/db" // Import db package to access the Client
	"sync"
	// "reflect"
)
// Helper Function
func createMapForUSerAndStatus(assignedToData []models.AssignedToStruct, statusData []models.StatusDataStruct)models.Mapstruct{
	assignedToMap := make(map[primitive.ObjectID]models.AssignedToStruct)
	statusMap := make(map[primitive.ObjectID]models.StatusDataStruct)
	for _,user := range assignedToData {
		assignedToMap[user.ID] = user
	}
	for _,status := range statusData {
		statusMap[status.ID] = status
	}
	return models.Mapstruct{
		AssignedToMap :assignedToMap,
		StatusMap:statusMap,
	}
}
func createMapForRWT(rwtData []models.RWTStruct)map[primitive.ObjectID]models.RWTStruct{
	rwtMap := make(map[primitive.ObjectID]models.RWTStruct)
	for _, rwt:=range rwtData {
		rwtMap[rwt.ID] = rwt
	}
	return rwtMap
}
func createBoardMap(boardData []models.BoardDataStruct)map[primitive.ObjectID]*models.BoardDataStruct{
	boardMap:= make(map[primitive.ObjectID]*models.BoardDataStruct)
	for board := range boardData {
		boardMap[boardData[board].ID] = &boardData[board]
	}
	
	return boardMap
}
func assignOwnerAndStatus(taskData *models.TaskDataStruct, userMap map[primitive.ObjectID]models.AssignedToStruct, statusMap map[primitive.ObjectID]models.StatusDataStruct){
	// fmt.Println("taskData",taskData)
	if assignedTo,ok := taskData.AssignedTo.(primitive.ObjectID);ok{
		if user,exists := userMap[assignedTo];exists{
			taskData.AssignedTo = user
		}
	}
	if status,exists := statusMap[taskData.Status.(primitive.ObjectID)];exists{
		taskData.Status = status
	}
}
func asisgnRWTvalues(taskData *models.TaskDataStruct,rwtMap map[primitive.ObjectID]models.RWTStruct){
	 if isRoleArray,ok := taskData.Role.(primitive.A);ok{
		if len(isRoleArray) > 0 {
			if roleOk,ok:=isRoleArray[0].(primitive.ObjectID);ok{
				if role,exists := rwtMap[roleOk];exists{
					taskData.Role = role
				}
			}
		}
	 }
	 if isTaskTypeArray,ok := taskData.TaskType.(primitive.A);ok{
		if len(isTaskTypeArray) > 0 {
			if taskTypeOk,ok:=isTaskTypeArray[0].(primitive.ObjectID);ok{
				if taskType,exists := rwtMap[taskTypeOk];exists{
					taskData.TaskType = taskType
				}
			}
		}
	 }
	 if isWorkStreamArray,ok := taskData.WorkStream.(primitive.A);ok{
		if len(isWorkStreamArray) > 0 {
			if workstreamok,ok:=isWorkStreamArray[0].(primitive.ObjectID);ok{
				if workstream,exists := rwtMap[workstreamok];exists{
					taskData.WorkStream = workstream
				}
			}
		}
	 }
}
func calculateOverdueCount(taskData []models.TaskDataStruct,statusMap map[primitive.ObjectID]models.StatusDataStruct) int {
	currentTime := time.Now()
	count := 0

	for _, task := range taskData {
		var plannedTo time.Time

		// Handle PlannedTo as primitive.DateTime
		if plannedToPrimitive, ok := task.PlannedTo.(primitive.DateTime); ok {
			plannedTo = plannedToPrimitive.Time()
		} else {
			// If PlannedTo is not primitive.DateTime, skip this task
			fmt.Println("PlannedTo is not a valid primitive.DateTime")
			continue
		}
		// fmt.Println("task status",task.Status.Status)
		// fmt.Println("task.Status",reflect.TypeOf(task.Status))
		// Check if the task is overdue
		if taskStatus,ok := task.Status.(models.StatusDataStruct);ok{
			if plannedTo.Before(currentTime) && (taskStatus.Category == "New" || taskStatus.Category == "Active") {
				count++
			}
		}
	}
	// fmt.Println("count of overdue tasks is", count)
	// fmt.Println("count of overdue tasks is", count)
	// fmt.Println("count of overdue tasks is", count)
	// fmt.Println("count of overdue tasks is", count)
	fmt.Println("count of overdue tasks is", count)
	return count
}


func GetBoardDetails() ([]models.TaskDataStruct, error) {
    // Use the Client from the db package
	fmt.Println("HELLLLLO")
    // fmt.Println("i am client", db.Client)
	var (
		boardData      []models.BoardDataStruct
		taskData       []models.TaskDataStruct
		assignedToData []models.AssignedToStruct
		statusData     []models.StatusDataStruct
		rwtData        []models.RWTStruct
		// count int
	)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    
    boardDB := db.Client.Database("google-wp").Collection("kt_m_boards")
	taskDB := db.Client.Database("google-wp").Collection("kt_t_taskLists")
	assignedToDB := db.Client.Database("google-wp").Collection("kt_m_users")
	statusDB := db.Client.Database("google-wp").Collection("kt_m_status")
	rwtDB := db.Client.Database("google-wp").Collection("kt_m_types")


    boardPipeline := mongo.Pipeline{
        {{"$match", bson.D{{"category", "custom"}, {"active", true}}}},
        {{"$project", bson.D{{"_id", 1}, {"title", 1}, {"category", 1}}}},
    }
	taskPipeline := mongo.Pipeline{
		{{"$match", bson.D{{"refBoardID", bson.D{{"$exists", true}}}, {"board", "custom"}, {"skip", false}}}},
		{{"$project", bson.D{{"_id", 1}, {"title", 1}, {"status", 1}, {"assignedTo", 1}, {"refBoardID", 1}, {"plannedFrom", 1}, {"plannedTo", 1}, {"startedOn", 1}, {"completedOn", 1}, {"revisedStartDate", 1}, {"revisedEndDate", 1}, {"participants", 1}, {"tags", 1}, {"role", 1}, {"type", 1}, {"workstream", 1}}}},
	}
	assignedToPipeline := mongo.Pipeline{
		{{"$project", bson.D{{"_id", 1}, {"fullName", 1}, {"email", 1}}}},
	}
	statusPipeline := mongo.Pipeline{
		{{"$match", bson.D{{"workItem", "Task"}}}},
		{{"$project", bson.D{{"_id", 1}, {"category", 1}, {"status", 1}, {"workItem", 1}}}},
	}
	rwtPipeline := mongo.Pipeline{
		{{"$match", bson.D{
			{"$or", bson.A{
				bson.D{{"__type", "role"}},
				bson.D{{"__type", "tasktype"}},
				bson.D{{"__type", "workstream"}},
			}},
		}}},
		{{"$project", bson.D{{"_id", 1}, {"name", 1}, {"__type", 1}}}},
	}
	var wg sync.WaitGroup
	errChan := make(chan error, 4)
	taskChan := make(chan error)
	// Add 5 for the number of goroutines
	wg.Add(4)
	go func() {
		defer wg.Done()
		getDBData(ctx, boardDB, boardPipeline, &boardData, errChan)
	}()
	go func() {
		defer wg.Done()
		getDBData(ctx, assignedToDB, assignedToPipeline, &assignedToData, errChan)
	}()
	go func() {
		defer wg.Done()
		getDBData(ctx, statusDB, statusPipeline, &statusData, errChan)
	}()
	go func() {
		defer wg.Done()
		getDBData(ctx, rwtDB, rwtPipeline, &rwtData, errChan)
	}()
	go func(){
		cursor,err := taskDB.Aggregate(ctx,taskPipeline)
		if err != nil {
			fmt.Println("ERROR IN TASK MONGO")
			taskChan <- fmt.Errorf("Error in mongo %v",err)
		}
		for cursor.Next(ctx){
			var data models.TaskDataStruct 
			decodedError := cursor.Decode(&data)
			if decodedError!= nil {
				fmt.Println("ERROR IN DECODING DATA")
				taskChan <- fmt.Errorf("Error in mongo %v",decodedError)
			}
			taskData = append(taskData,data)
		}
		taskChan <- nil
	}()
	// Wait for DB goroutines to complete
	go func() {
		wg.Wait()
		close(errChan)
	}()
	
	// Collect errors
	for err := range errChan {
		if err != nil {
			return nil, fmt.Errorf("error running aggregation: %v", err)
		}
	}
	// Task Chan 
	if err := <-taskChan; err != nil {
		fmt.Println("error")
	}
	
	// Creating Maps
	maps := createMapForUSerAndStatus(assignedToData,statusData)
	rwtMap := createMapForRWT(rwtData)
	boardMap:= createBoardMap(boardData)

	for task := range taskData {
		assignOwnerAndStatus(&taskData[task], maps.AssignedToMap,maps.StatusMap)
		asisgnRWTvalues(&taskData[task],rwtMap)
	}

	 for _, task := range taskData {
        if board, exists := boardMap[task.RefBoardID[0]]; exists {
            board.Tasks = append(board.Tasks, task)
        }
    }
	calculateOverdueCount(taskData,maps.StatusMap)
	// calculateStakeHolderDetails(taskData,maps.StatusMap)
	return taskData,nil
}
func getDBData(ctx context.Context, coll *mongo.Collection, pipeline mongo.Pipeline, result interface{}, errChan chan error) {
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		errChan <- fmt.Errorf("aggregation error: %w", err)
		return
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, result); err != nil {
		errChan <- fmt.Errorf("cursor error: %w", err)
		return
	}
	errChan <- nil
}
func GetBoardTasksDetails() (string, error) {
    hello := "I AM STARBOYYYYYYYYYYYYYYY"
    return hello, nil
}
