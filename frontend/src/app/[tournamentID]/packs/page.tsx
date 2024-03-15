"use client"

import { Header } from "@/components/header";
import Layout from "@/components/layout";
import { PackList } from "@/components/packlist";
import { ApiGetRequest, ApiPostRequest } from "@/requests/requests";
import { CardData } from "@/types/card";
import { BoosterPack, BoosterPackData, TournamentPlayer } from "@/types/tournamentPlayer";
import { Autocomplete, AutocompleteItem, Button, ButtonGroup, Dropdown, DropdownItem, DropdownMenu, DropdownTrigger, Input, Listbox, ListboxItem, Spinner } from "@nextui-org/react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { BsFillTrashFill } from "react-icons/bs";
import Image from "next/image"
import { OwnedBoosterPack } from "@/types/boosterPack";

export default function PacksPage(props: any) {
  let router = useRouter()
  let [tournamentPlayer, setTournamentPlayer] = useState<TournamentPlayer>()
  let [error, setError] = useState<string>("")
  let [currentCards, setCurrentCards] = useState<CardData[]>([])
  let [boostersLoading, setBoostersLoading] = useState<boolean>(false)
  let [boostersVisible, setBoostersVisible] = useState<boolean>(true)
  let [allFlipped, setAllFlipped] = useState<boolean>(false)


  useEffect(() => {
    refreshData()
  }, [])

  let refreshData = () => {
    setBoostersLoading(true)
    ApiGetRequest({
      route: "/tournament_player/tournament",
      query: { tournamentID: props.params.tournamentID },
      errorHandler: (err) => {
        setError(err)
        setBoostersLoading(false)
      },
      responseHandler: (res: { tournament_player: TournamentPlayer }) => {
        setBoostersLoading(false)
        setTournamentPlayer(res.tournament_player)
      }
    })
  }

  let sendOpenPackRequest = (booster_pack_data: BoosterPackData) => {
    setBoostersVisible(false)
    setBoostersLoading(true)
    ApiPostRequest({
      route: `/boosterpacks/open`,
      query: { tournamentID: props.params.tournamentID },
      body: {
        booster_pack_data
      },
      errorHandler: (err) => {
        setBoostersVisible(true)
        setBoostersLoading(false)
        setError(err)
      },
      responseHandler: (res: { card_data: CardData[] }) => {
        setBoostersLoading(false)
        setCurrentCards(res.card_data)
      }
    })
  }

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
                <AddPacks tournamentID={props.params.tournamentID} refreshFunc={refreshData} />
              </div>
            </>
          )
        }
        <Header title="Open packs" />
        {boostersLoading ? <div className="flex justify-center"> <Spinner /></div> :
          <div className="flex flex-row gap-2 justify-center">
            {boostersVisible && (
              <div className="bg-gray-800 w-[450px] border-small px-1 py-2 rounded-small border-default-200">
                <Listbox>
                  {
                    tournamentPlayer.game_resources.booster_packs.map((booster_pack) =>
                      <ListboxItem
                        className="text-white"
                        onPress={() => sendOpenPackRequest(booster_pack.data)}
                        key={booster_pack.data.set_code}
                        startContent={<div className="text-gray-500 w-16">{booster_pack.data.set_code}</div>}
                        endContent={<div className="text-gray-500 text-right">{`(${booster_pack.available})`}</div>}
                      >
                        {`${booster_pack.data.set_code} - ${booster_pack.data.set_name}`}
                      </ListboxItem>
                    )
                  }

                </Listbox>
              </div>
            )}
            {currentCards.length && (
              <>
                <div className="flex flex-col gap-8 justify-center">
                  <ButtonGroup>
                    <Button color="success" onClick={() => setAllFlipped(true)}>Flip all</Button>
                    <Button color="danger" onClick={() => { setCurrentCards([]); refreshData(); setAllFlipped(false); setBoostersVisible(true) }}>Close</Button>
                  </ButtonGroup>
                  <CardDisplay cardsData={currentCards} allFlipped={allFlipped} />
                </div>
              </>
            )}
          </div>
        }
      </div >
    </Layout >
  )
}

