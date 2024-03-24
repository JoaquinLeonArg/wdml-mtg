"use client"

import { Input } from "@nextui-org/react";
import { ApiPostRequest } from "@/requests/requests";
import { Button } from "@nextui-org/react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";

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
    ApiPostRequest({
      route: "/tournament_player",
      body: {
        tournament_code: joinCode,
      },
      errorHandler: (err) => {
        switch (err) {
          case "NOT_FOUND":
            setJoinError("Invite code is invalid"); break
          case "INTERNAL":
            setJoinError("An error ocurred"); break
          case "DUPLICATED_RESOURCE":
            setJoinError("You already joined this tournament"); break
        }
      },
      responseHandler: (res) => {
        router.push("/" + res.tournament_id)
      }
    })
  }

  let sendCreateRequest = () => {
    setCreateError("")
    if (createName == "") {
      setCreateError("Tournament name can't be empty")
      return
    }
    ApiPostRequest({
      route: "/tournament",
      body: {
        name: createName,
      },
      errorHandler: (err) => { setCreateError(err) },
      responseHandler: (res) => {
        router.push("/" + res.tournament_id)
      }
    })
  }

  return (
    <div className="flex flex-col items-center justify-center h-screen">
      <div className="w-[50%] p-2 flex flex-row items-center justify-between bg-gray-600 rounded-lg">
        <div className="flex flex-col w-[50%] h-full p-2">
          <h1 className="text-xl font-bold leading-tight tracking-tight text-white md:text-2xl">Join</h1>
          <h3 className="text-sm font-bold leading-tight tracking-tight text-white opacity-70 md:text-base mb-8">Join an existing tournament using an invitation code.</h3>
          <Input
            id="join-code"
            placeholder="xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
            label="Invitation code"
            onValueChange={(value) => setJoinCode(value)} />
          <p className="text-sm font-light text-red-400 h-2">{joinError}</p>
          <div className="h-4"></div>
          <Button color={joinCode.length > 0 ? "success" : "default"} onClick={sendJoinRequest}>Join</Button>
        </div>
        <div className="w-1 mx-2 h-[90%] bg-white opacity-20"></div>
        <div className="flex flex-col w-[50%] h-full p-2">
          <h1 className="text-xl font-bold leading-tight tracking-tight text-white md:text-2xl">Create</h1>
          <h3 className="text-sm font-bold leading-tight tracking-tight text-white opacity-70 md:text-base mb-8">Create a new tournament and invite people to join.</h3>
          <Input
            isDisabled
            id="name"
            placeholder="My cool tournament!"
            label="Tournament name"
            onValueChange={(value) => setCreateName(value)} />
          <p className="text-sm font-light text-red-400 h-2">{createError}</p>
          <div className="h-4"></div>
          <Button isDisabled color={createName.length > 0 ? "success" : "default"} onClick={sendCreateRequest}>Create</Button>
        </div>
      </div>
      <div className="mt-2">
        <Link href="/" className="font-medium ml-1 text-secondary-600 hover:underline"> Back to tournaments </Link>
      </div>
    </div >
  )
}