package Repository

import (
	"Loan_manager/Domain"
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type LoanRepository interface {
	CreateLoan(loan Domain.Loan) error
	GetLoanByID(id primitive.ObjectID) (*Domain.Loan, error)
	GetAllLoans(status string, order string) ([]Domain.Loan, error)
	UpdateLoan(loan *Domain.Loan) error // Updated method signature
	DeleteLoan(id primitive.ObjectID) error
}

type loanRepository struct {
	collection *mongo.Collection
}

func NewLoanRepository(collection *mongo.Collection) LoanRepository {
	return &loanRepository{collection: collection}
}

func (lr *loanRepository) CreateLoan(loan Domain.Loan) error {
	_, err := lr.collection.InsertOne(context.TODO(), loan)
	return err
}

func (lr *loanRepository) GetLoanByID(id primitive.ObjectID) (*Domain.Loan, error) {
	var loan Domain.Loan
	err := lr.collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&loan)
	if err != nil {
		return nil, err
	}
	return &loan, nil
}

func (lr *loanRepository) GetAllLoans(status string, order string) ([]Domain.Loan, error) {
	filter := bson.M{}
	if status != "all" {
		filter["status"] = status
	}

	opts := options.Find()
	if order == "desc" {
		opts.Sort = bson.D{{Key: "created_at", Value: -1}}
	} else {
		opts.Sort = bson.D{{Key: "created_at", Value: 1}}
	}

	cursor, err := lr.collection.Find(context.TODO(), filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var loans []Domain.Loan
	for cursor.Next(context.TODO()) {
		var loan Domain.Loan
		if err := cursor.Decode(&loan); err != nil {
			return nil, err
		}
		loans = append(loans, loan)
	}

	return loans, cursor.Err()
}

func (lr *loanRepository) UpdateLoan(loan *Domain.Loan) error { // Accept pointer
	filter := bson.M{"_id": loan.ID}
	update := bson.M{
		"$set": bson.M{
			"status":      loan.Status,
			"approved_at": loan.ApprovedAt,
		},
	}

	_, err := lr.collection.UpdateOne(context.TODO(), filter, update)
	return err
}

func (lr *loanRepository) DeleteLoan(id primitive.ObjectID) error {
	_, err := lr.collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}
