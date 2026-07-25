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

## [ ] Phase 2: Asynchronous Messaging & In-App Notifications (Service #2)

**Mål:** Når en karakter gives eller en student oprettes, oprettes der automatisk en in-app notifikation via RabbitMQ.

### [x] Step 2.1: RabbitMQ Infrastructure

* **Handling:** Tilføj `rabbitmq:3-management` til `docker-compose.yml`.
* **Test:** Tilgå `http://localhost:15672` og bekræft, at brokeren kører.

### [ ] Step 2.2: Notification Service (Service #2 - NoSQL / In-App Feed)

* **Handling:** Opret `notification-service` med en NoSQL/In-App DB (MongoDB/SQLite) til notifikationshistorik.
* **Implementering:**
1. Opret en RabbitMQ consumer, der lytter på events (f.eks. `grade.published`, `student.created`).
2. Tilføj en gRPC-endpoint (`GetUserNotifications`, `MarkAsRead`), så frontenden kan vise ulæste notifikationer ved login.

* **Test:** Send en test-event til RabbitMQ og verificer via gRPC, at notifikationen kan hentes frem.

### [ ] Step 2.3: Event Publishing fra Services

* **Handling:** Udvid services til at udgive events på RabbitMQ ved ændringer.
* **Integrationstest:**
1. Udfør en handling i systemet (f.eks. opret student).
2. Tjek **RabbitMQ UI** (beskeden blev sendt).
3. Tjek **Notification Service Logs & DB** (notifikationen er gemt og klar til at blive vist ved login).



---

## [ ] Phase 3: Core Domain Expansion (Services #3, #4, #5 & #6)

**Mål:** Tilføj resten af universitetets domæneservices og introducer NoSQL til kursusindhold.

### [ ] Step 3.1: Course Catalogue Service (Service #3 - SQL)

* **Handling:** Byg `catalogue-service` med ansvar for fagkataloget og **Enrollments** (M2M relation mellem `student_id` og `course_id`).
* **Implementering:** Eksponer gRPC endpoints: `GetCourse`, `ListCourses`, `EnrollStudent`.
* **Test:** Verificer, at tilmeldinger gemmes korrekt uden direkte databaselink til Profile DB.

### [ ] Step 3.2: Course Content Service (Service #4 - NoSQL / MongoDB)

* **Handling:** Byg `content-service` med **MongoDB** til opbevaring af fleksible lektionsblokke (tekst, video-links, PDF'er).
* **Implementering:** Gem hele moduler og deres indholdsblokke som dokumenter i MongoDB.
* **Test:** Hent et modul ud på ét enkelt databaseopslag og verificer, at dokument-strukturen returneres korrekt via gRPC.

### [ ] Step 3.3: Assignment & Grading Service (Service #5 - SQL)

* **Handling:** Byg `grading-service` til håndtering af opgaveafleveringer, deadlines og karakterer.
* **Implementering:** Når en karakter gemmes via `GradeSubmission`, publiceres en `grade.published` event til RabbitMQ.
* **Test:** Giv en karakter $\rightarrow$ Verificer at `notification-service` modtager eventen og opretter en in-app notifikation.

### [ ] Step 3.4: Authentication Service (Service #6 - Security Boundary)

* **Handling:** Byg `auth-service` med ansvar for `UserAccount` (Email, PasswordHash) og JWT-udstedelse.
* **Implementering:**
1. Konfigurer Gateway/BFF til at kræve `Authorization: Bearer <JWT>` på beskyttede ruter.
2. Sørg for, at `user_id` fra JWT videreføres i gRPC Metadata til de underliggende services.


* **Test:** Test adgang uden token ($\rightarrow 401$) vs. med gyldigt JWT ($\rightarrow 200$).

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