# urulink Chat Application

## Overview

This is a **backend chat application** designed for private or local chat environments. Built using **Golang** and the **Fiber** framework, this application leverages modern technologies to facilitate real-time communication and efficient data management.

## Features

- **Real-Time Communication**: Utilizes **WebSocket** for instant messaging capabilities.
- **User Connection Management**: Employs **Redis** for storing and managing user connections.
- **Message Broker**: Integrates **RabbitMQ** for efficient message handling and queuing.
- **Object Storage**: Uses **MinIO** for object storage, ensuring that all media and files are securely managed.
- **Microservice Architecture**: Built using a microservice architecture to promote modularity and scalability.

## Technologies Used

- **Golang**: The primary programming language for backend development.
- **Fiber**: A fast, minimalist web framework for Golang, used to create the server.
- **WebSocket**: For real-time, bidirectional communication between clients and the server.
- **Redis**: A powerful in-memory data store used for managing user sessions and connections.
- **RabbitMQ**: A robust message broker that enables reliable communication between microservices.
- **MinIO**: A high-performance object storage solution compatible with Amazon S3 APIs.
- **JWT (JSON Web Tokens)**: Used for secure authentication and authorization of users in the application.

## Getting Started

To set up the chat application on your local machine, follow these steps:

### Prerequisites

- Install [Go](https://golang.org/dl/) (version 1.16 or higher)
- Install [Docker](https://www.docker.com/get-started)
- Install [Fiber](https://docs.gofiber.io)
- Install [RabbitMQ](https://www.rabbitmq.com/download.html)
- Install [Redis](https://redis.io/download)
- Install [MinIO](https://min.io/download)

### Installation

- Clone all microservices and use the Dockerfile inside each one to build the service.
- Make sure to fill all required environment variables in `.env` files before building the Docker image.
- Now WebSocket is started, you can access it from `ws://your-ip-or-domain:8083`.

## Note

**Please note that this application is a work in progress, and not all features are complete. There are many additional features planned for future development. If you have any features in mind that you would like to implement, feel free to fork the repository and create a pull request with your contributions!**

## License

This project is licensed under the GNU General Public License v3.0 (GPL-3.0). See the LICENSE file for more information.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, feel free to open an issue or submit a pull request.

## Author

Created by **Mustafa Naseer**. For more information, feel free to contact me.

---

Thank you for checking out this chat application! Enjoy real-time communication with your peers.
