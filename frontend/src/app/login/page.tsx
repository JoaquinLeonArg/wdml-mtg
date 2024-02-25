"use client"

import { ChangeEvent, useState } from 'react';
import Image from "next/image"
import { Button } from "@/components/buttons";
import { TextFieldWithLabel } from "@/components/field";
import { Checkbox } from "@/components/checkbox";
import { useRouter } from "next/navigation";
import { ApiPostRequest } from "@/requests/requests";

enum PageState {
    PS_LOGIN,
    PS_REGISTER
}

type RegisterForm = {
    username: string
    email: string
    password: string
    repeatPassword: string
}

type LoginForm = {
    username: string
    password: string
}

export default function Login() {
    let router = useRouter()
    let [currentState, setCurrentState] = useState<PageState>(PageState.PS_LOGIN)
    let [isLoading, setIsLoading] = useState<boolean>(false)
    let [registerForm, setRegisterForm] = useState<RegisterForm>({ username: "", email: "", password: "", repeatPassword: "" })
    let [registerError, setRegisterError] = useState<string>("")
    let [loginForm, setLoginForm] = useState<LoginForm>({ username: "", password: "" })
    let [loginError, setLoginError] = useState<string>("")

    let sendLoginRequest = () => {
        setLoginError("")
        ApiPostRequest({
            body: {
                username: loginForm.username,
                password: loginForm.password
            },
            route: "/auth/login",
            responseHandler: (res) => {
                router.push("/")
            },
            errorHandler: (err) => {
                switch (err) {
                    case "INVALID_AUTH":
                        setLoginError("Invalid credentials")
                }
            }
        })
    }

    let sendRegisterRequest = () => {
        setRegisterError("")
        if (registerForm.password != registerForm.repeatPassword) {
            setRegisterError("Passwords must match")
            return
        }
        ApiPostRequest({
            body: {
                username: registerForm.username,
                email: registerForm.email,
                password: registerForm.password
            },
            route: "/auth/register",
            responseHandler: (res) => {
                alert("User created succesfully!")
            },
            errorHandler: (err) => {
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
                    <h2 className="text-3xl mt-4 text-primary-50 font-sans">WDML</h2>
                    {
                        (currentState == PageState.PS_LOGIN) && (
                            <div className="flex flex-col items-center justify-center px-6 py-8 my-8 lg:py-0">
                                <div className="p-6 space-y-4 md:space-y-6 sm:p-8">
                                    <h1 className="text-xl font-bold leading-tight tracking-tight text-white md:text-2xl">
                                        Sign in to your account
                                    </h1>
                                    <form className="space-y-4 md:space-y-6 max-w-min min-w-96" action="#">
                                        <TextFieldWithLabel
                                            onChange={(e: ChangeEvent<HTMLInputElement>) => setLoginForm({ ...loginForm, username: e.target.value })}
                                            size={40}
                                            id="username"
                                            label="Username"
                                            placeholder="Username" />
                                        <TextFieldWithLabel
                                            onChange={(e: ChangeEvent<HTMLInputElement>) => setLoginForm({ ...loginForm, password: e.target.value })}
                                            size={40}
                                            id="password"
                                            label="Password"
                                            placeholder="**********" />
                                        <div className="flex items-center justify-between">
                                            <Checkbox>Remember me</Checkbox>
                                            <a href="#" className="text-sm font-medium text-secondary-600 hover:underline">Forgot password?</a>
                                        </div>
                                        <div className="text-sm font-light text-red-400">{loginError}</div>
                                        <Button fullWidth icon="arrow" onClick={sendLoginRequest}>Sign in</Button>
                                        <p className="text-sm font-light text-gray-400 justify-center flex items-center">
                                            {"New user?"}
                                            <a href="#" className="font-medium ml-1 text-secondary-600 hover:underline"
                                                onClick={() => setCurrentState(PageState.PS_REGISTER)}>
                                                Sign up
                                            </a>
                                        </p>
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
                                    <form className="space-y-4 md:space-y-6 max-w-min min-w-96" action="#">
                                        <TextFieldWithLabel
                                            onChange={(e: ChangeEvent<HTMLInputElement>) => setRegisterForm({ ...registerForm, username: e.target.value })}
                                            size={40}
                                            id="username"
                                            label="Username"
                                            placeholder="Username" />
                                        <TextFieldWithLabel
                                            onChange={(e: ChangeEvent<HTMLInputElement>) => setRegisterForm({ ...registerForm, email: e.target.value })}
                                            size={40}
                                            htmlFor="email"
                                            type="email"
                                            id="email"
                                            label="Email"
                                            placeholder="Email" />
                                        <TextFieldWithLabel
                                            onChange={(e: ChangeEvent<HTMLInputElement>) => setRegisterForm({ ...registerForm, password: e.target.value })}
                                            size={40}
                                            id="password"
                                            label="Password"
                                            placeholder="**********" />
                                        <TextFieldWithLabel
                                            onChange={(e: ChangeEvent<HTMLInputElement>) => setRegisterForm({ ...registerForm, repeatPassword: e.target.value })}
                                            size={40}
                                            id="repeatpassword"
                                            label="Repeat password"
                                            placeholder="**********" />
                                        <p className="text-sm font-light text-red-400">{registerError}</p>
                                        <Button
                                            onClick={sendRegisterRequest}
                                            fullWidth
                                            icon="arrow">Sign up</Button>
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