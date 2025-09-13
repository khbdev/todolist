
# Golang TodoList Project

This is a **TodoList** application built with **Golang** following **Clean Architecture** principles. The project demonstrates a modern backend setup using Gin, GORM, MySQL, Redis, JWT authentication, and RabbitMQ for asynchronous tasks.

## Features

* **Gin & GORM**

  * RESTful API built with Gin
  * MySQL database with relationships

* **Authentication**

  * JWT access and refresh tokens for secure authentication

* **Caching with Redis**

  * Implements **read-through** and **write-through** caching patterns

* **RabbitMQ Integration**

  * Sends a "Welcome to TodoList" email to users upon registration
  * Uses **exchange, queue binding, retry mechanism**, and **DLQ (Dead Letter Queue)** for reliable message processing

* **Scheduled Cron Jobs**

  * Uses **Anqinq** library to schedule daily jobs
  * Sends a reminder email to all users once every 24 hours to encourage them to add new todos

* **Real-time Notifications**

  * Powered by **Redis Pub/Sub** and **Gorilla WebSocket**
  * Frontend instantly receives updates (create, update, delete) whenever a todo is modified

* **Clean Architecture**

  * Separation of concerns with `domain`, `usecase`, `repository`, and `handler` layers
  * Maintainable and scalable project structure

## Technologies Used

* Golang
* Gin
* GORM
* MySQL
* Redis
* RabbitMQ
* JWT (Access & Refresh tokens)
* Gorilla WebSocket
* Anqinq

