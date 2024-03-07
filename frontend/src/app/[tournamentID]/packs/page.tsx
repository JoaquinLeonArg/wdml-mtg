"use client"

import { CardDisplay, CardImage } from "@/components/card";
import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests";
import { BoosterPack, BoosterPackData, TournamentPlayer } from "@/types/tournamentPlayer";
import { Autocomplete, AutocompleteItem, Button, Input } from "@nextui-org/react";
import { useEffect, useState } from "react";
import { BsFillTrashFill } from "react-icons/bs";



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
  set: string
}[]

function AddPacks(props: AddPacksProps) {
  let [availablePacks, setAvailablePacks] = useState<BoosterPack[]>([])
  let [error, setError] = useState<string>("")
  let [addedPacks, setAddedPacks] = useState<AddedPacks>([])

  let sendAddPacksRequest = () => {
    setError("")
    ApiPostRequest({
      body: {
        booster_packs: addedPacks
      },
      route: "/tournament/" + props.tournamentID + "/boosters",
      responseHandler: (res) => {
        console.log(res)
      },
      errorHandler: (err) => {
        switch (err) {
          case "INVALID_AUTH":
            setError("Invalid credentials")
        }
      }
    })
  }

  useEffect(() => {
    ApiGetRequest({
      route: "/tournament/" + props.tournamentID + "/boosters",
      responseHandler: (res: { booster_packs: BoosterPack[] }) => {
        console.log(res)
        setAvailablePacks(res.booster_packs)
      },
      errorHandler: (err) => {
        setError(err)
      }
    })
  }, [props.tournamentID])

  return (
    <div>
      <form onSubmit={sendAddPacksRequest} className="flex flex-col gap-2 items-center">
        {/* TODO: Add more options, like players to add it to */}
        <div className="flex flex-row gap-2 items-center">
          <Input
            onChange={(e) => setAddedPacks([{ ...addedPacks[0], count: Number(e.target.value) }])}
            variant="bordered"
            type="number"
            label="Amount"
            placeholder="Number of packs"
            labelPlacement="inside"
            className="text-white"
            endContent={
              <div className="pointer-events-none flex items-center">
                <span className="text-gray-300 text-small">packs</span>
              </div>
            }
          />
          <Autocomplete
            onInputChange={(value) => setAddedPacks([{ ...addedPacks[0], set: availablePacks.find(v => `${v.set_code} - ${v.set_name}` == value)?.set_code || "" }])}
            id="set"
            label="Booster pack"
            labelPlacement="inside"
            placeholder="Select a boster pack"
            className="text-white max-w-xs"
            defaultItems={availablePacks.map((val) => { return { value: val.set_code, label: `${val.set_code} - ${val.set_name}` } })}
          >
            {(item) => <AutocompleteItem className="text-white" key={item.value}>{item.label}</AutocompleteItem>}
          </Autocomplete>
          <Button size="md" isIconOnly color="danger" aria-label="Like" disabled>
            <BsFillTrashFill className="w-6 h-6" />
          </Button>
        </div>
        <Button onPress={sendAddPacksRequest} size="md" className="w-full bg-primary-500" aria-label="Like">Add packs</Button>
      </form >
    </div >
  )
}