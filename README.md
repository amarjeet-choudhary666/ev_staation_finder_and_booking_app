//Work in progress

// The EV Finder Station and Booking App is a backend RESTful API service designed to help electric vehicle users locate nearby charging stations, view details, and book charging slots efficiently. Built using Golang, this application leverages powerful tools such as Gin for routing, GORM for database interaction with PostgreSQL, and Bcrypt for secure password hashing.

🔐 User Authentication & Authorization

Secure sign up and login using JWT

Password encryption with Bcrypt

Role-based access (admin, user)

🗺️ Charging Station Management

Add, update, delete, and retrieve EV charging stations

Filter stations by location, availability, and charger type

📅 Booking System

Book charging slots

View and manage user bookings

Admin slot and station management

📍 Search and Filter

Geolocation-based search for nearby stations

Filtering based on charger type, availability, and ratings

🧾 Admin Panel (API-based)

Manage stations, users, and bookings

Analytics and usage reports (optional future extension)


🛠️ Tech Stack
Technology	Purpose
Go (Golang)	Core backend language
Gin	HTTP web framework for RESTful APIs
PostgreSQL	Relational database
GORM	ORM for database handling
Bcrypt	Password hashing and authentication
JWT	Secure token-based authentication


RUN PROJECT: go run cmd/main.go || make run