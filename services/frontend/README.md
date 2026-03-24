# 🖥️ TruckGuard Frontend

### 1. What is it?

The **TruckGuard Frontend** is a modern, responsive web application built with **SvelteKit** and **Tailwind CSS**. It serves as the primary interface for operators, administrators, and accountants to manage customs weighing operations.

### 2. Purpose & How it Works

It provides a centralized dashboard for system management and real-time monitoring:

- **Live Dashboard**: Real-time overview of active permits and vehicle movements.
- **Plate Events**: Detailed log of all recognition events with images and metadata.
- **Permit Management**: Workflow for creating, editing, and closing vehicle passes.
- **Configuration**: Management of cameras, scales, and system-wide settings.
- **Role-Based Access**: Specialized interfaces for Operators, Controllers, and Accountants.

### 3. Tech Stack

- **Framework**: [SvelteKit](https://kit.svelte.dev/)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/)
- **UI Components**: [Shadcn-Svelte](https://www.shadcn-svelte.com/)
- **State Management**: Svelte Runes (v5)
- **Language**: [TypeScript](https://www.typescriptlang.org/)

### 4. Getting Started

#### **Prerequisites**

- Node.js (v20+)
- Yarn

#### **Run Commands**

1.  **Install dependencies:**
    ```bash
    yarn install
    ```
2.  **Start the development server:**
    ```bash
    yarn dev
    ```
3.  **Build for production:**
    ```bash
    yarn build
    ```

### 5. Configuration (Environment Variables)

Ensure you have a `.env` file with the following:

```env
PUBLIC_CORE_API_URL=http://localhost:8080
PUBLIC_AUTH_API_URL=http://localhost:8081
```
