"use client"

import { Button } from "@/components/buttons";
import { TextFieldWithLabel } from "@/components/field";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { ChangeEvent, useState } from "react";

type JoinResponse = {
  tournamentCode: string
}

export default function Home() {
  let router = useRouter()
  let [joinCode, setJoinCode] = useState<string>("")
  let [joinError, setJoinError] = useState<string>("")
  let [createName, setCreateName] = useState<string>("")
  let [createError, setCreateError] = useState<string>("")

  let sendJoinRequest = () => {
    setJoinError("")
    if (joinCode == "") {
      setJoinError("Join code can't be empty")
      return
    }
    fetch(`http://localhost:8080/tournament_player`, {
      method: "POST",
      credentials: 'include',
      body: JSON.stringify({
        username: joinCode,
      }),
    }).
      then((res: Response) => res.json()).
      then((res: JoinResponse) => {
        router.push("/" + res.tournamentCode)
      }).catch((reason: string) => { setJoinError("Something went wrong: " + reason) })
  }

  let sendCreateRequest = () => {
    setCreateError("")
    if (joinCode == "") {
      setCreateError("Tournament name can't be empty")
      return
    }
    fetch(`http://localhost:8080/tournament`, {
      method: "POST",
      credentials: 'include',
      body: JSON.stringify({
        name: createName,
      }),
    }).
      then((res: Response) => res.json()).
      then((res: JoinResponse) => {
        router.push("/" + res.tournamentCode)
      }).catch((reason: string) => { setJoinError("Something went wrong: " + reason) })
  }

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <div className="w-[50%] p-2 flex flex-row items-center justify-between bg-slate-600 rounded-lg">
        <div className="flex flex-col w-[50%] h-full p-2">
          <h1 className="text-xl font-bold leading-tight tracking-tight text-white md:text-2xl">Join</h1>
          <h3 className="text-sm font-bold leading-tight tracking-tight text-white opacity-70 md:text-base mb-8">Join an existing tournament using an invitation code.</h3>
          <TextFieldWithLabel
            id="join-code"
            placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
            label="Invitation code"
            onChange={(e: ChangeEvent<HTMLInputElement>) => setJoinCode(e.target.value)} />
          <div className="h-4"></div>
          <Button onClick={() => sendJoinRequest()}>Join</Button>
        </div>
        <div className="w-1 mx-2 h-[90%] bg-white opacity-20"></div>
        <div className="flex flex-col w-[50%] h-full p-2">
          <h1 className="text-xl font-bold leading-tight tracking-tight text-white md:text-2xl">Create</h1>
          <h3 className="text-sm font-bold leading-tight tracking-tight text-white opacity-70 md:text-base mb-8">Create a new tournament and invite people to join.</h3>
          <TextFieldWithLabel
            id="name"
            placeholder="My cool tournament!"
            label="Tournament name"
            onChange={(e: ChangeEvent<HTMLInputElement>) => setCreateName(e.target.value)} />
          <div className="h-4"></div>
          <Button onClick={() => sendCreateRequest()}>Create</Button>
        </div>
      </div>
      <div className="mt-2">
        <Link href="#" className="font-medium ml-1 text-secondary-600 hover:underline"> Back to tournaments </Link>
      </div>
    </div >
  )
}