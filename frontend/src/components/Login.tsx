import { ButtonProps, Button } from "@mantine/core";
import { GoogleIcon } from "../icons/GoogleIcon";

export const Login = () => {
    const GoogleButton = (props: ButtonProps & React.ComponentPropsWithoutRef<'button'>) => {
        return <Button leftSection={<GoogleIcon />} variant="default" {...props} />;
    }

    const loginWithGoogle = () => {
        window.location.href = "http://localhost:8080/auth/google/login"
    }

    return <GoogleButton onClick={loginWithGoogle}>Continue with Google</GoogleButton>

}