Her er den opdaterede og tilpassede udrulningsplan for din **INFS605 platform**, udvidet til at omfatte alle **6 mikrotjenester** og deres specifikke databaseteknologier (SQL vs. NoSQL), som vi har defineret.

Planen er fortsat bygget op omkring en **Vertical Slice strategi**: Vi færdiggør hele værdikæden for én funktion ad gangen (Frontend $\rightarrow$ Gateway $\rightarrow$ Service $\rightarrow$ DB/RabbitMQ), før vi bygger videre.

---

## Mål-Arkitektur & Tjenester

1. **Authentication Service** (gRPC, SQL DB – Identitet & Tokens)
2. **Student Profile Service** (gRPC, SQL DB – Brugerstamdata)
3. **Course Catalogue Service** (gRPC, SQL DB – Kurser, Fag & Tilmeldinger/Enrollments)
4. **Course Content Service** (gRPC, **NoSQL/MongoDB** – Dynamiske lektionsblokke)
5. **Assignment/Grading Service** (gRPC, SQL DB – Afleveringer & Karakterer)
6. **Notification Service** (RabbitMQ Consumer + gRPC, **NoSQL/In-App DB** – In-App Notifikationer)

---

## [x] Phase 1: Foundation & First "Vertical Slice" (Student Profile)

**Mål:** En knap i din frontend henter en studerendes profil hele vejen igennem stakken (*Frontend $\rightarrow$ API Gateway $\rightarrow$ Profile Service $\rightarrow$ SQLite DB*).

### [x] Step 1.1: Shared Schemas & Contracts (gRPC)

* **Handling:** Opret mappen `/proto` med `.proto` kontrakter.
* **Implementering:** Definer `student.proto` med `GetProfile` og `CreateStudent`.
* **Test:** Generer Go-kode med `protoc` uden kompileringsfejl.

### [x] Step 1.2: Student Profile Service & Database (GORM + SQLite)

* **Handling:** Byg **Student Profile Service** i Go.
* **Implementering:** Implementer GORM repository og in-memory unit tests (`*_test.go`).
* **Test:** Kør `go test ./...` og verificer gRPC-kald direkte på port `50051`.

### [x] Step 1.3: API Gateway & Docker Compose Integration

* **Handling:** Tilføj Nginx/Traefik som API Gateway foran `profile-service`.
* **Implementering:** Rute indkommende HTTP REST-kald (`/api/v1/students`) videre til den interne gRPC profile-service.
* **Test:** Kør `docker compose up --build`. Send et HTTP GET/POST-kald via cURL/Postman til Gateway og modtag JSON-svar.

### [x] Step 1.4: Minimal Frontend Integration

* **Handling:** Opret en simpel Go-baseret webfrontend (BFF / HTML templates).
* **Implementering:** Opret siden `/profile?id=xxx`, der henter data fra Gateway.
* **Test:** Åbn browseren. Hvis du kan se profiloplysningerne på skærmen, er første vertical slice komplet.

---

## [x] Phase 2: Asynchronous Messaging & In-App Notifications (Service #2)

**Mål:** Når en karakter gives eller en student oprettes, oprettes der automatisk en in-app notifikation via RabbitMQ.

### [x] Step 2.1: RabbitMQ Infrastructure

* **Handling:** Tilføj `rabbitmq:3-management` til `docker-compose.yml`.
* **Test:** Tilgå `http://localhost:15672` og bekræft, at brokeren kører.

### [x] Step 2.2: Notification Service (Service #2 - NoSQL / In-App Feed)

* **Handling:** Opret `notification-service` med en NoSQL/In-App DB (MongoDB/SQLite) til notifikationshistorik.
* **Implementering:**
1. Opret en RabbitMQ consumer, der lytter på events (f.eks. `grade.published`, `student.created`).
2. Tilføj en gRPC-endpoint (`GetUserNotifications`, `MarkAsRead`), så frontenden kan vise ulæste notifikationer ved login.
3. Vis Notifikationer i frontend via `/notifications` side.

* **Test:** Send en test-event til RabbitMQ og verificer via gRPC, at notifikationen kan hentes frem.

---

## Phase 3: Core Domain Expansion

### [ ] Step 3.2: Course Content Service (Service #4 - NoSQL / CloverDB)

> **Note:** Siden du har valgt **CloverDB** (embedded NoSQL), kører databasen lokalt i din Go-proces og gemmer i `./data/nosql` i stedet for en MongoDB-container.

#### **1. DB & Storage Setup**

* [ ] **CloverDB Initialization:** Opret og initialiser CloverDB i `internal/repository/clover.go` og opret collectionen `"modules"`.
* [ ] **Domain Models:** Definer `domain.Module` og `domain.File` med tilhørende `json:"_id"` og `json:"..."` tags.
* [ ] **Seed Data:** Opret en `SeedCloverData(db)` funktion, der indsætter test-moduler og fil-arrays, hvis collectionen er tom.
* [ ] **Docker Volume:** Verificer at `./course-catalogue-service/data:/app/data` er mounted i `docker-compose.yml`, så NoSQL-dataene overlever container-genstarter.

