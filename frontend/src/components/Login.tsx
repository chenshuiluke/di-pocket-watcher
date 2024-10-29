import { Button, ButtonProps } from "@mantine/core";
import { GoogleIcon } from "../icons/GoogleIcon";
import { useGoogleLogin } from "@/features/auth/api/auth";



export const Login = () => {
    const loginMutation = useGoogleLogin()

    const GoogleButton = (props: ButtonProps & React.ComponentPropsWithoutRef<'button'>) => {
        return <Button leftSection={<GoogleIcon />} variant="default" {...props} />;
    }

    

    return (
        <GoogleButton
            onClick={() => loginMutation.mutate()}
            loading={loginMutation.isPending}
            disabled={loginMutation.isPending}
        >
            Continue with Google
        </GoogleButton>
    );
}