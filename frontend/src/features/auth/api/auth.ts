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

const loginWithGoogle = () => {
  return new Promise((resolve, reject) => {
      // Open the Google login in a popup
      const width = 500;
      const height = 600;
      const left = window.screenX + (window.outerWidth - width) / 2;
      const top = window.screenY + (window.outerHeight - height) / 2.5;

      const popup = window.open(
          'http://localhost:8080/api/auth/google/login',
          'Google Login',
          `width=${width},height=${height},left=${left},top=${top}`
      );

      // Setup message listener for the token
      const handleMessage = (event: MessageEvent) => {
          if (event.origin === 'http://localhost:8080') {
              if (event.data.token) {
                  // Store the token
                  localStorage.setItem('jwt_token', event.data.token);
                  // Close the popup
                  popup?.close();
                  resolve(event.data.token)
              }
              reject()
          }
      };

      window.addEventListener('message', handleMessage);

      // Cleanup
      const checkPopup = setInterval(() => {
          if (!popup || popup.closed) {
              clearInterval(checkPopup);
              window.removeEventListener('message', handleMessage);
              
          }
      }, 1000);
  })

}

export const useGoogleLogin = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: async () => {
      await loginWithGoogle()
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: authKeys.current() });
    },
  });
};
