package main

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"time"
)

type SpoolService struct {
	Mail *MailService
}

func NewSpoolService(mailService *MailService) *SpoolService {
	spoolService = &SpoolService{
		Mail: mailService,
	}

	return spoolService
}

var spoolService *SpoolService
var queue *sqs.SQS
var queueName string

func StartSQSSpool() {
	for {
		time.Sleep(time.Minute)
		PollSqsMessage()
	}
}

func PollSqsMessage() {
	session := GetAWSSession()
	queue = sqs.New(session)
	queueName, _ = GetQueueName()

	result, err := queue.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &queueName,
		MaxNumberOfMessages: aws.Int64(10),
		VisibilityTimeout:   aws.Int64(360),
		WaitTimeSeconds:     aws.Int64(20),
	})
	if err != nil {
		log.Printf("[ERROR] [%s] SQS Error: %s", ServiceName, err.Error())
		time.Sleep(time.Minute)
		return
	}

	for index, msg := range result.Messages {
		log.Printf("[INFO] [%s] Processing message %d of %d recieved from queue", ServiceName, index+1, len(result.Messages))
		ProcessSQSMessage(msg)
	}
}

func ProcessSQSMessage(msg *sqs.Message) {
	snsMsg := &sns.PublishInput{}
	if err := json.Unmarshal(bytes.NewBufferString(*msg.Body).Bytes(), snsMsg); err != nil {
		log.Printf("[ERROR] [%s] SQS Error: %s", ServiceName, err.Error())
		return
	}
	if snsMsg.Message == nil {
		log.Printf("[ERROR] [%s] No message body", ServiceName)
		return
	}
	buf := bytes.NewBufferString(*snsMsg.Message).Bytes()
	emsg := &Message{}
	err := json.Unmarshal(buf, emsg)
	if nil != err {
		log.Printf("[ERROR] [%s] Unmarshall Error: %s", ServiceName, err.Error())
	}
	spoolService.Mail.QueueEmail(emsg)

	_, err = queue.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &queueName,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if nil != err {
		log.Printf("[ERROR] [%s] SQS Error: %s", ServiceName, err.Error())
	}
}
