import axios from "axios";

const BASE_URL = "http://localhost:8080/api";

export const apiClient = axios.create({
  baseURL: BASE_URL,
});

// Request interceptor to add auth token
apiClient.interceptors.request.use((config) => {
  const token = localStorage.getItem("jwt_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor for error handling
apiClient.interceptors.response.use(
  (response) => response,
  async (error) => {
    // Handle 401 Unauthorized errors (expired token, etc.)
    if (error.response?.status === 401) {
      localStorage.removeItem("jwt_token");
      window.location.href = "/login";
    }
    return Promise.reject(error);
  }
);

// Type for API error responses
export type ApiError = {
  message: string;
  code?: string;
  details?: Record<string, unknown>;
};
