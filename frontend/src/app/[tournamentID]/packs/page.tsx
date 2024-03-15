"use client"

import { CardDisplay } from "@/components/card";
import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests";
import { CardData } from "@/types/card";
import { BoosterPack, BoosterPackData, TournamentPlayer } from "@/types/tournamentPlayer";
import { Autocomplete, AutocompleteItem, Button, Input } from "@nextui-org/react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { BsFillTrashFill } from "react-icons/bs";



export default function PacksPage(props: any) {
  let router = useRouter()
  let [tournamentPlayer, setTournamentPlayer] = useState<TournamentPlayer>()
  let [error, setError] = useState<string>("")
  let [currentCards, setCurrentCards] = useState<CardData[]>([])

  useEffect(() => {
    ApiGetRequest({
      route: "/tournament_player",
      errorHandler: (err) => { setError(err) },
      responseHandler: (res: { tournament_players: TournamentPlayer[] }) => {
        console.log(res)
        if (!res.tournament_players) {
          router.push("/")
        }
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
  if (error) {
    return error
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
                <AddPacks tournamentID={props.params.tournamentID} />
              </div>
            </>
          )
        }
        <Header title="Open packs" />
        <div className="flex flex-row gap-2 justify-between">
          <PackList packs={tournamentPlayer.game_resources.booster_packs} openPackHandler={setCurrentCards} tournamentID={props.params.tournamentID} />
          <CardDisplay cardImageURLs={currentCards} />
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
  type: string
}[]

function AddPacks(props: AddPacksProps) {
  let [availablePacks, setAvailablePacks] = useState<BoosterPackData[]>([])
  let [error, setError] = useState<string>("")
  let [addedPacks, setAddedPacks] = useState<AddedPacks>([{ type: "bt_draft", count: 0, set: "" }])

  let isButtonDisabled = addedPacks.length == 0 || addedPacks[0].count == 0 || !addedPacks[0].set

  let sendAddPacksRequest = () => {
    setError("")
    ApiPostRequest({
      body: {
        booster_packs: addedPacks
      },
      route: "/boosterpacks",
      query: { tournamentID: props.tournamentID },
      responseHandler: (res) => {

      },
      errorHandler: (err) => {
        switch (err) {
          case "INVALID_AUTH":
            setError("Invalid credentials")
          case "NO_DATA":
            setError("Invalid number or type of packs")
          default:
            setError("An error ocurred")
        }
      }
    })
  }

  useEffect(() => {
    ApiGetRequest({
      route: "/boosterpacks",
      query: { tournamentID: props.tournamentID },
      responseHandler: (res: { booster_packs: BoosterPackData[] }) => {
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
          <Button size="md" isIconOnly aria-label="Like" disabled>
            <BsFillTrashFill className="w-6 h-6" />
          </Button>
        </div>
        <Button onPress={sendAddPacksRequest} disabled={isButtonDisabled} size="md" className={`w-full ${isButtonDisabled || "bg-primary-500"}`} aria-label="Like">Add packs</Button>
      </form >
      <p className="text-sm font-light text-red-400">{error}</p>
    </div >
  )
}