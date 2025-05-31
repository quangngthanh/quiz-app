Dưới đây là bản hướng dẫn README.md tổng hợp, rõ ràng cho dự án quiz-app, chia thành hai phần: **Backend** và **Frontend**. Nội dung bao gồm công nghệ sử dụng, lệnh khởi động, migrate, và các thông tin cần thiết khác.

---

# Quiz App

Ứng dụng Quiz Realtime gồm hai phần: **Backend (Golang)** và **Frontend (React + Vite)**.

---

## 1. Backend

### Công nghệ sử dụng
- **Ngôn ngữ:** Go (Golang)
- **Framework:** Gin
- **ORM:** GORM
- **Database:** PostgreSQL
- **Cache:** Redis
- **WebSocket:** Gin + Gorilla WebSocket
- **Migration:** SQL scripts (migrations folder)
- **Container:** Docker, Docker Compose

### Cấu trúc thư mục
```
backend/
  ├── cmd/                # Entry point (main.go)
  ├── config/             # Cấu hình (YAML)
  ├── internal/
  │   ├── config/         # Đọc config
  │   ├── handler/        # HTTP/WebSocket handlers
  │   ├── model/          # Định nghĩa models
  │   ├── repository/     # Truy cập DB, Redis
  │   └── service/        # Business logic
  ├── migrations/         # SQL migration scripts
  ├── Dockerfile
  ├── go.mod, go.sum
```

### Các lệnh khởi động & migrate

#### 1. Cài đặt dependencies
```bash
cd backend
go mod tidy
```

#### 2. Khởi động server ()
```bash
# Chạy trực tiếp (yêu cầu đã có postgres và redis như trong docker compose file: `docker compose up postgres redis`)
go run cmd/server/main.go

# Hoặc dùng Docker in production
docker-compose up --build
```


#### 3. Các API chính
- `POST   /api/quiz`                : Tạo quiz mới
- `POST   /api/quiz/:quizID/join`   : Tham gia quiz
- `POST   /api/quiz/:quizID/answer` : Gửi đáp án
- `GET    /api/quiz/:quizID`        : Lấy thông tin quiz
- `GET    /api/quiz/:quizID/leaderboard` : Lấy bảng xếp hạng

#### 4. WebSocket
- `GET /ws/quiz/:quizID/leaderboard` : Nhận realtime leaderboard

---

## 2. Frontend

### Công nghệ sử dụng
- **React** (TypeScript)
- **Vite** (build tool)
- **WebSocket** (realtime)
- **CSS modules**
- **TailwindCSS**

### Cấu trúc thư mục
```
frontend/
  ├── public/
  ├── src/
  │   ├── components/      # Các component chính (Leaderboard, QuizCreator, QuizParticipant)
  │   ├── hooks/           # Custom hooks (useWebSocket)
  │   ├── utils/           # Hàm tiện ích (api.ts)
  │   ├── types/           # Định nghĩa type
  │   └── App.tsx, main.tsx
  ├── package.json
  ├── tsconfig.json
  └── vite.config.ts
```

### Các lệnh khởi động

#### 1. Cài đặt dependencies
```bash
cd frontend
npm install
```

#### 2. Khởi động dev server
```bash
npm run dev
```
- Ứng dụng sẽ chạy tại: [http://localhost:5173](http://localhost:5173)

#### 3. Build production
```bash
npm run build
```
---

## 3. Docker Compose

- Có thể sử dụng `docker-compose.yaml` ở thư mục gốc để khởi động toàn bộ hệ thống (backend, frontend, database, redis).

```bash
docker-compose up --build
```

---

## 4. Thông tin bổ sung

- **Redis:** Cần chạy Redis để cache leaderboard.
- **Cổng mặc định:**
  - Backend: `:8088`
  - Frontend: `:5173`
  - PostgreSQL: `:5432`
  - Redis: `:6379`
- **Tài khoản mặc định:** Không có, user tự nhập username khi join quiz.

---

## 5. Liên hệ & Đóng góp

- Repo: [https://github.com/quangngthanh/quiz-app](https://github.com/quangngthanh/quiz-app)
- Đóng góp: Pull request, issue, hoặc liên hệ trực tiếp qua Github.

