import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { apiClient } from "@/lib/api-client";

export type User = {
  id: number;
  email: string;
};

// Query keys
export const authKeys = {
  all: ["auth"] as const,
  current: () => [...authKeys.all, "current"] as const,
  user: (id: number) => [...authKeys.all, "user", id] as const,
};

// Get current user
export const useCurrentUser = () => {
  return useQuery({
    queryKey: authKeys.current(),
    queryFn: async () => {
      const { data } = await apiClient.get<{ user: User }>("/auth");
      return data.user;
    },
  });
};

// Logout mutation
export const useLogout = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      localStorage.removeItem("jwt_token");
      await queryClient.invalidateQueries({ queryKey: authKeys.all });
    },
  });
};

// Google OAuth login (handled via popup, as previously implemented)
export const useGoogleLogin = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      // Implementation remains in the Login component
      // This is just a placeholder for additional login-related logic
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: authKeys.current() });
    },
  });
};
