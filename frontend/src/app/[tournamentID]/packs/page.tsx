"use client"

import { CardDisplay, CardImage } from "@/components/card";
import { ComboBox } from "@/components/combobox";
import { CounterInput } from "@/components/counter";
import { TextFieldWithLabel } from "@/components/field";
import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { ApiGetRequest } from "@/requests/requests";
import { BoosterPack, BoosterPackData, TournamentPlayer } from "@/types/tournamentPlayer";
import { useEffect, useState } from "react";



export default function PacksPage(props: any) {
  let [tournamentPlayer, setTournamentPlayer] = useState<TournamentPlayer>()
  let [error, setError] = useState<string>("")

  useEffect(() => {
    ApiGetRequest({
      route: "/tournament_player",
      errorHandler: (err) => { setError(err) },
      responseHandler: (res: { tournament_players: TournamentPlayer[] }) => {
        console.log(res)
        let thisTournamentPlayer = res.tournament_players.filter((tp) => tp.tournament_id == props.params.tournamentID)[0]
        if (!thisTournamentPlayer) {
          setError("An error ocurred")
          return
        }
        setTournamentPlayer(thisTournamentPlayer)
      }
    })
  }, [])

  if (!tournamentPlayer) {
    return ""
  }

  return (
    <Layout tournamentID={props.params.tournamentID}>
      <div className="mx-16 my-16">
        {
          tournamentPlayer.access_level == "al_administrator" &&
          (
            <>
              <Header title="Add packs" />
              <div className="flex flex-row gap-2 justify-left">
                {/* <PackList packs={tournamentPlayer} /> */}
                <AddPacks tournamentID={props.params.tournamentID} />
              </div>
            </>
          )
        }
        <Header title="Open packs" />
        <div className="flex flex-row gap-2 justify-center">
          {/* <PackList packs={tournamentPlayer} /> */}
          <CardDisplay cardImageURLs={[]} />
        </div>
      </div >
    </Layout >
  )
}

type AddPacksProps = {
  tournamentID: string
}

type AddedPacks = {
  count: number
  pack: BoosterPack
}[]

function AddPacks(props: AddPacksProps) {
  let [availablePacks, setAvailablePacks] = useState<BoosterPack[]>([])
  let [error, setError] = useState<string>("")
  let [addedPacks, setAddedPacks] = useState<AddedPacks>([])

  useEffect(() => {
    ApiGetRequest({
      route: "/tournament/" + props.tournamentID + "/boosters",
      responseHandler: (res: { booster_packs: BoosterPack[] }) => {
        setAvailablePacks(res.booster_packs)
      },
      errorHandler: (err) => {
        setError(err)
      }
    })
  }, [])

  return (
    <div>
      <form>
        {/* TODO: Add more options, like players to add it to */}
        <div className="flex flex-row">
          <CounterInput id="" initial={1} label="" max={99} min={1} onChange={(value) => console.log(value)} />
          <ComboBox default="test" data={["test", "test2", "other"]} />
        </div>
      </form>
    </div>
  )
}