# Microservices programming project
INFS605 - course project

# Task
- Osborne.AI
- Student Services Dashboard for university operations
- Microservices architecture - connected with DockerCompose, maybe Kubernetes if time permits

# Requirements
- Minimum 3 new services
- Docker and DockerCompose, git, Yaml, SQL or noSQL,
- Optional languages.

## Services (Ideas)
- Student Profile Service
- Course Catalogue Service
- Feedback Service
- Notification Service
- Grades Service
- Timetable Service
- Assignment Tracker
-  Room Booking Service (for team meetings), 
come more generic services such as – an Image Upload Service, a 
- Search Service, 
- a Logging or Audit Trail Service,
- Frontend UI,
- PDF Generator, 
- System Metrics Service. 
- API Gateway (Use Nginx or Traefik)

## Functional
- Must expose a least 2 RESTful APIs with 2 endpoints
- HTTP, TCP(RPC, maybe gRPC) and simple message queue 
- Small frontend for interaction
- Authentication to one or more services

## Non-Functional
- Each service in its own Docker container
- Must be defined in a DockerCompose file
- Include documentation and usage instructions
- System must support logging, basic error handing