type AddPacksProps = {
  tournamentID: string
  refreshFunc: () => void
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
  let [loading, setLoading] = useState<boolean>(false)

  let isButtonDisabled = addedPacks.length == 0 || addedPacks[0].count <= 0 || !addedPacks[0].set

  let sendAddPacksRequest = () => {
    setError("")
    setLoading(true)
    ApiPostRequest({
      body: {
        booster_packs: addedPacks
      },
      route: "/boosterpacks",
      query: { tournamentID: props.tournamentID },
      responseHandler: (_) => {
        setLoading(false)
        props.refreshFunc()
      },
      errorHandler: (err) => {
        setLoading(false)
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
    refreshAvailableBoosters()
  }, [])

  let refreshAvailableBoosters = () => {
    ApiGetRequest({
      route: "/boosterpacks",
      query: { tournamentID: props.tournamentID },
      responseHandler: (res: { booster_packs: BoosterPackData[] }) => {
        setAvailablePacks(res.booster_packs)
      },
      errorHandler: (err) => {
        setError(err)
      }
    })
  }

  return (
    <div>
      <form onSubmit={sendAddPacksRequest} className="flex flex-col gap-2 items-center">
        <div className="flex flex-row gap-2 items-center">
          <Input
            disabled={loading}
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
            disabled={loading}
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
        <Button isLoading={loading} onPress={sendAddPacksRequest} disabled={isButtonDisabled} size="md" className={`w-full ${isButtonDisabled || "bg-primary-500"}`} aria-label="Like">Add packs</Button>
      </form >
      <p className="text-sm font-light text-red-400">{error}</p>
    </div >
  )
}


export type CardImageProps = CardData & {
  startFaceUp?: boolean
}

export function CardImage(props: CardImageProps) {
  let [isFaceUp, setIsFaceUp] = useState<boolean>(props.startFaceUp || false)
  let borderRarityColor = {
    "common": "border-rarity-common",
    "uncommon": "border-rarity-uncommon",
    "rare": "border-rarity-rare",
    "mythic": "border-rarity-mythic",
  }[props.rarity]
  let shadowRarityColor = {
    "common": "shadow-rarity-common",
    "uncommon": "shadow-rarity-uncommon",
    "rare": "shadow-rarity-rare",
    "mythic": "shadow-rarity-mythic",
  }[props.rarity]
  return (
    <div className="group w-[256px] h-[355px] hover:scale-110 will-change-transform scale-100 duration-75 z-[100] hover:z-[110] [perspective:1000px]">
      <div onClick={() => { if (!isFaceUp) setIsFaceUp(true) }} className={
        `absolute rounded-xl w-full h-full duration-500 transition-all [transform-style:preserve-3d] ${!isFaceUp && "[transform:rotateY(180deg)]"}`
      }>
        <div className="absolute inset-0 [backface-visibility:hidden]">
          <Image className={`duration-75 border-2 ${borderRarityColor} rounded-xl ${isFaceUp && "shadow-[0px_0px_20px_1px_rgba(0,0,0,0.3)]"} ${shadowRarityColor}`} unoptimized priority src={props.image_url} alt="" width={1024} height={768} quality={100} />
        </div>
        <div className="absolute inset-0 [backface-visibility:hidden] [transform:rotateY(180deg)]">
          <Image className="duration-75 border-2 border-white rounded-xl" src="/cardback.webp" alt="" width={1024} height={1024} quality={100} layout="" />
        </div>
      </div>
    </div >
  )
}

export type CardDisplayProps = {
  cardsData: CardData[]
  allFlipped: boolean
}

export function CardDisplay(props: CardDisplayProps) {
  return (
    <div className="flex flex-wrap flex-row gap-2 items-center justify-center">
      {props.cardsData.map((cardData: CardData, i: number) => {
        return (
          <CardImage key={`${props.allFlipped}-${cardData.image_url}`} {...cardData} startFaceUp={props.allFlipped} />
        )
      })}
    </div>
  )
}