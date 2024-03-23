"use client"

import { ChangeEvent, SyntheticEvent, useState } from 'react';
import Image from "next/image"
import { Button, Input, Checkbox } from "@nextui-org/react";
import { useRouter } from "next/navigation";
import { ApiPostRequest } from "@/requests/requests";

enum PageState {
    PS_LOGIN,
    PS_REGISTER
}


export default function Login() {
    let router = useRouter()
    let [currentState, setCurrentState] = useState<PageState>(PageState.PS_LOGIN)

    let [loginUsername, setLoginUsername] = useState<string>("")
    let [loginPassword, setLoginPassword] = useState<string>("")
    let [registerUsername, setRegisterUsername] = useState<string>("")
    let [registerEmail, setRegisterEmail] = useState<string>("")
    let [registerPassword, setRegisterPassword] = useState<string>("")
    let [registerRepeatPassword, setRegisterRepeatPassword] = useState<string>("")

    let [isLoading, setIsLoading] = useState<boolean>(false)
    let [registerError, setRegisterError] = useState<string>("")
    let [loginError, setLoginError] = useState<string>("")

    let sendLoginRequest = () => {
        setLoginError("")
        setIsLoading(true)
        ApiPostRequest({
            body: {
                username: loginUsername,
                password: loginPassword
            },
            route: "/auth/login",
            responseHandler: () => {
                router.push("/")
            },
            errorHandler: (err) => {
                setIsLoading(false)
                switch (err) {
                    case "INVALID_AUTH":
                        setLoginError("Invalid credentials")
                }
            }
        })
    }

    let sendRegisterRequest = () => {
        setRegisterError("")
        if (registerPassword != registerRepeatPassword) {
            setRegisterError("Passwords must match")
            return
        }
        setIsLoading(false)
        ApiPostRequest({
            body: {
                username: registerUsername,
                email: registerEmail,
                password: registerPassword
            },
            route: "/auth/register",
            responseHandler: (_) => {
                setIsLoading(false)
                alert("User created succesfully!")
            },
            errorHandler: (err) => {
                setIsLoading(false)
                switch (err) {
                    case "USERNAME_INVALID":
                        setRegisterError("Invalid username")
                    case "PASSWORD_WEAK":
                        setRegisterError("Password is too weak")
                    case "EMAIL_INVALID":
                        setRegisterError("Email is invalid")
                    case "DUPLICATED_RESOURCE":
                        setRegisterError("User already exists")
                }
            }
        })
    }

    return (
        <div className='flex flex-row'>
            <div className='bg-background-300 relative w-full bg-[url(/intrude-on-the-mind.jpg)] bg-cover'>
            </div>
            <div className='flex h-screen justify-center'>
                <div className='my-auto px-8 flex flex-col items-center'>
                    <Image src="/logo.png" alt="" width={100} height={100}></Image>
                    <h2 className="text-3xl mt-4 text-primary-50 font-sans">Tolarian Archives</h2>
                    {
                        (currentState == PageState.PS_LOGIN) && (
                            <div className="flex flex-col items-center justify-center px-6 py-8 my-8 lg:py-0">
                                <div className="p-6 space-y-4 md:space-y-6 sm:p-8">
                                    <h1 className="text-xl font-bold leading-tight tracking-tight text-white md:text-2xl">
                                        Sign in to your account
                                    </h1>
                                    <form onSubmit={sendLoginRequest} className="space-y-4 md:space-y-6 max-w-min min-w-96" action="#">
                                        <Input
                                            id="username"
                                            type="username"
                                            label="Username"
                                            placeholder="Username"
                                            onValueChange={(value) => setLoginUsername(value)}
                                            isDisabled={isLoading}
                                            required />
                                        <Input
                                            id="password"
                                            type="password"
                                            label="Password"
                                            placeholder="**********"
                                            onValueChange={(value) => setLoginPassword(value)}
                                            isDisabled={isLoading}
                                            required />
                                        <div className="text-sm font-light text-red-400">{loginError}</div>
                                        <Button color="success" isLoading={isLoading} onClick={sendLoginRequest} fullWidth>Sign in</Button>
                                        <div className="flex flex-col gap-4 text-sm font-light text-gray-400 justify-center items-center">
                                            <div className="flex flex-row">
                                                {"New user?"}
                                                <a href="#" className="font-medium ml-1 text-secondary-600 hover:underline"
                                                    onClick={() => setCurrentState(PageState.PS_REGISTER)}>
                                                    Sign up
                                                </a>
                                            </div>
                                            {/* <a href="#" className="text-sm font-medium text-secondary-600 hover:underline">Forgot password?</a> */}
                                        </div>
                                    </form>
                                </div >
                            </div >
                        )
                    }
                    {
                        (currentState == PageState.PS_REGISTER) && (
                            <div className="flex flex-col items-center justify-center px-6 py-8 my-8 mx-auto lg:py-0">
                                <div className="p-6 space-y-4 md:space-y-6 sm:p-8">
                                    <h1 className="text-xl font-bold leading-tight tracking-tight text-white md:text-2xl">
                                        Sign up
                                    </h1>
                                    <form onSubmit={sendRegisterRequest} className="space-y-4 md:space-y-6 max-w-min min-w-96" action="#">
                                        <Input
                                            id="username"
                                            label="Username"
                                            placeholder="Username"
                                            onValueChange={(value) => setRegisterUsername(value)}
                                            isDisabled={isLoading}
                                            required />
                                        <Input
                                            type="email"
                                            id="email"
                                            label="Email"
                                            placeholder="Email"
                                            onValueChange={(value) => setRegisterEmail(value)}
                                            isDisabled={isLoading}
                                            required />
                                        <Input
                                            id="password"
                                            label="Password"
                                            type="password"
                                            placeholder="**********"
                                            onValueChange={(value) => setRegisterPassword(value)}
                                            isDisabled={isLoading}
                                            required />
                                        <Input
                                            id="repeatPassword"
                                            label="Repeat password"
                                            type="password"
                                            placeholder="**********"
                                            onValueChange={(value) => setRegisterRepeatPassword(value)}
                                            isDisabled={isLoading}
                                            required />
                                        <p className="text-sm font-light text-red-400">{registerError}</p>
                                        <Button color="success" isLoading={isLoading} onClick={sendRegisterRequest} fullWidth>Sign up</Button>
                                        <a href="#" className="text-sm justify-center flex items-center font-medium ml-1 text-secondary-600 hover:underline"
                                            onClick={() => setCurrentState(PageState.PS_LOGIN)}>
                                            Back to sign in
                                        </a>
                                    </form>
                                </div >
                            </div >
                        )
                    }
                </div>
            </div>
        </div>
    );
}