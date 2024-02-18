"use client"

import { AnimatePresence, motion } from "framer-motion"
import { useState } from 'react';
import Image from "next/image"
import { Button } from "@/components/buttons";
import { TextFieldWithLabel } from "@/components/field";
import { Checkbox } from "@/components/checkbox";

enum PageState {
    PS_LOGIN,
    PS_REGISTER
}


export default function Login() {
    let [currentState, setCurrentState] = useState<PageState>(PageState.PS_LOGIN)
    return (
        <div className='flex flex-row'>
            <div className='bg-background-300 relative w-full bg-[url(/intrude-on-the-mind.jpg)] bg-cover'>
            </div>
            <div className='flex h-screen justify-center min-w-[50%]'>
                <div className='my-auto px-8 flex flex-col items-center'>
                    <Image src="/logo.png" alt="" width={100} height={100}></Image>
                    <h2 className="text-3xl mt-4 text-primary-50 font-sans">WDML</h2>
                    {
                        (currentState == PageState.PS_LOGIN) && (
                            <LoginForm toSignUp={() => { setCurrentState(PageState.PS_REGISTER) }} />
                        )
                    }
                    {
                        (currentState == PageState.PS_REGISTER) && (
                            <RegisterForm toSignIn={() => { setCurrentState(PageState.PS_LOGIN) }} />
                        )
                    }
                </div>
            </div>
        </div>
    );
}

function LoginForm({ toSignUp }: any) {
    return (
        <div className="flex flex-col items-center justify-center px-6 py-8 my-8 mx-auto lg:py-0">
            <div className="p-6 space-y-4 md:space-y-6 sm:p-8">
                <h1 className="text-xl font-bold leading-tight tracking-tight text-white md:text-2xl">
                    Sign in to your account
                </h1>
                <form className="space-y-4 md:space-y-6" action="#">
                    <TextFieldWithLabel size={40} id="username" label="Username" placeholder="Username" />
                    <TextFieldWithLabel size={40} id="password" label="Password" placeholder="**********" />
                    <div className="flex items-center justify-between">
                        <Checkbox>Remember me</Checkbox>
                        <a href="#" className="text-sm font-medium text-secondary-600 hover:underline">Forgot password?</a>
                    </div>
                    <Button fullWidth icon="arrow">Sign in</Button>
                    <p className="text-sm font-light text-gray-400 justify-center flex items-center">
                        {"New user?"}<a href="#" className="font-medium ml-1 text-secondary-600 hover:underline" onClick={() => toSignUp()}>Sign up</a>
                    </p>
                </form>
            </div >
        </div >
    )
}

function RegisterForm({ toSignIn }: any) {
    return (
        <div className="flex flex-col items-center justify-center px-6 py-8 my-8 mx-auto lg:py-0">
            <div className="p-6 space-y-4 md:space-y-6 sm:p-8">
                <h1 className="text-xl font-bold leading-tight tracking-tight text-white md:text-2xl">
                    Sign up
                </h1>
                <form className="space-y-4 md:space-y-6" action="#">
                    <TextFieldWithLabel size={40} id="username" label="Username" placeholder="Username" />
                    <TextFieldWithLabel size={40} htmlFor="email" type="email" id="email" label="Email" placeholder="Email" />
                    <TextFieldWithLabel size={40} id="password" label="Password" placeholder="**********" />
                    <TextFieldWithLabel size={40} id="repeatpassword" label="Repeat password" placeholder="**********" />
                    <Button fullWidth icon="arrow">Sign up</Button>
                    <a href="#" className="text-sm justify-center flex items-center font-medium ml-1 text-secondary-600 hover:underline" onClick={() => toSignIn()} > Back to sign in</a>
                </form>
            </div >
        </div >
    )
}