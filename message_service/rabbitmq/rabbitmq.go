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
	"fmt"

	"github.com/streadway/amqp"
	"urulink.go/message_service/env"
)

// RabbitMQManager manages the connection and channel to RabbitMQ
type RabbitMQManager struct {
	conn    *amqp.Connection // connection to the RabbitMQ server
	channel *amqp.Channel    // channel for communication with RabbitMQ
}

// UruLinkInit initializes RabbitMQ connection and declares the exchange and queue for message handling
func UruLinkInit(envManager *env.EnvManger) (*RabbitMQManager, error) {
	// Format the RabbitMQ URL with environment configuration settings
	rabbitMQURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		envManager.RabbitMQUser,
		envManager.RabbitMQPassword,
		envManager.RabbitMQHost,
		envManager.RabbitMQPort,
	)

	// Establish a connection to RabbitMQ
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err // return error if connection fails
	}

	// Open a channel on the RabbitMQ connection
	ch, err := conn.Channel()
	if err != nil {
		return nil, err // return error if channel creation fails
	}

	// Declare a topic exchange for routing messages
	err = ch.ExchangeDeclare(
		envManager.RabbitMQExchangeName, // the name of the exchange
		"topic",                         // type of exchange (topic-based)
		true,                            // durable - persists after server restarts
		false,                           // auto-delete when no consumers are bound
		false,                           // internal - not used by publishers
		false,                           // no-wait - wait for confirmation
		nil,                             // additional arguments
	)
	if err != nil {
		return nil, err // return error if exchange declaration fails
	}

	// Declare a queue to hold messages before consumption
	_, err = ch.QueueDeclare(
		envManager.RabbitMQQueueName, // name of the queue
		true,                         // durable - persists after server restarts
		false,                        // not auto-deleted
		false,                        // not exclusive to this connection
		false,                        // no-wait - wait for confirmation
		nil,                          // additional arguments
	)
	if err != nil {
		return nil, err // return error if queue declaration fails
	}

	// Bind the queue to the exchange with a routing key pattern
	err = ch.QueueBind(
		envManager.RabbitMQQueueName,    // name of the queue
		"#",                             // binding key to match any routing key
		envManager.RabbitMQExchangeName, // name of the exchange
		false,                           // no-wait - wait for confirmation
		nil,                             // additional arguments
	)
	if err != nil {
		return nil, err // return error if queue binding fails
	}

	// Return a new RabbitMQManager instance with the connection and channel
	return &RabbitMQManager{
		conn:    conn,
		channel: ch,
	}, nil
}
