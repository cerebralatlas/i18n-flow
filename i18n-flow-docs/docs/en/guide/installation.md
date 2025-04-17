# Installation

## Backend Setup

1. Clone the repository:

```bash
git clone https://github.com/ilukemagic/i18n-flow.git
cd i18n-flow
```

2. Install dependencies:

```bash
cd admin-backend
go mod tidy
```

3. Configure database:

```bash
# Edit config/config.yaml with your database credentials
```

4. Start the server:

```bash
go run main.go
```

## Frontend Setup

1. Navigate to frontend directory:

```bash
cd admin-frontend
```

2. Install dependencies:

```bash
pnpm install
```

3. Start development server:

```bash
npm run dev
```
