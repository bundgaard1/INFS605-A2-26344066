Here is the complete, actionable rollout plan for your **INFS605 project**, translated into English and structured around a **vertical slice strategy** (building one end-to-end feature at a time to avoid the "empty microservices" trap).

---

## Target Architecture & Service Selection

To meet your requirements (minimum 3 services + infrastructure), we will build these 3 core domains:

1. **Student Profile Service** (Synchronous gRPC, SQL database)
2. **Notification Service** (Asynchronous via RabbitMQ, NoSQL/In-memory storage)
3. **Course Catalogue Service** (Synchronous gRPC or REST)
---

## Phase 1: Foundation & First "Vertical Slice"

**Goal:** Make a button on your frontend successfully retrieve a student profile all the way through the stack (*Frontend $\rightarrow$ API Gateway $\rightarrow$ Service $\rightarrow$ DB*).

### Step 1.1: Shared Schemas & Contracts

* **Action:** Create a central directory (`/proto`) for your Protobuf definitions.
* **Implementation:** Define `student.proto` containing your `GetStudentProfile(StudentRequest) returns (StudentResponse)` method.
* **Testing:** Generate Go code locally using `protoc` and verify no compilation errors occur.

### Step 1.2: Student Profile Service & Database

* **Action:** Build the **Student Profile Service** in Go with a PostgreSQL or SQLite database.
* **Implementation:** Expose one gRPC endpoint (`GetStudent` / `CreateStudent`).
* **Testing:**
* Write unit tests (`*_test.go`) for your database access layer.
* Use **grpcurl** or Postman to invoke your gRPC service directly on port `50051`.



### Step 1.3: Docker Compose & API Gateway

* **Action:** Create your `docker-compose.yml` including the **Student Service**, **Postgres**, and **Nginx/Traefik**.
* **Implementation:** Configure Nginx/Traefik as an API Gateway to route incoming REST calls (`/api/v1/students`) to the internal gRPC Student Service.
* **Testing:** Run `docker compose up --build`. Send an HTTP GET request via cURL or Postman to the Gateway, confirming you receive the student JSON payload back.

### Step 1.4: Minimal Frontend (Go + HTML Templates)

* **Action:** Build a simple Go web server for the UI.
* **Implementation:** Create a page (`/student?id=1`) that fetches data from the API Gateway and renders it in an HTML template.
* **Testing:** Open your browser. If you can see student profile data on the screen, **your first vertical slice is complete.**

---

## Phase 2: Asynchronous Messaging & Service #2

**Goal:** Automatically trigger an asynchronous notification whenever a student profile is updated.

### Step 2.1: RabbitMQ Infrastructure

* **Action:** Add `rabbitmq:3-management` to your `docker-compose.yml`.
* **Testing:** Spin up containers, navigate to `http://localhost:15672` (RabbitMQ Management Dashboard), and verify the broker is running.

### Step 2.2: Notification Service (Service #2)

* **Action:** Create a Go service that has no direct REST/gRPC endpoints—it only consumes messages from a RabbitMQ queue (`student_events`).
* **Implementation:** When a message (e.g., `StudentUpdated`) lands in the queue, the Notification Service logs the event or saves it to its own isolated store (e.g., Redis or MongoDB).
* **Testing:** Write a small Go test script to publish a test message directly to RabbitMQ, verifying the Notification Service consumes and logs it.

### Step 2.3: Connect Student Service to RabbitMQ

* **Action:** Extend the **Student Profile Service** to publish an event to RabbitMQ whenever a student is created or updated.
* **Testing (Integration Test):**
1. Send a `POST /api/v1/students` request from your UI or Postman.
2. Inspect **Student DB** (data is saved).
3. Inspect **RabbitMQ UI** (message was published and acknowledged).
4. Inspect **Notification Service Logs** (`docker logs notification-service`) to verify execution.



---

## Phase 3: Service #3 & Security

**Goal:** Add your 3rd service and implement JWT Authentication at the Gateway boundary.

### Step 3.1: Course Catalogue Service (Service #3)

* **Action:** Build a Go service with its own isolated database (e.g., course list, room locations).
* **Implementation:** Expose a REST or gRPC API to read and list course data.
* **Testing:** Run integration tests to verify database isolation (confirming zero direct database calls to the Student DB).

### Step 3.2: Centralized JWT Authentication at the Gateway

* **Action:** Configure Nginx/Traefik (or a small auth endpoint inside the Student Service) to issue and validate JWT tokens.
* **Implementation:** Protect sensitive routes (e.g., `POST /api/v1/courses`) at the Gateway level so they require an `Authorization: Bearer <token>` header.
* **Testing:**
* Call `POST /api/v1/courses` without a token $\rightarrow$ Expect `401 Unauthorized`.
* Generate a token via your auth route, re-send the request $\rightarrow$ Expect `200 OK`.



---

## Phase 4: Observability, Documentation & Final Check

**Goal:** Satisfy all course rubric non-functional requirements and finalize your submission.

### Step 4.1: Centralized Logging & Error Handling

* **Action:** Ensure all Go services use structured logging (e.g., Go's native `slog` or `zap`) and print logs to `stdout`/`stderr`.
* **Testing:** Run `docker compose logs -f` and confirm you can track requests flowing through the system across containers.

### Step 4.2: Documentation

* **Action:** Create a clean `README.md` containing:
1. **Architecture Diagram:** A simple visual or Mermaid.js diagram illustrating service communication.
2. **Run Instructions:** A single command to start the entire environment: `docker compose up --build`.
3. **API Specs:** A cURL collection or an `api.http` file (VS Code REST Client format) for easy testing of your endpoints.



---

## Incremental Testing & Verification Matrix

| Step | Component | Test Tool / Approach | Success Criteria |
| --- | --- | --- | --- |
| **1** | Protobuf Schemas | `protoc` CLI | Code compiles clean in Go without missing dependencies. |
| **2** | Isolated Service (gRPC) | `grpcurl` / Go Unit Tests | Service responds to payload on local port `50051`. |
| **3** | Docker Compose & Gateway | Postman / `curl` | REST calls to port `80/443` correctly proxy to internal services. |
| **4** | RabbitMQ Messaging | RabbitMQ Management Web UI | Messages appear in queue and clear out upon consumption. |
| **5** | JWT Authentication | Postman / Insomnia | Protected endpoints reject unauthenticated calls with HTTP 401. |