#### **2. Service & Repository-lag**

* [ ] **Clover Repository Implementering:**
* [ ] `GetModulesByCourseID(ctx, courseID)` $\rightarrow$ Udfører `query.NewQuery("modules").Where(query.Field("course_id").IsEq(courseID))` og unmarshaller til `[]*domain.Module`.
* [ ] `SaveModule(ctx, module)` $\rightarrow$ Bruger `document.NewDocumentOf(module)` og indsætter/opdaterer i CloverDB.


* [ ] **Service Layer Business Logic:** Opret `ContentService` der binder repository sammen med eventuelle valideringer (fx tjek om kurset findes via gRPC-kald til SQL-delen).

#### **3. gRPC Server Setup**

* [ ] **Proto Specifikation:** Definer `content.proto` med beskeder som `Module`, `File`, `GetCourseContentRequest` og `GetCourseContentResponse`.
* [ ] **gRPC Handler Implementation:** Implementer `GetCourseContent` handleren, som kalder `ContentService` og mapper `domain.Module` og `domain.File` over til Proto structs.
* [ ] **Server Registration:** Registrer `ContentServer` i din `main.go`.

#### **4. Frontend / BFF Integration & Test**

* [ ] **BFF Client:** Tilføj `ContentClient` til din Go Frontend BFF og konfigurer `COURSE_CONTENT_SERVICE_ADDR`.
* [ ] **HTTP Handler:** Opret `/courses/{id}/content` endpoint i frontenden.
* [ ] **Template / UI View:** Render modulernes titler, beskedtekster og fil-links (`<a href="...">filename.pdf</a>`) i HTML-skabelonen.
* [ ] **Verificering:** Åbn et kursus i frontenden og bekræft, at moduler og filer hentes ud i ét enkelt læsekald via gRPC fra CloverDB.

---

### [ ] Step 3.3: Assignment & Grading Service (Service #5 - SQL)

#### **1. Database Setup**

* [ ] **GORM Schema / Migrations:** Opret tabeller for `assignments` (`id`, `course_id`, `title`, `due_date`) og `submissions` (`id`, `assignment_id`, `student_id`, `grade`, `submitted_at`).
* [ ] **Database Seeding:** Seed et par test-opgaver til eksisterende kurser.

#### **2. Service & Repository-lag**

* [ ] **Grading Repository:** Implementer `CreateSubmission` og `UpdateGrade`.
* [ ] **Service Layer & Event Triggering:**
* [ ] I `GradeSubmission()` gemmes karakteren i SQL-databasen.
* [ ] Så snart DB-opdateringen lykkes, udgives en `grade.published` event på RabbitMQ med `{student_id, course_id, grade}`.



#### **3. gRPC Server Setup**

* [ ] **Proto Specifikation:** Definer `grading.proto` med `SubmitAssignment` og `GradeSubmission` RPC-kald.
* [ ] **gRPC Handler:** Opret server-implementering og kobl den på gRPC-porten (f.eks. `:50054`).

#### **4. Frontend / BFF Integration & Test**

* [ ] **UI Form:** Opret en simpel side/form i frontenden, hvor en underviser kan vælge en studerende og indtaste en karakter.
* [ ] **Verificering af Event Flow:**
1. Underviser trykker "Gem Karakter" $\rightarrow$ Frontend kalder `/api/grades`.
2. Frontend kalder `grading-service` via gRPC.
3. `grading-service` gemmer i sin DB og sender `grade.published` event på RabbitMQ.
4. `notification-service` opfanger eventet og gemmer en notifikation ("Du har modtaget karakteren A i CS101").
5. Den studerende logger ind og ser notifikationen på sit dashboard.

---

### [ ] Step 3.4: Event Publishing fra Services

#### **1. Database & Domain Event Setup**

* [ ] **Domain Events Definition:** Opret en fælles struct/proto for events (f.eks. `StudentCreatedEvent`, `EnrollmentCreatedEvent`).
* [ ] **RabbitMQ Producer Wrapper:** Byg en genanvendelig `publisher.go` i din service, som håndterer forbindelse, kanaler og genforbindelse til RabbitMQ.
* [ ] **JSON/Protobuf Serialization:** Konverter dit domæne-event til JSON eller Protobuf før det udgives på switchen (exchange).

#### **2. Service-lag Integration**

* [ ] **Publish ved DB mutation:** Kald `publisher.Publish("student.created", payload)` lige efter en succesfuld SQL-transaktion (f.eks. i `CreateProfile`).
* [ ] **Error Handling / Fallback:** Sørg for at logge en klar fejl, hvis DB-ændringen lykkedes, men RabbitMQ-kaldet fejler (eller implementer Outbox pattern, hvis du vil være ekstra grundig).

#### **3. Server & Consumer Setup**

* [ ] **Notification Consumer Setup:** I `notification-service`, lyt på kørselstidspunktet på RabbitMQ-køen `notification-queue` bundet til de relevante routing keys (`*.created`, `*.published`).
* [ ] **Consumer Handler:** Opret en ny notifikation i Notification DB, når en besked modtages.

#### **4. Frontend Integration & Test**

