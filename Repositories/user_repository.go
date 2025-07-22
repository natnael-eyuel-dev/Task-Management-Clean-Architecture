package repositories

// imports
import (
	"context";
	"time";
	"go.mongodb.org/mongo-driver/bson";
	"go.mongodb.org/mongo-driver/bson/primitive";
	"go.mongodb.org/mongo-driver/mongo";
	"github.com/natnael-eyuel-dev/Task-Management-Clean-Architecture/Domain";
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(col *mongo.Collection) domain.UserRepository {
	return &userRepository{collection: col}
}

//  register user in to database
func (userRepo *userRepository) CreateUser(user *domain.User) error {
	
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	// generate new ObjectID if not set
	if user.ID.IsZero() {
		user.ID = primitive.NewObjectID()
	}

	// save user to database
	_, err := userRepo.collection.InsertOne(contx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrUserExists
		}
		return err
	}

	return nil        // success
}

// find user from database by username
func (userRepo *userRepository) GetByUsername(username string) (*domain.User, error) {
	
	var user domain.User
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()
	
	// find user by username
	err := userRepo.collection.FindOne(contx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil        // success
}

// find user from database by id
func (userRepo *userRepository) GetUserById(userID primitive.ObjectID) (*domain.User, error) {
	
	var user domain.User
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()
	
	// find user by id
	err := userRepo.collection.FindOne(contx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil         // success
}

// count users in the database currently
func (userRepo *userRepository) GetUserCount() (int64, error) {
	
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	// count users in user collection currently
	count, err := userRepo.collection.CountDocuments(contx, bson.M{})
	if err != nil {
		return 0, err
	}

	return count, nil        // success
}

// update user role to admin in database (only admins can perform this operation)
func (userRepo *userRepository) UpdateRole(id primitive.ObjectID, role string) error {
	
	contx, cancel := context.WithTimeout(context.Background(), 5*time.Second)        // set timeout
	defer cancel()

	// update user's role to admin
	result, err := userRepo.collection.UpdateOne(
		contx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"role": role}},
	)

	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return domain.ErrUserNotFound
	}

	return nil        // success
}