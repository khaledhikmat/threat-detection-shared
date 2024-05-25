package pubsub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/khaledhikmat/threat-detection-shared/models"
	"github.com/khaledhikmat/threat-detection-shared/service/config"
)

const (
	region = "us-east-2"
)

func NewAwsPubsub(cfgsvc config.IService) IService {
	c, err := makeSNSClient(context.Background())
	if err != nil {
		fmt.Printf("Couldn't create SNS client. Here's why: %v\n", err)
		panic(err)
	}

	q, err := makeSQSClient(context.Background())
	if err != nil {
		fmt.Printf("Couldn't create SQS client. Here's why: %v\n", err)
		panic(err)
	}

	return &awsPubsub{
		SnsClient: c,
		SqsClient: q,
		CfgSvc:    cfgsvc,
	}
}

type awsPubsub struct {
	SnsClient *sns.Client
	SqsClient *sqs.Client
	CfgSvc    config.IService
}

func (s *awsPubsub) CreateTopic(ctx context.Context, topicName string) (string, error) {
	// List all topics
	resp, err := s.SnsClient.ListTopics(ctx, &sns.ListTopicsInput{})
	if err != nil {
		return "", fmt.Errorf("failed to list topics: %w", err)
	}

	// Check if the topic already exists
	for _, t := range resp.Topics {

		if *t.TopicArn == topicName {
			// The topic already exists, return its ARN
			return *t.TopicArn, nil
		}
	}

	fmt.Printf("***** Creating topic %s\n", topicName)

	// If the topic doesn't exist, create it
	respCreate, err := s.SnsClient.CreateTopic(ctx, &sns.CreateTopicInput{
		Name: aws.String(topicName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create topic: %w", err)
	}

	// Return the topic ARN
	return *respCreate.TopicArn, nil
}

func (s *awsPubsub) CreateQueue(ctx context.Context, queueName, topicARN string) (string, string, error) {
	queueURL := ""
	newQueue := false

	// Try to get the URL of the queue
	resp, err := s.SqsClient.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(queueName),
	})

	if err != nil {
		var queueDoesNotExist *sqstypes.QueueDoesNotExist
		if errors.As(err, &queueDoesNotExist) {
			newQueue = true
			// If the queue doesn't exist, create it
			respCreate, err := s.SqsClient.CreateQueue(ctx, &sqs.CreateQueueInput{
				QueueName: aws.String(queueName),
			})
			if err != nil {
				return "", "", fmt.Errorf("failed to create queue: %w", err)
			}

			// Return the URL of the newly created queue
			queueURL = *respCreate.QueueUrl

		} else {
			// If the error is not QueueDoesNotExist, return the error
			return "", "", err
		}
	} else {
		// If the queue already exists, get its URL
		queueURL = *resp.QueueUrl
	}

	// Using the queue URL, Get the attributes of the queue so we can get its ARN
	respAttributes, err := s.SqsClient.GetQueueAttributes(ctx, &sqs.GetQueueAttributesInput{
		QueueUrl:       aws.String(queueURL),
		AttributeNames: []sqstypes.QueueAttributeName{sqstypes.QueueAttributeNameQueueArn},
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to get queue attributes: %w", err)
	}

	queueARN := respAttributes.Attributes[string(sqstypes.QueueAttributeNameQueueArn)]

	// Add a policy to allow SNS to send messages to the queue
	if newQueue {
		policy := map[string]interface{}{
			"Version": "2012-10-17",
			"Id":      "SNStoSQSPolicy",
			"Statement": []map[string]interface{}{
				{
					"Sid":       "Allow-SNS-SendMessage",
					"Effect":    "Allow",
					"Principal": map[string]string{"AWS": "*"},
					"Action":    "sqs:SendMessage",
					"Resource":  queueARN,
					"Condition": map[string]map[string]string{
						"ArnEquals": {"aws:SourceArn": topicARN},
					},
				},
			},
		}

		// Convert the policy to JSON
		policyJSON, err := json.Marshal(policy)
		if err != nil {
			fmt.Printf("failed to serialize policy: %v\n", err)
		}

		// Set the policy on the queue
		_, err = s.SqsClient.SetQueueAttributes(ctx, &sqs.SetQueueAttributesInput{
			QueueUrl: aws.String(queueURL),
			Attributes: map[string]string{
				"Policy": string(policyJSON),
			},
		})
		if err != nil {
			fmt.Printf("failed to attach policy: %v\n", err)
		}
	}

	// Return the queue URL and ARN
	return queueURL, queueARN, nil
}

func (s *awsPubsub) PublishRecordingClip(ctx context.Context, _, topicArn string, clip models.RecordingClip) error {
	// Serialize the clip to JSON
	clipJSON, err := json.Marshal(clip)
	if err != nil {
		return fmt.Errorf("failed to serialize clip: %w", err)
	}

	// Publish the clip to the topic
	input := sns.PublishInput{
		TopicArn: aws.String(topicArn),
		Message:  aws.String(string(clipJSON)),
	}

	_, err = s.SnsClient.Publish(ctx, &input)
	return err
}

type snsMessage struct {
	Message string `json:"Message"`
}

func (s *awsPubsub) Subscribe(ctx context.Context, topicArn, queueURL, queueArn string) (chan string, error) {
	stream := make(chan string)

	go func() {
		defer close(stream)

		// Subscribe the queue to the topic
		attributes := map[string]string{}
		_, err := s.SnsClient.Subscribe(ctx, &sns.SubscribeInput{
			Protocol:              aws.String("sqs"),
			TopicArn:              aws.String(topicArn),
			Attributes:            attributes,
			Endpoint:              aws.String(queueArn),
			ReturnSubscriptionArn: true,
		})
		if err != nil {
			fmt.Printf("Couldn't subscribe queue %v to topic %v. Here's why: %v\n", queueArn, topicArn, err)
			return
		}

		// Poll the queue for messages
		for {
			// If context is cancelled, stop polling
			if ctx.Err() != nil {
				fmt.Printf("Context cancelled. Stopping polling for messages from queue %v\n", queueArn)
				return
			}

			output, err := s.SqsClient.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
				QueueUrl:            aws.String(queueURL),
				MaxNumberOfMessages: 10,
				WaitTimeSeconds:     5, // Atomic operation should not last more than 5 seconds
			})
			if err != nil {
				fmt.Printf("Couldn't receive message from queue %v. Here's why: %v\n", queueArn, err)
				return
			}

			// The assumption is that we receive messages from SNS in batches of 10
			for _, message := range output.Messages {
				var snsMsg snsMessage
				err := json.Unmarshal([]byte(*message.Body), &snsMsg)
				if err != nil {
					fmt.Printf("Record but ignore - Failed to decode message from queue %v: %v\n", queueArn, err)
					continue
				}

				stream <- snsMsg.Message
				_, err = s.SqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
					QueueUrl:      aws.String(queueURL),
					ReceiptHandle: message.ReceiptHandle,
				})
				if err != nil {
					fmt.Printf("NOT NORMAL - record but ignore - Couldn't delete message from queue %v. Here's why: %v\n", queueArn, err)
				}
			}
		}
	}()

	return stream, nil
}

func (s *awsPubsub) Finalize() {
}

func makeSNSClient(ctx context.Context) (*sns.Client, error) {
	var client *sns.Client
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(region),
	)

	if err != nil {
		return nil, err
	}

	client = sns.NewFromConfig(cfg)

	return client, nil
}

func makeSQSClient(ctx context.Context) (*sqs.Client, error) {
	var client *sqs.Client
	cfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(region),
	)

	if err != nil {
		return nil, err
	}

	client = sqs.NewFromConfig(cfg)

	return client, nil
}