* [ ] **UI Trigger:** Opret en ny studerende eller tilmeld et kursus via Frontend.
* [ ] **RabbitMQ Dashboard Check:** Tjek http://localhost:15672 og verificer, at beskeden tæller op under `Publish` rate på køen.
* [ ] **Notification Badge i Frontend:** Åbn notifications-siden i frontenden og verificer, at den nyligt oprettede notifikation vises for brugeren.

---

### [ ] Step 3.5: Authentication Service (Service #6 - Security Boundary)

#### **1. Database Setup**

* [ ] **Auth Schema / Migrations:** Opret `user_accounts` tabel med `id`, `email`, `password_hash`, og `role` (`student`, `teacher`, `admin`).
* [ ] **Crypto Setup:** Implementer `bcrypt` til sikker hashing og sammenligning af adgangskoder.

#### **2. Service & JWT Implementation**

* [ ] **JWT Generator:** Opret en hjælpefunktion til at udstede og signere JWT tokens (indeholdende `user_id`, `email`, `role` og `exp`).
* [ ] **Auth Service Methods:** Implementer `Register` og `Login` metoder i servicelaget.

#### **3. gRPC Server & Gateway Interceptors**

* [ ] **Proto Specifikation:** Definer `auth.proto` med `Login` og `ValidateToken` RPCs.
* [ ] **gRPC Auth Interceptor:** Opret en gRPC Interceptor på tværs af microservices, der læser JWT tokenet ud af gRPC Context Metadata (`authorization: bearer <token>`) og verificerer signatur.

#### **4. Frontend / BFF Integration & Test**

* [ ] **Session / Cookie Handling:** Når brugeren logger ind på `/login` i Go Frontenden, kaldes `auth-service`. Ved succes gemmes JWT-tokenet i en sikker, `HttpOnly` cookie i browseren.
* [ ] **BFF Middleware (`h.Authenticate`):** Opdater din `Authenticate` middleware i Go Frontenden til at læse cookien og vedhæfte JWT-tokenet som gRPC Metadata på *alle* udgående microservice-kald.
* [ ] **Test Scenarier:**
* [ ] Test adgang til `/dashboard` uden login-cookie $\rightarrow$ Omdirigeres til `/login` (eller returnerer HTTP `401`).
* [ ] Test adgang med gyldigt login $\rightarrow$ gRPC-kaldene modtager `user_id` direkte fra gRPC metadata og returnerer den korrekte brugers data (`200 OK`).

---

## [ ] Phase 4: Observability, Documentation & Final Check

**Mål:** Opfylod alle ikke-funktionelle krav i fagets bedømmelseskriterier.

### [ ] Step 4.1: Centralized Logging & Error Handling

* **Handling:** Sørg for at alle 6 Go-services bruger struktureret logging (`slog` eller `zap`) til `stdout`/`stderr`.
* **Test:** Kør `docker compose logs -f` og følg en anmodnings vej gennem Gateway og services.

### [ ] Step 4.2: Documentation

* **Handling:** Opret en grundig `README.md` med:
1. **Arkitekturdiagram:** Visualisering af de 6 services, RabbitMQ, Gateway, SQL og NoSQL databaser.
2. **Run Instructions:** `docker compose up --build`.
3. **API Collection:** En `.http` fil eller Postman collection til test af alle væsentlige flow.

---

## [ ] Phase 5: High Availability & Instance Scaling

**Mål:** Demonstrer horisontal skalering og belastningsfordeling på tværs af servicerne.

### [ ] Step 5.1: Multi-Instance Docker Compose Configuration

* **Handling:** Fjern specifikke port-bindings på interne mikrotjeneste-containere i `docker-compose.yml`.
* **Eksekvering:**
```bash
docker compose up -d --scale profile-service=3 --scale notification-service=2

```

* **Test:** Kør `docker compose ps` og bekræft, at alle instanser kører på det fælles Docker-netværk.

### [ ] Step 5.2: Gateway Round-Robin Load Balancing

* **Handling:** Konfigurer Nginx som load balancer for gRPC og REST backends.
* **Test:** Log `os.Hostname()` i Go-servicerne, send 6 opkald igennem Gatewayen, og verificer i loggen, at forespørgslerne fordeles ligeligt på tværs af container-ID'erne.

### [ ] Step 5.3: Asynchronous Competing Consumers (RabbitMQ)

* **Handling:** Sørg for, at alle skalerede instanser af `notification-service` lytter på den **samme RabbitMQ-kø**.
* **Test:** Send 10 karakter-events hurtigt efter hinanden. Bekræft i logs og i RabbitMQ Dashboard, at hver event kun behandles **én gang** af én af instanserne (round-robin).

### [ ] Step 5.4: Database Connection Pooling & Statelessness Audit

* **Handling:**
* Begræns database-forbindelser i Go drivers (`db.SetMaxOpenConns(10)`).
* Sikr at alle services validerer adgang via tilstandsløse (stateless) JWT-signaturer.


* **Test:** Kør en stresstest med `hey` eller `ab`:
```bash
hey -n 200 -c 20 http://localhost/api/v1/courses

```