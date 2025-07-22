package controllers

// imports
import (
	"net/http";
	"strings";
	"github.com/gin-gonic/gin";
	"github.com/natnael-eyuel-dev/Task-Management-Clean-Architecture/Domain";
	"github.com/natnael-eyuel-dev/Task-Management-Clean-Architecture/Usecases";
	"go.mongodb.org/mongo-driver/bson/primitive";
)

// task controller
type TaskController struct {
	taskUseCase usecases.TaskUseCase        // task usecase for task operations
}

// user controller
type UserController struct {
	userUseCase usecases.UserUseCase        // user usecase for user operations 
}

// new task controller
func NewTaskController(uc usecases.TaskUseCase) *TaskController {
	return &TaskController{taskUseCase: uc}        // return new task controller instance
}

// new user controller
func NewUserController(uc usecases.UserUseCase) *UserController {
	return &UserController{userUseCase: uc}        // return new user controller instance
}

func (taskContr *TaskController) CreateTask(c *gin.Context) {
	
	var task domain.Task
	err := c.ShouldBindJSON(&task)      // parse request body into task struct
	if err != nil {
		// handle specific date format error case
		if strings.Contains(err.Error(), "numeric literal") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid date format. Use ISO 8601 format like '2023-12-31T00:00:00Z'",
				"example": gin.H{
					"due_date": "2025-07-22T00:00:00Z",
				},
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create task through usecase layer
	createdTask, err := taskContr.taskUseCase.CreateTask(&task)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdTask)        // return created task with 201 status
}

func (taskContr *TaskController) DeleteTask(c *gin.Context) {
	
	id := c.Param("id")       // get task id from request parameter

	_, err := primitive.ObjectIDFromHex(id)       // validate it is a valid ObjectID 
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	// delete task through usecase layer
	err = taskContr.taskUseCase.DeleteTask(id)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message":"task deleted successfully"})    // success response
}

func (taskContr *TaskController) GetAllTasks(c *gin.Context) {
	
	// get all tasks through usecase layer
	tasks, err := taskContr.taskUseCase.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(tasks) == 0 {
		c.JSON(http.StatusOK, []domain.Task{})
		return
	}

	c.JSON(http.StatusOK, tasks)       // return all tasks
}

func (taskContr *TaskController) GetTaskByID(c *gin.Context) {
	
	id := c.Param("id")        // get task id from request parameter

	_, err := primitive.ObjectIDFromHex(id)      // validate it is a valid ObjectID
	if err != nil {      
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	// get specific task through usecase layer
	task, err := taskContr.taskUseCase.GetTaskByID(id)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)       // return found task 
}

func (taskContr *TaskController) UpdateTask(c *gin.Context) {
	
	id := c.Param("id")       // get task id from request parameter

	_, err := primitive.ObjectIDFromHex(id)        // validate it is a valid ObjectID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID format"})
		return
	}

	var task domain.Task
	err = c.ShouldBindJSON(&task)       // parse request body into task struct
	if err != nil {
		// handle specific date format error case
		if strings.Contains(err.Error(), "numeric literal") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid date format. Use ISO 8601 format like '2025-7-16T00:00:00Z'",
				"example": gin.H{
					"due_date": "2025-07-22T00:00:00Z",
				},
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// update task through usecase layer
	updatedTask, err := taskContr.taskUseCase.UpdateTask(id, &task)
	if err != nil {
		if err == domain.ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})       
		return
	}

	c.JSON(http.StatusOK, gin.H{ "message":"task updated successfully", "updated_task":updatedTask})       // success response
}

func (uc *UserController) Register(c *gin.Context) {
	
	var user domain.User
	err := c.ShouldBindJSON(&user)       // parse request body into user struct
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// create user through usecase layer
	if err := uc.userUseCase.Register(&user); err != nil {
		if err == domain.ErrUserExists {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})       // success response
}

func (uc *UserController) Login(c *gin.Context) {
	
	var creds domain.Credentials
	err := c.ShouldBindJSON(&creds)        // parse request body into user struct
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// authenticate user through usecase layer
	token, user, err := uc.userUseCase.Login(&creds)
	if err != nil {
		if err == domain.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// return token, user info (excluding sensitive data)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func (uc *UserController) PromoteToAdmin(c *gin.Context) {
	
	userID := c.Param("id")       // get user id from request parameter
	 
	_, err := primitive.ObjectIDFromHex(userID)       // validate it is a valid ObjectID
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		return
	}

	// promote user through usecase layer
	err = uc.userUseCase.PromoteToAdmin(userID) 
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin successfully"})       // success response
}