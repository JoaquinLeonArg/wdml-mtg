"use client"

import { ApiGetRequest } from "@/requests/requests";
import { useRouter } from "next/navigation";
import { useEffect } from "react";

export default function Home() {
  let router = useRouter()

  useEffect(() => {
    ApiGetRequest({
      route: "/tournament_player",
      responseHandler: (res) => {
        if (!res.tournament_players) {
          router.push("/join")
          return
        }
        router.push("/" + res.tournament_players[0].tournament_id)
      },
      errorHandler: (err) => {
        console.log(err)
      }
    })
  }, [])
}
