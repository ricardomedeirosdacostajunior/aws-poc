package awsrepository

import (
	"aws-poc/pkg/awssession"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var (
	hashKeyType   = aws.String("HASH")
	stringType    = aws.String("S")
	payPerRequest = aws.String("PAY_PER_REQUEST")
)

type (
	record interface {
		ID() string
	}
	mapMarshaller   func(in interface{}) (map[string]*dynamodb.AttributeValue, error)
	mapUnmarshaller func(m map[string]*dynamodb.AttributeValue, out interface{}) error
	adapter         interface {
		PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error)
		DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error)
		GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error)
	}
	repository interface {
		put(rec record) error
		delete(rec record) error
	}
	dynamoRepository struct {
		sess      *session.Session
		tableName *string
		mapMarshaller
		mapUnmarshaller
		adapter
	}
)

func newRegister(sess *session.Session, tableName *string) dynamoRepository {
	return dynamoRepository{sess, tableName, dynamodbattribute.MarshalMap, dynamodbattribute.UnmarshalMap, svc()}
}

func svc() (svc *dynamodb.DynamoDB) {
	sess := awssession.NewLocalSession()
	svc = dynamodb.New(sess)
	return
}

func (r dynamoRepository) put(rec record) error {
	item, err := r.mapMarshaller(rec)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      item,
		TableName: r.tableName,
	}

	item["ID"] = &dynamodb.AttributeValue{
		S: aws.String(rec.ID()),
	}

	if _, err = r.PutItem(input); err != nil {
		return err
	}

	return nil
}

func (r dynamoRepository) delete(rec record) error {
	input := &dynamodb.DeleteItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(rec.ID()),
			},
		},
	}

	if _, err := r.DeleteItem(input); err != nil {
		return err
	}

	return nil
}

func (r dynamoRepository) get(rec record, item interface{}) (interface{}, error) {
	input := &dynamodb.GetItemInput{
		TableName: r.tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(rec.ID()),
			},
		},
	}

	ri, err := r.GetItem(input)
	if err != nil {
		return nil, err
	}

	if err := r.mapUnmarshaller(ri.Item, &item); err != nil {
		return nil, err
	} else {
		return item, err
	}
}
