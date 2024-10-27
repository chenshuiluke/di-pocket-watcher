import { Button, ButtonProps } from "@mantine/core";
import { GoogleIcon } from "../icons/GoogleIcon";
import { useState } from "react";

export const Login = () => {
    const [isLoading, setIsLoading] = useState(false);

    const GoogleButton = (props: ButtonProps & React.ComponentPropsWithoutRef<'button'>) => {
        return <Button leftSection={<GoogleIcon />} variant="default" {...props} />;
    }

    const loginWithGoogle = () => {
        setIsLoading(true);
        
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
                    // You might want to trigger a refresh or redirect here
                    window.location.reload();
                }
            }
        };

        window.addEventListener('message', handleMessage);

        // Cleanup
        const checkPopup = setInterval(() => {
            if (!popup || popup.closed) {
                clearInterval(checkPopup);
                window.removeEventListener('message', handleMessage);
                setIsLoading(false);
            }
        }, 1000);
    }

    return (
        <GoogleButton 
            onClick={loginWithGoogle} 
            loading={isLoading}
            disabled={isLoading}
        >
            Continue with Google
        </GoogleButton>
    );
}