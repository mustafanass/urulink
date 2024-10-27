/*
 * Copyright (C) 2024 Mustafa Naseer
 *
 * This file is part of urulink chat application.
 *
 * urulink is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation version 3 of the License .
 *
 *
 * urulink is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with urulink. If not, see <http://www.gnu.org/licenses/>.
 */

package rabbitmq

import (
	"context"
	"fmt"

	"github.com/streadway/amqp"
	"urulink.go/message_service/env"
	"urulink.go/message_service/helper"
)

// PublishMessage sends a message to the specified RabbitMQ exchange with the given user ID as the routing key.
func (rm *RabbitMQManager) PublishMessage(ctx context.Context, userId string, messageBody []byte, envManager *env.EnvManger) error {
	// Publish a message to RabbitMQ using the configured exchange name and routing key.
	err := rm.channel.Publish(
		envManager.RabbitMQExchangeName, // name of the exchange
		userId,                          // routing key (target user ID)
		true,                            // mandatory flag
		false,                           // immediate flag
		amqp.Publishing{
			ContentType: "application/json", // setting the content type to JSON
			Body:        messageBody,        // message body content
		},
	)
	return err // return any error that occurs during publishing
}

// ListenForMessages subscribes to a RabbitMQ queue for the specified user ID, processing messages with the provided handler.
func (rm *RabbitMQManager) ListenForMessages(ctx context.Context, userId string, handler func([]byte) error, envManager *env.EnvManger) {
	// Consume messages from the RabbitMQ queue using the configured queue name and consumer tag.
	msgs, err := rm.channel.Consume(
		envManager.RabbitMQQueueName, // name of the queue
		userId,                       // consumer tag (typically a unique identifier for the consumer)
		false,                        // auto-acknowledge flag (set to false for manual acknowledgment)
		false,                        // exclusive flag
		false,                        // no-local flag
		false,                        // no-wait flag
		nil,                          // additional arguments
	)
	if err != nil {
		// Log an error if message consumption fails
		helper.LogError(nil, fmt.Sprintf("Failed to consume messages for userId %s", userId), err)
		return
	}

	// Start a goroutine to listen for incoming messages asynchronously.
	go func() {
		for msg := range msgs {
			// Process each message using the provided handler function
			err := handler(msg.Body)
			if err != nil {
				// Log an error and requeue the message if the handler returns an error
				helper.LogError(nil, fmt.Sprintf("Error processing message for userId %s", userId), err)
				msg.Nack(false, true) // negative acknowledgment to requeue the message
			} else {
				// Acknowledge the message if processing is successful
				if err := msg.Ack(false); err != nil {
					helper.LogError(nil, fmt.Sprintf("Failed to acknowledge message for userId %s", userId), err)
				}
			}
		}
	}()

	// Log that the message listener has started for the user
	helper.LogInfo("Started listening for messages", map[string]interface{}{"userId": userId})
}

// CancelConsume stops consuming messages for the specified consumer tag in RabbitMQ.
func (rm *RabbitMQManager) CancelConsume(consumerTag string) error {
	// Cancel the consumer in RabbitMQ using the provided consumer tag.
	if err := rm.channel.Cancel(consumerTag, false); err != nil {
		// Return an error if cancellation fails
		return fmt.Errorf("[CancelConsume] Failed to cancel consumer %s: %v", consumerTag, err)
	}
	// Log the successful cancellation of the consumer
	helper.LogInfo("Consumer cancelled successfully", map[string]interface{}{
		"consumerTag": consumerTag,
	})
	return nil
}
